package main

import "fastlrys"

func main() {
	Server := fastlrys.New(80)
	Server.AddGET("/", func(ctx *fastlrys.RequestCtx) {
		ctx.HTML("index.html")
	})
	_ = Server.Do()
}
