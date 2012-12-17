package main

import (
	"flag"
	"fmt"
	"github.com/dustinrc/gonzo/mdp"
	"github.com/dustinrc/gonzo/mdp/client"
)

var (
	broker = flag.String("b", "tcp://127.0.0.1:5555", "broker connection point")
	data = flag.String("d", "ping", "data to send to the service")
	service = flag.String("s", "echo", "service requested")
)

func main() {
	flag.Parse()

	c, err := client.New(*broker)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	m := mdp.CreateMessage([]byte(*data))

	reply, err := c.Send(*service, m)
	if err != nil {
		panic(err)
	}
	for i, v := range reply {
		fmt.Printf("frame[%d]: %s\n", i, v)
	}
}
