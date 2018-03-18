package gokvadmin

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"errors"
	"reflect"
	"path/filepath"
	"go/build"
	"strings"
)

type GoKVAdmin struct {
	Engine Engine
	Config Config
	Server *http.Server
	Router *mux.Router
}

var (
	enginesMap = map[string]Engine{}
	assetsPath = ""
)

func RegisterEngine(name string, engine Engine) {
	enginesMap[name] = engine
}

var errUnknownEngine = errors.New("Unknow engine")

const (
	engineBoltDB = "boltdb"
	engineBadger = "badger"
	gkvImportPath = "github.com/sigurniv/gokvadmin"
)

func NewGoKVAdmin(engine string, db interface{}, config Config) (GoKVAdmin, error) {
	var gkv GoKVAdmin
	engineInstance, ok := enginesMap[engine]
	if !ok {
		return gkv, errUnknownEngine
	}

	e := reflect.New(reflect.TypeOf(engineInstance).Elem()).Interface().(Engine)
	e.SetDB(db)

	gkv = GoKVAdmin{Engine: e}
	err := gkv.setupServer(config)

	return gkv, err
}

func (gkv *GoKVAdmin) setupServer(config Config) error {
	var err error

	_, err = setupAssetsPath()
	if err != nil {
		return err
	}

	router := gkv.setupRouter()

	// add custom engine routes
	if instance, ok := gkv.Engine.(RoutedEngine); ok {
		instance.AddRoutes(router)
	}

	gkv.Config = config
	gkv.Router = router

	gkv.Server = &http.Server{
		Addr: ":" + config.Port,
		Handler: gkv.Router,
	}

	return err
}

func (gkv *GoKVAdmin) setupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", gkv.HomeHandler)
	r.SkipClean(true)

	//get the value by key & optionally bucket
	r.HandleFunc("/api/auth", gkv.Auth).Methods(http.MethodPost)

	//get the value by key & optionally bucket
	r.HandleFunc("/api/init", gkv.Init).Methods(http.MethodGet)

	//get the value by key & optionally bucket
	r.HandleFunc("/api/key/{key}", gkv.GetKey).Methods(http.MethodGet)

	//set the value by key & optionally bucket
	r.HandleFunc("/api/key/{key}", gkv.SetKey).Methods(http.MethodPost)

	//delete the value by key & optionally bucket
	r.HandleFunc("/api/key/{key}", gkv.DeleteKey).Methods(http.MethodDelete)

	//get the value by key prefix & optionally bucket
	r.HandleFunc("/api/prefix/key/{key}", gkv.GetKeyByPrefix).Methods(http.MethodGet)
	r.HandleFunc("/api/prefix/key/", gkv.GetKeyByPrefix).Methods(http.MethodGet)

	r.NotFoundHandler = http.HandlerFunc(gkv.HomeHandler)

	//static file server for web interface
	r.PathPrefix("/dist/").Handler(http.FileServer(http.Dir(assetsPath)))

	r.Use(gkv.AuthMiddleware)

	return r
}

func (gkv *GoKVAdmin) Run() error {
	var err error
	if gkv.Config.TLS != nil {
		err = gkv.Server.ListenAndServeTLS(gkv.Config.TLS.CertFile, gkv.Config.TLS.KeyFile)
	} else {
		err = gkv.Server.ListenAndServe()
	}

	return err
}

func (gkv *GoKVAdmin) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if gkv.Config.Auth == nil {
			next.ServeHTTP(w, r)
			return
		}

		if r.URL.Path == "/api/auth" || r.URL.Path == "/" || strings.HasPrefix(r.URL.Path, "/dist") {
			next.ServeHTTP(w, r)
			return
		}

		token := r.Header.Get("Authorization")
		if token != gkv.Config.Auth.token || gkv.Config.Auth.tokenExpires.Before(time.Now().Local()) {
			gkv.Config.Auth.token = ""
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func setupAssetsPath() (string, error) {
	appPkg, err := build.Import(gkvImportPath, "", build.FindOnly)
	if err != nil {
		return "", err
	}

	path := filepath.Join(appPkg.Dir, "assets")
	assetsPath = path
	return path, err
}