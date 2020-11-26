package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	oas "gitlab.lan.athonet.com/core/netconfd/server/go"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	comm "gitlab.lan.athonet.com/core/netconfd/common"
	"gitlab.lan.athonet.com/core/netconfd/logger"
	"gitlab.lan.athonet.com/core/netconfd/swaggerui"
)

//RouteHandler is the handler for a certain route OperationID
type RouteHandler func(body interface{}) (interface{}, error)

//Manager of the machine takes HTTP requests, perform actions and give results back [blocking]
type Manager struct {
	Conf           Conf
	HTTPServer     *http.Server
	ServeMux       *http.ServeMux
	SyncHTTPServer *http.Server
	SyncServeMux   *http.ServeMux
	client         http.Client
	m2MEndpointURL string
	openapi        *openapi3.Swagger
	router         *openapi3filter.Router
	routeHandlers  *mux.Router
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

	logger.Log.Info("Loading config from " + *conffile)
	js, err := os.Open(*conffile)
	if err != nil {
		fail(-1, fmt.Sprintf("Config File %v access error [%v]: your network will not be configured. Aborting...", *conffile, err.Error()))
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

	m.HTTPServer = &http.Server{
		Addr:           m.Conf.Global.Mgmt.Host + ":" + strconv.FormatUint(uint64(m.Conf.Global.Mgmt.Port), 10),
		Handler:        m.ServeMux,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

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
		httpError(w, err, http.StatusBadRequest)
		return
	}

	wrw := &wrapperResponseWriter{
		ResponseWriter: w,
		buf:            &bytes.Buffer{},
	}

	m.routeHandlers.ServeHTTP(wrw, r)

	responseValidationInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: requestValidationInput,
		Status:                 wrw.Status,
		Header:                 w.Header(),
	}

	responseValidationInput.SetBodyBytes(wrw.buf.Bytes())

	// Validate response.
	if err := openapi3filter.ValidateResponse(ctx, responseValidationInput); err != nil {
		panic(err)
	}

	if (wrw.Status != http.StatusOK) && (wrw.Status != http.StatusCreated) {
		logger.Log.Warning(fmt.Sprintf("HTTP Error %v: %v", wrw.Status, string(wrw.buf.Bytes())))
	} else {
		logger.Log.Debug(fmt.Sprintf("HTTP Ok %v: %v", wrw.Status, string(wrw.buf.Bytes())))
	}
}

//NewManager creates new manager
func NewManager() *Manager {
	swaggeruiFS, _ := fs.NewFromZipData(swaggerui.Init())

	openapiJSON, _ := swaggeruiFS.Open("/openapi.yaml")

	data, _ := ioutil.ReadAll(openapiJSON)
	openapi, err := openapi3.NewSwaggerLoader().LoadSwaggerFromData(data)
	if openapi == nil && err != nil {
		panic(err)
	}
	router := openapi3filter.NewRouter().WithSwagger(openapi)

	//Opt in IPv4 and IPv6 validation
	openapi3.DefineIPv4Format()
	openapi3.DefineIPv6Format()

	serveMux := http.NewServeMux()

	networkApiService := oas.NewNetworkApiService()
	networkApiController := oas.NewNetworkApiController(networkApiService)
	systemApiService := oas.NewSystemApiService()
	systemAPIController := oas.NewSystemApiController(systemApiService)
	m := &Manager{
		openapi:       openapi,
		router:        router,
		ServeMux:      serveMux,
		routeHandlers: oas.NewRouter(networkApiController, systemAPIController),
	}
	serveMux.Handle("/", m)
	serveMux.Handle("/swaggerui/", http.StripPrefix("/swaggerui", http.FileServer(swaggeruiFS)))

	return m
}

//Start activates manager
func (m *Manager) Start() {

	logger.Log.Notice("Started netConfD manager...")

	go m.HTTPServer.ListenAndServe()
}

type wrapperResponseWriter struct {
	http.ResponseWriter
	buf    *bytes.Buffer
	Status int
}

func (wrw *wrapperResponseWriter) Write(p []byte) (int, error) {
	wrw.ResponseWriter.Write(p)
	return wrw.buf.Write(p)
}

//WriteHeader wrapper to record status
func (wrw *wrapperResponseWriter) WriteHeader(status int) {
	wrw.Status = status
	wrw.ResponseWriter.WriteHeader(status)
}
