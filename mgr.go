package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rakyll/statik/fs"
	_ "gitlab.lan.athonet.com/riccardo.manfrin/netconfd/statik"

	"gitlab.lan.athonet.com/riccardo.manfrin/netconfd/logger"
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

func (m *Manager) overrideWithEnv() error {

	clusterLocal, clusterLocalFound := os.LookupEnv("CLUSTER_LOCAL")
	clusterRemote, clusterRemoteFound := os.LookupEnv("CLUSTER_REMOTE")

	if clusterLocalFound {
		logger.Log.Info("Overriding sync local host with " + clusterLocal + " CLUSTER_LOCAL env var")
		m.Conf.Sync.Local.Host = clusterLocal
	}
	if clusterRemoteFound {
		logger.Log.Info("Overriding sync remote host with " + clusterRemote + " CLUSTER_REMOTE env var")
		m.Conf.Sync.Remote.Host = clusterRemote
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

	statikFS, _ := fs.New()

	m.ServeMux = http.NewServeMux()
	m.HTTPServer = &http.Server{
		Addr:           m.Conf.Global.Mgmt.Host + ":" + strconv.FormatUint(uint64(m.Conf.Global.Mgmt.Port), 10),
		Handler:        m.ServeMux,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	m.ServeMux.Handle("/", m)
	m.ServeMux.Handle("/swaggerui/", http.StripPrefix("/swaggerui", http.FileServer(statikFS)))

	return nil
}

func (m *Manager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	default:
		logger.Log.Warning("Unsupported API:" + r.Method + "-" + r.URL.Path)
		http.NotFound(w, r)
		return
	}
	if todo {
		logger.Log.Warning("Implement API:" + r.Method + "-" + r.URL.Path)
	}

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
