# fastlrys

A lightweight HTTP framework built on top of `fasthttp`, focused on simple APIs, middleware composition, and pragmatic built-in helpers.
Chinese version: [`README_zh_CN.md`](README_zh_CN.md)

## Highlights
- High-performance serving powered by `github.com/valyala/fasthttp`.
- Simple route registration: `AddGET`, `AddPOST`, `AddAny`, and `AddPath`.
- Wildcard path matching with `*` suffix (for example `/api/*`).
- Global middleware pipeline via `AddMiddle`.
- Built-in per-IP rate limiting (default: `24 req/s`, burst `48`).
- Optional TLS serving, with optional HTTP + HTTPS dual listening.
- Helpers for static file serving and lightweight reverse proxying.
- Convenient `RequestCtx` methods for JSON, HTML, file responses, query/form access, and headers.

## Requirements
- Go `1.26rc1` or newer (based on `go.mod`).

## Quick Start
```go
package main

import "fastlrys"

func main() {
    s := fastlrys.New(8080)

    s.AddMiddle("log", fastlrys.RequestLog)
    s.AddMiddle("cors", fastlrys.NormalCORS)

    s.AddGET("/", func(ctx *fastlrys.RequestCtx) {
        ctx.JSON(map[string]interface{}{"message": "hello"})
    })

    s.AddPOST("/echo", func(ctx *fastlrys.RequestCtx) {
        ctx.Write(ctx.Body())
        ctx.SetStatus(200)
    })

    _ = s.Do()
}
```

## Core API
- `New(port int) *Server`: create server with sane defaults.
- `Do() error`: start serving.
- `SetAddr(addr string)`: set listen address prefix (default `0.0.0.0:`).
- `ListenIPV6(enabled bool)`: toggle IPv6 listen (`true` by default).
- `SetMaxBodySize(bytes int)`: configure max request body size (default `256MB`).
- `SetTLS(tlsPort int, pemPath, keyPath string)`: enable HTTPS listener.
- `SetShowHttp(show bool)`: when TLS is enabled, keep HTTP listener running too.
- `SetNoRequest(func(*RequestCtx))`: fallback handler when no route matches.

## Routing
- `AddGET(path, handler)`
- `AddPOST(path, handler)`
- `AddAny(path, handler)`
- `AddPath(path, method, handler)`

Wildcard example:
```go
s.AddAny("/api/*", func(ctx *fastlrys.RequestCtx) {
    ctx.String(200, ctx.Path())
})
```

## Middleware
Register global middleware with:
```go
s.AddMiddle("name", func(ctx *fastlrys.RequestCtx) { /* ... */ })
```

Built-in middleware:
- `RequestLog`: request logging.
- `NormalStatus`: default status `200`.
- `NormalCORS`: permissive CORS headers + handles `OPTIONS`.

## Static Files and Proxy
- `AddStaticFile("/static/*", "./public")`: map URL path to local files.
- `AddRemotePath("/proxy/*", "https://example.com")`: forward requests upstream.

## Response Helpers (`RequestCtx`)
- `JSON`, `JSONList`, `String`, `Write`, `WriteString`
- `HTML`, `File`, `SendFile`, `OriFile`
- `SetStatus`, `AddHeader`, `DelHeader`, `Close`
- `Body`, `Path`, `Query`, `QureyString`, `PostArgString`
- `GetFormFile`, `GetFormValue`, `ReturnStruct`

## QoS / Rate Limiting
`fastlrys` applies per-IP rate limiting in the main handler.
Tune it directly on `server.QOS`:
- `QosNumber`: token fill rate per second.
- `QosBursts`: burst size.
- `QosBucket`: limiter buckets by client IP.

## Notes
- Route storage uses a Go map; avoid ambiguous overlapping wildcard routes if deterministic match order matters.
- `AddRemotePath` creates upstream requests with method `GET` in current implementation.
