package client

import "github.com/dustinrc/gonzo/mdp"

type Client interface {
	Send(service string, message mdp.Message) (mdp.Message, error)
	Close()
}

type client struct {
	conn *connection
}

func New(brokerURL string) (Client, error) {
	conn, err := newConnection(brokerURL)
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
	request := message.PrependFrames([]byte(mdp.CV01), []byte(service))
	err := c.conn.send(request)
	if err != nil {
		return nil, err
	}

	reply, err := c.conn.recv()
	if err != nil {
		return nil, err
	}

	reply = reply[2:]
	return reply, nil
}
