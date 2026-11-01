package fastlrys

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"

	tool "github.com/Xiaxiaobaii/autotool"
	"github.com/dgrr/http2"
	"github.com/fasthttp/websocket"

	"github.com/valyala/fasthttp"
	"golang.org/x/time/rate"
)

type CtxHandler func(ctx *RequestCtx)

func (server *Server) Do() error {
	var e error
	var ln net.Listener
	tool.FlagCall("debug", func(str string) {
		tool.Print(server.addr + tool.Itoa(server.Port))
	})
	if server.listenIpv6 {
		ln, e = net.Listen("tcp", server.addr+tool.Itoa(server.Port))
	} else {
		ln, e = net.Listen("tcp4", server.addr+tool.Itoa(server.Port))
	}

	if e != nil {
		return e
	}
	server.FastServer.MaxRequestBodySize = server.maxBodySize
	if k, v := server.https.tls.Get(); k != "" {
		tln, e := net.Listen("tcp", server.addr+tool.Itoa(server.TlsPort))
		if e != nil {
			return e
		}
		if server.ifHttp {
			go func() {
				_ = server.FastServer.Serve(ln)
			}()
		}
		http2.ConfigureServer(server.FastServer, http2.ServerConfig{})
		return server.FastServer.ServeTLS(tln, tool.IntoS(k), tool.IntoS(v))
	} else {
		return server.FastServer.Serve(ln)
	}
}

func (server *Server) SetAddr(addr string) {
	server.addr = addr
}

func (server *Server) SetShowHttp(Show bool) {
	server.ifHttp = Show
}

func (server *Server) ListenIPV6(Show bool) {
	server.listenIpv6 = Show
}

func (server *Server) AddPath(path string, typ HttpType, call CtxHandler) error {
	path = tool.UrlFormat(path)
	Pat := Path{
		path: path,
		typ:  typ,
	}
	if _, ok := server.path[Pat]; ok {
		return tool.Error("Path already exists")
	}
	server.path[Pat] = call
	tool.LogPrint("Path "+path+" Add", tool.INFO, 1)
	return nil
}

func (server *Server) SetNoRequest(call CtxHandler) {
	server.NoRequest = true
	server.NoRequ = call
}

func (server *Server) AddMiddle(Name string, middle CtxHandler) {
	server.MiddleWare[Name] = middle
}

var Limitmu sync.Mutex

func (server *Server) GetLimit(ip string) *rate.Limiter {
	Limitmu.Lock()
	defer Limitmu.Unlock()

	bucket, exist := server.QOS.QosBucket[ip]

	if !exist {
		server.QOS.QosBucket[ip] = rate.NewLimiter(rate.Limit(server.QOS.QosNumber), server.QOS.QosBursts)
		return server.QOS.QosBucket[ip]
	}

	return bucket

}

func (server *Server) handCall(ctx *fasthttp.RequestCtx) {
	if !server.GetLimit(ctx.RemoteIP().String()).Allow() {
		ctx.SetBodyString("别急, 先喝茶")
		ctx.SetStatusCode(418)
		return
	}
	path := tool.UrlFormat(string(ctx.Path()))
	fCtx := &RequestCtx{
		FastCtx: ctx,
		QueryCtx: Path{
			path: path,
		},
		Method: HttpType(ctx.Method()),
	}

	if server.HSTS {
		fCtx.AddHeader("Strict-Transport-Security", "max-age=300")
		if !strings.Contains(string(ctx.URI().FullURI()), "https://") {
			ctx.SetStatusCode(307)
			return
		}
	}

	for k, v := range server.path {
		if path == k.path || (strings.Contains(path, k.path[:len(k.path)-1]) && k.path[len(k.path)-1:] == "*") {
			if k.typ == HttpType(ctx.Method()) || k.typ == ANY {
				for _, v := range server.MiddleWare {
					v(fCtx)
				}
				v(fCtx)
			}
			return
		}
	}
	if server.NoRequest {
		for _, v := range server.MiddleWare {
			v(fCtx)
		}
		server.NoRequ(fCtx)
		return
	}
	ctx.SetStatusCode(404)

}

func (server *Server) SetTLS(tlsPort int, pem, key string) {
	server.TlsPort = tlsPort
	server.https.tls.Set(pem, key)
}

func (server *Server) SetMaxBodySize(BodySize int) {
	server.maxBodySize = BodySize
}

func (server *Server) AddGET(path string, call CtxHandler) {
	_ = server.AddPath(tool.UrlFormat(path), GET, call)
}

func (server *Server) AddWebSocket(path string, call WebSocketHandler) {
	_ = server.AddPath(tool.UrlFormat(path), GET, func(ctx *RequestCtx) {
		if CtxConvWs(ctx, call) != nil {
			ctx.SetStatus(504)
		}
	})
}

func (server *Server) AddPOST(path string, call CtxHandler) {
	_ = server.AddPath(tool.UrlFormat(path), POST, call)
}

func (server *Server) AddAny(path string, call CtxHandler) {
	_ = server.AddPath(tool.UrlFormat(path), ANY, call)
}

func (server *Server) AddRemotePath(path string, link string) {
	server.AddPath(tool.UrlFormat(path), ANY, func(ctx *RequestCtx) {
		pat := ctx.Path()
		pat = strings.TrimPrefix(pat, strings.ReplaceAll(path, "*", ""))
		header := make(http.Header)
		// for k, v := range ctx.FastCtx.Request.Header. {
		// 	header.Add(string(k), string(v))
		// }
		if ctx.FastCtx.Request.Header.ConnectionUpgrade() {
			header.Del("Sec-Websocket-Extensions")
			header.Del("Sec-Websocket-Key")
			header.Del("Sec-Websocket-Version")
			header.Del("Connection")
			header.Del("Upgrade")
			l := ""
			if ctx.FastCtx.TLSConnectionState() != nil {
				l = fmt.Sprintf("wss:%s%s?%s", strings.Join(strings.Split(link, ":")[1:], ":"), pat, ctx.QureyString())
			} else {
				l = fmt.Sprintf("ws:%s%s?%s", strings.Join(strings.Split(link, ":")[1:], ":"), pat, ctx.QureyString())
			}
			conn, resp, e := websocket.DefaultDialer.Dial(l, header)
			if e != nil {
				ctx.SetStatus(504)
				return
			}
			for key, value := range resp.Header {
				for _, v := range value {
					ctx.AddHeader(key, v)
				}
			}
			CtxConvWs(ctx, func(c *WebSocketConn) {
				defer conn.Close()
				readRemote(conn, c.conn)
			})

		} else {
			get, e := http.NewRequest(string(ctx.FastCtx.Method()), fmt.Sprintf("%s%s?%s", link, pat, ctx.QureyString()), bytes.NewReader(ctx.Body()))
			if e != nil {
				ctx.SetStatus(504)
				return
			}
			get.Header = header
			resq, e := http.DefaultClient.Do(get)
			if e != nil {
				ctx.SetStatus(504)
				return
			}
			ctx.SetStatus(resq.StatusCode)
			body, _ := io.ReadAll(resq.Body)
			ctx.Write(body)
			for key, value := range resq.Header {
				for _, v := range value {
					ctx.AddHeader(key, v)
				}
			}
		}

	})
}

func (server *Server) AddStaticFile(path string, DirPath string) {
	_ = server.AddPath(tool.UrlFormat(path), ANY, func(ctx *RequestCtx) {
		ctx.FastCtx.SendFile(DirPath + tool.UrlFormat(string(ctx.FastCtx.Path()))[len(tool.UrlFormat(path))-1:])
	})
}
