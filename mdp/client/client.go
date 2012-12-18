package client

import (
	"fmt"
	"github.com/dustinrc/gonzo/mdp"
)

type Client interface {
	Send(service string, message mdp.Message) (mdp.Message, error)
	Close()
}

type client struct {
	conn *connection
}

type clientMismatchError struct {
	mismatchType   string
	want, received interface{}
}

func (e clientMismatchError) Error() string {
	return fmt.Sprintf("%v mismatch: want \"%v\", received \"%v\"", e.mismatchType,
		e.want, e.received)
}

func New(brokerURL string, timeout float64) (Client, error) {
	conn, err := newConnection(brokerURL, timeout)
	if err != nil {
		return nil, err
	}

	newClient := client{conn}

	return &newClient, nil
}

func (c *client) Close() {
	c.conn.close()
}

func (c *client) Send(service string, message mdp.Message) (mdp.Message, error) {
	request := message.Prepend([]byte(mdp.CV01), []byte(service))
	err := c.conn.send(request)
	if err != nil {
		return nil, err
	}

	reply, err := c.conn.recv()
	if err != nil {
		return nil, err
	}

	if proto := string(reply[0]); proto != mdp.CV01 {
		err = clientMismatchError{"protocol", mdp.CV01, proto}
		return nil, err
	} else if srvc := string(reply[1]); srvc != service {
		err = clientMismatchError{"service", service, srvc}
		return nil, err
	}
	reply = reply[2:]
	return reply, nil
}
