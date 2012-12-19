package client

import (
	"fmt"
	zmq "github.com/alecthomas/gozmq"
	"github.com/dustinrc/gonzo"
)

type connection struct {
	ctx  zmq.Context
	sock zmq.Socket
}

type timeoutError struct {
	msg string
}

func (e timeoutError) Error() string { return fmt.Sprintf("%v", e.msg) }

func newConnection(url string) (*connection, error) {
	ctx, err := zmq.NewContext()
	if err != nil {
		return nil, err
	}

	sock, err := ctx.NewSocket(zmq.REQ)
	if err != nil {
		return nil, err
	}

	sock.Connect(url)
	conn := connection{ctx, sock}

	return &conn, nil
}

func (conn *connection) close() {
	conn.sock.Close()
	conn.sock = nil
	conn.ctx.Close()
	conn.ctx = nil
}

func (conn *connection) send(message gonzo.Message, timeout float64) (err error) {
	pi := zmq.PollItem{Socket: conn.sock, Events: zmq.POLLOUT}
	pis := zmq.PollItems{pi}
	_, err = zmq.Poll(pis, int64(timeout*1e6))
	if err != nil {
	} else if i := pis[0]; i.REvents&zmq.POLLOUT != 0 {
		err = conn.sock.SendMultipart(message, 0)
	} else {
		err = timeoutError{"connection.send() timeout"}
	}
	return
}

func (conn *connection) recv(timeout float64) (message gonzo.Message, err error) {
	pi := zmq.PollItem{Socket: conn.sock, Events: zmq.POLLIN}
	pis := zmq.PollItems{pi}
	_, err = zmq.Poll(pis, int64(timeout*1e6))
	if err != nil {
	} else if i := pis[0]; i.REvents&zmq.POLLIN != 0 {
		message, err = conn.sock.RecvMultipart(0)
	} else {
		err = timeoutError{"connection.recv() timeout"}
	}
	return
}
