package ws

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func NewDefaultUpGrader(c *gin.Context) (*websocket.Conn, error) {
	upGrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return nil, err
	}

	return ws, nil
}

type WsConn struct {
	conn            *websocket.Conn
	secWebsocketKey string

	listenHeartBeatTimeout time.Duration
	messageListener        func(*WsConnReadMessage)
	closeListener          func()
	isServing              bool
}

func NewWsConn(c *gin.Context) (r *WsConn, resultErr error) {
	ws, err := NewDefaultUpGrader(c)
	if err != nil {
		resultErr = err
		return
	}

	r = &WsConn{
		conn:            ws,
		secWebsocketKey: c.Request.Header.Get("Sec-Websocket-Key"),
	}

	r.setDefault()

	return
}

func (w *WsConn) setDefault() {
	w.isServing = false
	w.SetListenHeartBeatTimeout(0)
	w.SetMessageListener(nil)
	w.SetCloseListener(nil)
}

func (w *WsConn) GetSecWebsocketKey() string {
	return w.secWebsocketKey
}

func (w *WsConn) SetListenHeartBeatTimeout(timeout time.Duration) {
	if timeout == 0 {
		timeout = time.Minute * 5
	}
	w.listenHeartBeatTimeout = timeout
}

func (w *WsConn) SetMessageListener(messageListener func(*WsConnReadMessage)) {
	if messageListener == nil {
		messageListener = func(wcrm *WsConnReadMessage) {}
	}
	w.messageListener = messageListener
}

func (w *WsConn) SetCloseListener(listener func()) {
	if listener == nil {
		listener = func() {}
	}
	w.closeListener = listener
}

func (w *WsConn) Send(
	messageType int, p []byte,
) (resultErr error) {
	return w.conn.WriteMessage(messageType, p)
}

func (w *WsConn) Serve() {
	if w.isServing {
		return
	}
	w.isServing = true

	wg := sync.WaitGroup{}
	wg.Add(2)

	listenChan := make(chan *WsConnReadMessage, 1)

	go func(listenChan chan *WsConnReadMessage) {
		defer wg.Done()
		w.listen(listenChan)
	}(listenChan)

	go func(listenChan chan *WsConnReadMessage) {
		defer wg.Done()

	LOOP:
		for {
			select {
			case m := <-listenChan:
				w.messageListener(m)
			case <-time.After(w.listenHeartBeatTimeout):
				w.Close()
				break LOOP
			}
		}

	}(listenChan)

	go func(listenChan chan *WsConnReadMessage) {
		defer close(listenChan)
		wg.Wait()
	}(listenChan)
}

func (w WsConn) listen(listenChan chan *WsConnReadMessage) {
	for {
		messageType, p, err := w.conn.ReadMessage()

		msg := NewWsConnReadMessage(messageType, p, err)
		listenChan <- msg

		if err != nil {
			break
		}
	}
}

func (w WsConn) Close() (resultErr error) {
	if err := w.conn.Close(); err != nil {
		resultErr = err
		return
	}
	w.closeListener()
	return
}

type WsConnReadMessage struct {
	MessageType int
	P           []byte
	Err         error
}

func NewWsConnReadMessage(
	messageType int, p []byte, err error,
) *WsConnReadMessage {
	return &WsConnReadMessage{
		MessageType: messageType,
		P:           p,
		Err:         err,
	}
}
