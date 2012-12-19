package client

import (
	"fmt"
	"github.com/dustinrc/gonzo/mdp"
)

type Client interface {
	Dial() error
	Send(service string, message mdp.Message) (mdp.Message, error)
	Close()
}

type client struct {
	conn     *connection
	url      string
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
	newClient := client{url: brokerURL, timeout: timeout, attempts: attempts}
	return &newClient, nil
}

func (c *client) Dial() error {
	if c.conn != nil {
		return nil
	}

	conn, err := newConnection(c.url)
	if err != nil {
		return err
	}

	c.conn = conn
	return nil
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
