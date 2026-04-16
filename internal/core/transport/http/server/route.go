package core_http_server

import "net/http"

type Route struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

func NewRoute(
	mathod string,
	path string,
	handler http.HandlerFunc,
) Route {
	return Route{
		Method:  mathod,
		Path:    path,
		Handler: handler,
	}
}
