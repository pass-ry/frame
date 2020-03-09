package router

import (
	"net/http"
	_ "net/http/pprof"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.ifchange.com/data/cordwood/log"
	handler "gitlab.ifchange.com/data/cordwood/rpc/rpc-handler"
)

var (
	_ http.Handler = (*Router)(nil)

	_ http.Handler = (*notFoundHandler)(nil)
	_ http.Handler = (*methodNotAllowedHandler)(nil)
)

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) { router.router.ServeHTTP(w, r) }

type Router struct{ router *mux.Router }

func (router *Router) Handler(path string, handlerFunc handler.Handler) *Router {
	router.router.HandleFunc(path, handler.Wrap(handlerFunc)).
		Methods(http.MethodPost)
	return router
}

func NewRouter() *Router {
	router := &Router{router: mux.NewRouter()}

	router.router.NotFoundHandler = new(notFoundHandler)
	router.router.MethodNotAllowedHandler = new(methodNotAllowedHandler)
	return router
}

func (r *Router) WithPPROF() *Router {
	r.router.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)
	return r
}

var (
	heartbeat = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "heartbeat",
			Help: "node heartbeat",
		})
	once sync.Once
)

func (r *Router) WithMetrics() *Router {
	r.router.Handle("/metrics", promhttp.Handler())
	once.Do(func() {
		go func() {
			ticker := time.Tick(time.Duration(10) * time.Second)
			for {
				heartbeat.SetToCurrentTime()
				<-ticker
			}
		}()
	})

	return r
}

type (
	notFoundHandler         struct{}
	methodNotAllowedHandler struct{}
)

func (h *notFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r.Body != nil {
			r.Body.Close()
		}
	}()

	log.Errorf("HTTP URI %s PATH Not Found",
		r.RequestURI)

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("PATH NOT FOUND"))
}

func (h *methodNotAllowedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r.Body != nil {
			r.Body.Close()
		}
	}()

	log.Errorf("HTTP URI %s Method %s Not Allowed",
		r.RequestURI, r.Method)

	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte("METHOD NOT ALLOWED"))
}
