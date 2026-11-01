# fastlrys

`fastlrys` 是一个基于 `fasthttp` 的轻量级 Web 框架，强调接口简洁、性能可用和常见能力内置。
English version: [`README.md`](README.md)

## 项目特点
- 基于 `github.com/valyala/fasthttp`，性能开销较低。
- 路由注册简单：`AddGET`、`AddPOST`、`AddAny`、`AddPath`。
- 支持 `*` 后缀通配路由（如 `/api/*`）。
- 支持全局中间件链（`AddMiddle`）。
- 内置按 IP 限流（默认 `24 req/s`，突发 `48`）。
- 支持 TLS，同时可配置为 HTTP 与 HTTPS 并行监听。
- 内置静态文件映射与简易反向代理能力。
- `RequestCtx` 提供常用响应与请求读取方法。

## 环境要求
- Go `1.26rc1` 或更高版本（以 `go.mod` 为准）。

## 快速开始
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

## 核心 API
- `New(port int) *Server`：创建服务实例。
- `Do() error`：启动服务。
- `SetAddr(addr string)`：设置监听地址前缀，默认 `0.0.0.0:`。
- `ListenIPV6(enabled bool)`：开启或关闭 IPv6 监听（默认开启）。
- `SetMaxBodySize(bytes int)`：设置请求体大小上限（默认 `256MB`）。
- `SetTLS(tlsPort int, pemPath, keyPath string)`：启用 HTTPS。
- `SetShowHttp(show bool)`：启用 TLS 时是否保留 HTTP 监听。
- `SetNoRequest(func(*RequestCtx))`：未命中路由时的兜底处理。

## 路由能力
- `AddGET(path, handler)`
- `AddPOST(path, handler)`
- `AddAny(path, handler)`
- `AddPath(path, method, handler)`

通配路由示例：
```go
s.AddAny("/api/*", func(ctx *fastlrys.RequestCtx) {
    ctx.String(200, ctx.Path())
})
```

## 中间件
注册方式：
```go
s.AddMiddle("name", func(ctx *fastlrys.RequestCtx) { /* ... */ })
```

内置中间件：
- `RequestLog`：打印请求日志。
- `NormalStatus`：默认状态码设置为 `200`。
- `NormalCORS`：设置宽松 CORS，并处理 `OPTIONS` 请求。

## 静态资源与反向代理
- `AddStaticFile("/static/*", "./public")`：把 URL 映射到本地文件。
- `AddRemotePath("/proxy/*", "https://example.com")`：把请求转发到上游。

## RequestCtx 常用方法
- 输出相关：`JSON`、`JSONList`、`String`、`Write`、`WriteString`
- 文件相关：`HTML`、`File`、`SendFile`、`OriFile`
- Header/状态：`SetStatus`、`AddHeader`、`DelHeader`、`Close`
- 请求读取：`Body`、`Path`、`Query`、`QureyString`、`PostArgString`
- 表单与结构体：`GetFormFile`、`GetFormValue`、`ReturnStruct`

## QoS 限流
框架在主处理流程中对每个 IP 做限流，可通过 `server.QOS` 调整：
- `QosNumber`：每秒令牌生成速率。
- `QosBursts`：突发容量。
- `QosBucket`：按 IP 维护的令牌桶。

## 说明
- 路由内部基于 Go map 存储；如果定义了重叠规则，匹配顺序可能不稳定。
- 目前 `AddRemotePath` 转发时固定使用 `GET` 方法。
