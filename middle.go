package fastlrys

import (
	tool "github.com/Kasesiki/autotool"
)

func RequestLog(ctx *RequestCtx) {
	tool.LogPrint(ctx.FastCtx.Request.URI().String()+"  "+string(ctx.FastCtx.Method())+"  "+ctx.FastCtx.RemoteIP().String(), tool.INFO, 1)
}

func NormalStatus(ctx *RequestCtx) {
	ctx.SetStatus(200)
}

func NormalCORS(ctx *RequestCtx) {
	ctx.AddHeader("Access-Control-Allow-Origin", "*")
	ctx.AddHeader("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization, Accept, X-Requested-With , yourHeaderFeild")
	ctx.AddHeader("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
	if ctx.Method == OPTIONS {
		ctx.SetStatus(200)
	}
}

