package client

import (
	zmq "github.com/alecthomas/gozmq"
	"github.com/dustinrc/gonzo/mdp"
)

type connection struct {
	ctx  zmq.Context
	sock zmq.Socket
}

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
	conn.ctx.Close()
}

func (conn *connection) send(message mdp.Message) error {
	err := conn.sock.SendMultipart(message, 0)
	return err
}

func (conn *connection) recv() (mdp.Message, error) {
	reply, err := conn.sock.RecvMultipart(0)
	return reply, err
}
