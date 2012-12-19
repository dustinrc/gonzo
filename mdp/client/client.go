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
	conn     *connection
	timeout  float64
	attempts int
}

type clientMismatchError struct {
	mismatchType   string
	want, received interface{}
}

func (e clientMismatchError) Error() string {
	return fmt.Sprintf("%v mismatch: want \"%v\", received \"%v\"", e.mismatchType,
		e.want, e.received)
}

func New(brokerURL string, timeout float64, attempts int) (Client, error) {
	conn, err := newConnection(brokerURL)
	if err != nil {
		return nil, err
	}

	newClient := client{conn, timeout, attempts}

	return &newClient, nil
}

func (c *client) Close() {
	c.conn.close()
}

func (c *client) Send(service string, message mdp.Message) (mdp.Message, error) {
	request := message.Prepend([]byte(mdp.CV01), []byte(service))
	err := c.conn.send(request, c.timeout)
	if err != nil {
		return nil, err
	}

	var reply mdp.Message
	for attempt := 1; attempt <= c.attempts; attempt++ {
		reply, err = c.conn.recv(c.timeout)
		if err == nil {
			break
		}
		if err != nil {
			fmt.Println("Failed attempt", attempt, "of", c.attempts)
		}
		if attempt == c.attempts {
			return nil, err
		}
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
