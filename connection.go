package gonzo

import (
	"fmt"
	zmq "github.com/alecthomas/gozmq"
)

type Connection struct {
	ctx  zmq.Context
	sock zmq.Socket
}

type timeoutError struct {
	msg string
}

func (e timeoutError) Error() string { return fmt.Sprintf("%v", e.msg) }

func NewConnection(url string) (*Connection, error) {
	ctx, err := zmq.NewContext()
	if err != nil {
		return nil, err
	}

	sock, err := ctx.NewSocket(zmq.REQ)
	if err != nil {
		return nil, err
	}

	sock.Connect(url)
	conn := Connection{ctx, sock}

	return &conn, nil
}

func (conn *Connection) Close() {
	conn.sock.Close()
	conn.sock = nil
	conn.ctx.Close()
	conn.ctx = nil
}

func (conn *Connection) Send(message Message, timeout float64) (err error) {
	pi := zmq.PollItem{Socket: conn.sock, Events: zmq.POLLOUT}
	pis := zmq.PollItems{pi}
	_, err = zmq.Poll(pis, int64(timeout*1e6))
	if err != nil {
	} else if i := pis[0]; i.REvents&zmq.POLLOUT != 0 {
		err = conn.sock.SendMultipart(message, 0)
	} else {
		err = timeoutError{"Connection.Send() timeout"}
	}
	return
}

func (conn *Connection) Recv(timeout float64) (message Message, err error) {
	pi := zmq.PollItem{Socket: conn.sock, Events: zmq.POLLIN}
	pis := zmq.PollItems{pi}
	_, err = zmq.Poll(pis, int64(timeout*1e6))
	if err != nil {
	} else if i := pis[0]; i.REvents&zmq.POLLIN != 0 {
		message, err = conn.sock.RecvMultipart(0)
	} else {
		err = timeoutError{"Connection.Recv() timeout"}
	}
	return
}
