package fastlrys

import (
	"errors"

	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
)

var Normalupgrader = websocket.FastHTTPUpgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(ctx *fasthttp.RequestCtx) bool {
		return true
	},
}

const CloseMessage int = 8

type WebSocketConn struct {
	conn *websocket.Conn
}
type WebSocketHandler func(c *WebSocketConn)

func CtxConvWs(ctx *RequestCtx, handler WebSocketHandler) error {
	return Normalupgrader.Upgrade(ctx.FastCtx, func(c *websocket.Conn) {
		defer c.Close()
		handler(&WebSocketConn{conn: c})
	})
}

func (c *WebSocketConn) Close() error {
	return c.conn.Close()
}

func (c *WebSocketConn) WriteText(text string) error {
	return c.conn.WriteMessage(websocket.TextMessage, []byte(text))
}

func (c *WebSocketConn) WriteJson(v any) error {
	return c.conn.WriteJSON(v)
}

func (c *WebSocketConn) ReadText() (string, error) {
	t, p, e := c.conn.ReadMessage()
	if e != nil {
		return "", e
	}
	if t == websocket.BinaryMessage {
		return "", errors.New("message type not is Text")
	}
	return string(p), nil
}

func (c *WebSocketConn) ReadJson(v any) error {
	return c.conn.ReadJSON(v)
}

func (c *WebSocketConn) Read() (messageType int, p []byte, err error) {
	return c.conn.ReadMessage()
}

// var (
// 	newline = []byte{'\n'}
// 	space   = []byte{' '}
// )

// var upgrader = websocket.FastHTTPUpgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// }

// readRemote pumps messages from the websocket connection to the hub.
//
// The application runs readRemote in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func readRemote(a *websocket.Conn, b *websocket.Conn) {
	go func() {
		for {
			typ, p, e := b.ReadMessage()
			if e != nil {
				return
			}
			if typ == websocket.TextMessage {
				a.WriteMessage(websocket.TextMessage, p)
			} else {
				b.WriteMessage(CloseMessage, []byte{})
			}
		}
	}()
	for {
		typ, p, e := a.ReadMessage()
		if e != nil {
			return
		}
		if typ == websocket.TextMessage {
			b.WriteMessage(websocket.TextMessage, p)
		} else {
			a.WriteMessage(CloseMessage, []byte{})
		}
	}
}
