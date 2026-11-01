package fastlrys

import (
	"time"

	tool "github.com/Xiaxiaobaii/autotool"
	"github.com/valyala/fasthttp"
	"golang.org/x/time/rate"
)

func New(Port int) *Server {
	var Manager = &Server{
		Port:      Port,
		addr:      "0.0.0.0:",
		path:      make(map[Path]func(ctx *RequestCtx)),
		NoRequest: false,
		https: &TLSoption[string]{
			tls:     &tool.Kvalue[string]{},
			onlyTLS: false,
		},
		MiddleWare:  make(map[string]func(ctx *RequestCtx)),
		maxBodySize: 256 * 1024 * 1024,
		ifHttp:      true,
		listenIpv6:  true,
		HSTS:        false,
		QOS: QosSetting{
			QosNumber: 24,
			QosBucket: map[string]*rate.Limiter{},
			QosBursts: 48,
		},
	}
	Manager.FastServer = &fasthttp.Server{
		Handler:      Manager.handCall,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		Concurrency:  1000 * 500,
	}
	return Manager
}
