package framework

import "net/http"

type Context struct {
	ResponseWriter http.ResponseWriter
	Request *http.Request
	Params map[string]string
}
