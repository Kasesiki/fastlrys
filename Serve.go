package fastlrys

import (
	"encoding/json"
	"mime/multipart"
	"os"

	tool "github.com/Kasesiki/autotool"

	"github.com/valyala/fasthttp"
)

// 将map类型转换为json发送
func (ctx *RequestCtx) JSON(m map[string]interface{}) { ctx.String(200, tool.Mtoj(m)) }

func (ctx *RequestCtx) ReturnStruct(sct any) {
	m := make(map[string]interface{})
	data, _ := json.Marshal(sct)
	json.Unmarshal(data, &m)
	ctx.JSON(m)
}

// 设置返回code
func (ctx *RequestCtx) SetStatus(Code int) { ctx.FastCtx.SetStatusCode(Code) }

// 添加返回header
func (ctx *RequestCtx) AddHeader(key, value string) { ctx.FastCtx.Response.Header.Add(key, value) }

// 关闭连接
func (ctx *RequestCtx) Close() { _ = ctx.FastCtx.Conn().Close() }

// 删除返回header
func (ctx *RequestCtx) DelHeader(key string) { ctx.FastCtx.Response.Header.Del(key) }

// 返回文件
func (ctx *RequestCtx) OriFile(Path string) { fasthttp.ServeFile(ctx.FastCtx, Path) }

// 向返回写入字节
func (ctx *RequestCtx) Write(b []byte) { _, _ = ctx.FastCtx.Write(b) }

// 写入字符串
func (ctx *RequestCtx) WriteString(s string) { _, _ = ctx.FastCtx.WriteString(s) }

// 返回get参数
func (ctx *RequestCtx) QureyString() []byte { return ctx.FastCtx.QueryArgs().QueryString() }

// 返回post参数
func (ctx *RequestCtx) PostArgString() []byte { return ctx.FastCtx.PostArgs().QueryString() }

// 获取请求路径 /path/xxx
func (ctx *RequestCtx) Path() string {
	return string(ctx.FastCtx.Path())
}

// 获取请求body
func (ctx *RequestCtx) Body() []byte { return ctx.FastCtx.Request.Body() }

// 写入字符串和返回code
func (ctx *RequestCtx) String(status int, obj any) {
	_, _ = ctx.FastCtx.WriteString(tool.IntoS(obj))
	ctx.SetStatus(status)
}

func (ctx *RequestCtx) GetFormFile(key string) (*multipart.FileHeader, error) {
	return ctx.FastCtx.FormFile(key)
}

// 获取提交的参数
func (ctx *RequestCtx) GetFormValue(key string) []byte {
	return ctx.FastCtx.FormValue(key)
}

// 返回文件，指向目录时返回目录文件列表
func (ctx *RequestCtx) File(Path string) {
	//ctx.AddHeader("Content-Disposition", "attachment; filename="+tool.GetPathName(Path))
	fasthttp.ServeFile(ctx.FastCtx, Path)
}

// 将给定路径中的本地文件内容作为响应体发送
func (ctx *RequestCtx) SendFile(Path string) {
	ctx.FastCtx.SendFile(Path)
}

func (ctx *RequestCtx) Query(str string) string {
	return string(ctx.FastCtx.QueryArgs().Peek(str))
}

// 将Path指代的html返回，如无法读取则返回504
func (ctx *RequestCtx) HTML(Path string) {
	o, e := os.Open(Path)
	if e != nil {
		ctx.SetStatus(504)
		return
	}
	_, _ = ctx.FastCtx.WriteString(tool.ReadAllString(o))

	switch tool.FileNameGetSuffix(o.Name()) {
	case "html":
		ctx.FastCtx.Response.Header.Set("Content-Type", "text/html; charset=utf-8")
	case "css":
		ctx.FastCtx.Response.Header.Set("Content-Type", "text/css; charset=utf-8")
	}

	ctx.SetStatus(200)
}

// 返回一个数组，内容为json
func (ctx *RequestCtx) JSONList(l []map[string]interface{}) {
	ctx.WriteString("[")
	for i := range l {
		ctx.JSON(l[i])
		if i+1 != len(l) {
			ctx.WriteString(",")
		}
	}
	ctx.WriteString("]")
}
