package wstypes

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type ConnWriter struct {
	writeCh chan []byte
	conn    *websocket.Conn
	done    chan struct{}
}

func NewConnWriter(conn *websocket.Conn) *ConnWriter {
	cw := &ConnWriter{
		writeCh: make(chan []byte, 256),
		conn:    conn,
		done:    make(chan struct{}),
	}

	go cw.writeInConn()

	return cw
}

func (cw *ConnWriter) Send(msg []byte) error {
	select {
	case cw.writeCh <- msg:
	default:
		return fmt.Errorf("client is dead")
	}
	return nil
}

func (cw *ConnWriter) writeInConn() {
	defer close(cw.done)
	for msg := range cw.writeCh {
		if err := cw.conn.WriteMessage(websocket.TextMessage, msg); err != nil {

			// TODO: обработка ошибки

			return
		}
	}
}

func (cw *ConnWriter) Close() {
	close(cw.writeCh)
	<-cw.done
	cw.conn.Close()
}
