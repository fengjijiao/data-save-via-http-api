package httpLib

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/fengjijiao/data-save-via-http-api/pkg/conf"
	"path"
)

type Route struct {
	Name		string
	Method		string
	Pattern		string
	HandlerFunc	http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return router
}

var routes = Routes {
	Route {
		"Index",
		"GET",
		path.Join(conf.Config.BaseUrlPath, "/"),
		IndexHttpHandler,
	},
	Route {
		"Login",
		"POST",
		path.Join(conf.Config.BaseUrlPath, "/login"),
		LoginHttpHandler,
	},
	Route {
		"Register",
		"POST",
		path.Join(conf.Config.BaseUrlPath, "/register"),
		RegisterHttpHandler,
	},
	Route {
		"CreateDataSet",
		"POST",
		path.Join(conf.Config.BaseUrlPath, "/dataSet/new"),
		CreateDataSetHttpHandler,
	},
	Route {
		"ListDataSet",
		"GET",
		path.Join(conf.Config.BaseUrlPath, "/dataSet/list"),
		ListDataSetHttpHandler,
	},
	Route {
		"GetValJson",
		"GET",
		path.Join(conf.Config.BaseUrlPath, "/dataSet/{did}/json"),
		GetValJsonHttpHandler,
	},
	Route {
		"GetVal",
		"GET",
		path.Join(conf.Config.BaseUrlPath, "/dataSet/{did}"),
		GetValHttpHandler,
	},
	Route {
		"PutVal",
		"PUT",
		path.Join(conf.Config.BaseUrlPath, "/dataSet/{did}"),
		PutValHttpHandler,
	},
	Route {
		"DelVal",
		"DELETE",
		path.Join(conf.Config.BaseUrlPath, "/dataSet/{did}"),
		DelValHttpHandler,
	},
}