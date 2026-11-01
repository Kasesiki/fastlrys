package fastlrys

import (
	tool "github.com/Xiaxiaobaii/autotool"

	"github.com/valyala/fasthttp"
	"golang.org/x/time/rate"
)

type Server struct {
	addr        string
	Port        int
	TlsPort     int
	path        map[Path]func(ctx *RequestCtx)
	FastServer  *fasthttp.Server
	https       *TLSoption[string]
	NoRequest   bool
	NoRequ      func(ctx *RequestCtx)
	MiddleWare  map[string]func(ctx *RequestCtx)
	maxBodySize int
	ifHttp      bool
	listenIpv6  bool
	HSTS 	    bool
	QOS			QosSetting
}

type QosSetting struct {
	QosNumber   int
	QosBucket   map[string]*rate.Limiter
	QosBursts   int

}

type RequestCtx struct {
	FastCtx  *fasthttp.RequestCtx
	QueryCtx Path
	Method   HttpType
}

type Path struct {
	path string
	typ  HttpType
}

type TLSoption[T any] struct {
	tls     *tool.Kvalue[T]
	onlyTLS bool
}

type HttpType string
type HttpStatus string

const (
	ANY     HttpType = "ANY"
	GET     HttpType = "GET"
	POST    HttpType = "POST"
	HEAD    HttpType = "HEAD"
	PUT     HttpType = "PUT"
	DELETE  HttpType = "DELETE"
	CONNECT HttpType = "CONNECT"
	OPTIONS HttpType = "OPTIONS"
	TRACE   HttpType = "TRACE"
	PATCH   HttpType = "PATCH"
)
