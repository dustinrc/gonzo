package client

import (
	"fmt"
	zmq "github.com/alecthomas/gozmq"
	"github.com/dustinrc/gonzo/mdp"
)

type connection struct {
	ctx  zmq.Context
	sock zmq.Socket
	tmo  float64
}

type timeoutError struct {
	msg string
}

func (e timeoutError) Error() string { return fmt.Sprintf("%v", e.msg) }

func newConnection(url string, timeout float64) (*connection, error) {
	ctx, err := zmq.NewContext()
	if err != nil {
		return nil, err
	}

	sock, err := ctx.NewSocket(zmq.REQ)
	if err != nil {
		return nil, err
	}

	sock.Connect(url)
	conn := connection{ctx, sock, timeout}

	return &conn, nil
}

func (conn *connection) close() {
	conn.sock.Close()
	conn.ctx.Close()
}

func (conn *connection) send(message mdp.Message) (err error) {
	pi := zmq.PollItem{Socket: conn.sock, Events: zmq.POLLOUT}
	pis := zmq.PollItems{pi}
	_, err = zmq.Poll(pis, int64(conn.tmo * 1e6))
	if err != nil {
	} else if i := pis[0]; i.REvents&zmq.POLLOUT != 0 {
		err = conn.sock.SendMultipart(message, 0)
	} else {
		err = timeoutError{"connection.send() timeout"}
	}
	return
}

func (conn *connection) recv() (message mdp.Message, err error) {
	pi := zmq.PollItem{Socket: conn.sock, Events: zmq.POLLIN}
	pis := zmq.PollItems{pi}
	_, err = zmq.Poll(pis, int64(conn.tmo * 1e6))
	if err != nil {
	} else if i := pis[0]; i.REvents&zmq.POLLIN != 0 {
		message, err = conn.sock.RecvMultipart(0)
	} else {
		err = timeoutError{"connection.recv() timeout"}
	}
	return
}
