package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/rakyll/statik/fs"
	"github.com/santhosh-tekuri/jsonschema"
	comm "gitlab.lan.athonet.com/primo/susancalvin/common"
	"gitlab.lan.athonet.com/riccardo.manfrin/netconfd/logger"
	"gitlab.lan.athonet.com/riccardo.manfrin/netconfd/schemas"
	"gitlab.lan.athonet.com/riccardo.manfrin/netconfd/swaggerui"
)

//Manager of the machine takes HTTP requests, perform actions and give results back [blocking]
type Manager struct {
	Conf           Conf
	HTTPServer     *http.Server
	ServeMux       *http.ServeMux
	SyncHTTPServer *http.Server
	SyncServeMux   *http.ServeMux
	client         http.Client
	m2MEndpointURL string
	schema         *jsonschema.Schema
	openapi        *openapi3.Swagger
	router         *openapi3filter.Router
}

//PushContent describes the push content
type PushContent struct {
	FilePath       string `json:"filepath"`
	B64FileContent string `json:"b64filecontent"`
}

//File encodes the sync state of a file between local and remote
type File struct {
	Name      string `json:"name"`
	LocalMd5  string `json:"local_md5"`
	RemoteMd5 string `json:"remote_md5"`
	Aligned   bool   `json:"aligned"`
}

//FilesState encodes the state of all monitored files among remote and local
type FilesState struct {
	Files []File `json:"files"`
}

func urlMatch(in string, expected string) bool {
	return strings.Contains(in, expected)
}

func fail(err int, log string) {
	logger.Log.Error(log)
	os.Exit(err)
}

func httpError(w http.ResponseWriter, e error, code int) {
	logger.Log.Warning(e.Error())
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	io.WriteString(w, comm.ToJSONString(e))
}

func httpOk(w http.ResponseWriter) {
	w.WriteHeader(200)

}

func (m *Manager) overrideWithEnv() error {

	logger.Log.Info("Looking up NETCONFD_HOST from env")
	netconfdHost, netconfdHostFound := os.LookupEnv("NETCONFD_HOST")
	logger.Log.Info("Looking up NETCONFD_PORT from env")
	netconfdPort, netconfdPortFound := os.LookupEnv("NETCONFD_PORT")

	if netconfdHostFound {
		logger.Log.Info("Overriding .global.mgmt.host with " + netconfdHost + " NETCONFD_HOST env var")
		m.Conf.Global.Mgmt.Host = netconfdHost
	}
	if netconfdPortFound {

		overridePort, err := strconv.Atoi(netconfdPort)
		if err != nil {
			logger.Log.Info("Env NETCONFD_PORT " + netconfdPort + " is not a number")
			return err
		}
		if overridePort > int(math.Exp2(16))-1 {
			logger.Log.Info("Env NETCONFD_PORT " + netconfdPort + " is too high")
			return err
		}
		logger.Log.Info("Overriding .global.mgmt.port with " + netconfdPort + " NETCONFD_PORT env var")
		m.Conf.Global.Mgmt.Port = uint16(overridePort)
	}

	return nil
}

//LoadConfig loads the configuration
func (m *Manager) LoadConfig(conffile *string) error {
	js, err := os.Open(*configfile)
	if err != nil {
		fail(-1, "Missing config file...")
	}
	defer js.Close()

	if err := json.NewDecoder(js).Decode(&m.Conf); err != nil {
		fail(-1, "Bad configuration...")
	}

	logger.LoggerSetLevel(m.Conf.Global.LogLev)

	m.overrideWithEnv()

	res, _ := json.Marshal(m.Conf)
	logger.Log.Notice("Config: " + string(res))

	logger.Log.Notice("Starting mgmt service on " + m.Conf.Global.Mgmt.Host + ":" + strconv.FormatInt(int64(m.Conf.Global.Mgmt.Port), 10))
	logger.Log.Notice("Log level set to " + m.Conf.Global.LogLev)

	schemasFS, _ := fs.NewFromZipData(schemas.Init())
	swaggeruiFS, _ := fs.NewFromZipData(swaggerui.Init())

	f, _ := schemasFS.Open("/config.json")
	data, _ := ioutil.ReadAll(f)
	m.schema, err = jsonschema.CompileString("/config.json", string(data))
	if err != nil {
		logger.Fatal(err)
	}

	openapiJson, _ := swaggeruiFS.Open("/openapi.json")

	data, _ = ioutil.ReadAll(openapiJson)
	m.openapi, err = openapi3.NewSwaggerLoader().LoadSwaggerFromData(data)
	m.router = openapi3filter.NewRouter().WithSwagger(m.openapi)

	m.ServeMux = http.NewServeMux()
	m.HTTPServer = &http.Server{
		Addr:           m.Conf.Global.Mgmt.Host + ":" + strconv.FormatUint(uint64(m.Conf.Global.Mgmt.Port), 10),
		Handler:        m.ServeMux,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	m.ServeMux.Handle("/", m)
	m.ServeMux.Handle("/swaggerui/", http.StripPrefix("/swaggerui", http.FileServer(swaggeruiFS)))

	return nil
}

func (m *Manager) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	// Find route
	route, pathParams, _ := m.router.FindRoute(r.Method, r.URL)

	// Validate request
	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    r,
		PathParams: pathParams,
		Route:      route,
	}
	if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
		panic(err)
	}

	var (
		respStatus      = 200
		respContentType = "application/json"
		respBody        = bytes.NewBufferString(`{}`)
	)

	log.Println("Response:", respStatus)
	responseValidationInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: requestValidationInput,
		Status:                 respStatus,
		Header: http.Header{
			"Content-Type": []string{respContentType},
		},
	}
	if respBody != nil {
		data, _ := json.Marshal(respBody)
		responseValidationInput.SetBodyBytes(data)
	}

	// Validate response.
	if err := openapi3filter.ValidateResponse(ctx, responseValidationInput); err != nil {
		panic(err)
	}
	/*
		todo := false
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Log.Warning("Failed to process request body")
			return
		}
		if r.Method != "GET" {
			logger.Log.Notice("API:" + r.Method + ":" + r.URL.Path + ":" + string(body))
		} else {
			logger.Log.Info("API:" + r.Method + ":" + r.URL.Path + ":" + string(body))
		}

		if !strings.HasPrefix(r.URL.Path, "/api/1/") {
			logger.Log.Warning("Unsupported API:" + r.Method + "-" + r.URL.Path)
			http.NotFound(w, r)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/api/1")

		switch path {
		case "/state":
			if r.Method == "GET" {
				w.Header().Add("Content-Type", "application/json")
				io.WriteString(w, "ok")
			}
		case "/config":
			if r.Method == "PUT" {
				err = m.schema.Validate(strings.NewReader(string(body)))
				if err != nil {
					httpError(w, err, http.StatusBadRequest)
					return
				}
				w.Header().Add("Content-Type", "application/json")
				io.WriteString(w, "ok")
			}
		default:
			logger.Log.Warning("Unsupported API:" + r.Method + "-" + r.URL.Path)
			http.NotFound(w, r)
			return
		}
		if todo {
			logger.Log.Warning("Implement API:" + r.Method + "-" + r.URL.Path)
		}*/

}

//NewManager creates new manager
func NewManager(conffile *string) *Manager {
	m := &Manager{}
	m.LoadConfig(conffile)
	return m
}

//Start activates manager
func (m *Manager) Start() {

	logger.Log.Notice("Starting manager...")

	go m.HTTPServer.ListenAndServe()
}
