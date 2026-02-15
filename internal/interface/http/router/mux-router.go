package router

import (
	"maqhaa/library/logging"
	"maqhaa/library/middleware"
	"net/http"

	"github.com/gorilla/mux"
	//"github.com/gorilla/mux"
)

var (
	muxDispatcher = mux.NewRouter()
)

type muxRouter struct{}

func NewMuxRouter() Router {
	return &muxRouter{}
}

func (*muxRouter) GET(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	muxDispatcher.HandleFunc(uri, f).Methods("GET")
}
func (*muxRouter) POST(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	muxDispatcher.HandleFunc(uri, f).Methods("POST")
}

func (*muxRouter) PUT(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	muxDispatcher.HandleFunc(uri, f).Methods("PUT")
}

func (*muxRouter) DELETE(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	muxDispatcher.HandleFunc(uri, f).Methods("DELETE")
}

func (*muxRouter) GetRouter() *mux.Router {
	return muxDispatcher
}

func (*muxRouter) SERVE(port string) {
	logging.Log.Infof("Http server listen in port %s", port)
	//muxDispatcher.Use(apmgorilla.Middleware())
	muxDispatcher.Use(middleware.LoggingMiddleware)
	http.ListenAndServe(port, muxDispatcher)
}
