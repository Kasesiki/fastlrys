package main

import "github.com/Kasesiki/fastlrys"

func main() {
	Server := fastlrys.New(80)
	Server.AddGET("/", func(ctx *fastlrys.RequestCtx) {
		ctx.HTML("index.html")
	})
	_ = Server.Do()
}
