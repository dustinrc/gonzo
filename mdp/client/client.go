package client

import "github.com/dustinrc/gonzo/mdp"

type Client interface {
	Send(service string, message mdp.Message) error
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

func (c *client) Send(service string, message mdp.Message) error {
	message = message.PrependFrames([]byte(mdp.CV01), []byte(service))
	err := c.conn.send(message)
	return err
}
