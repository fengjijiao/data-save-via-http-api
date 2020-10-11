package httpLib

import (
	"net/http"
	"github.com/gorilla/mux"
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
		"Css",
		"GET",
		"/css/{cssName}.css",
		CssStaticHttpHandler,
	},
	Route {
		"Js",
		"GET",
		"/js/{jsName}.js",
		JsStaticHttpHandler,
	},
	Route {
		"Index",
		"GET",
		"/",
		IndexHttpHandler,
	},
	Route {
		"Login",
		"GET",
		"/login",
		LoginHttpHandler,
	},
	Route {
		"Login",
		"POST",
		"/login",
		LoginHttpHandler,
	},
	Route {
		"Register",
		"GET",
		"/register",
		RegisterHttpHandler,
	},
	Route {
		"Register",
		"POST",
		"/register",
		RegisterHttpHandler,
	},
	Route {
		"CreateDataSet",
		"POST",
		"/dataSet/new",
		CreateDataSetHttpHandler,
	},
	Route {
		"GetVal",
		"GET",
		"/dataSet/{did}",
		GetValHttpHandler,
	},
	Route {
		"PutVal",
		"PUT",
		"/dataSet/{did}",
		PutValHttpHandler,
	},
	Route {
		"DelVal",
		"DELETE",
		"/dataSet/{did}",
		DelValHttpHandler,
	},
}