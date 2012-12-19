package main

import (
	"flag"
	"fmt"
	"github.com/dustinrc/gonzo/mdp"
	"github.com/dustinrc/gonzo/mdp/client"
	"os"
)

var (
	broker   = flag.String("b", "tcp://127.0.0.1:5555", "broker connection point")
	service  = flag.String("s", "echo", "service requested")
	timeout  = flag.Float64("t", 5.0, "request/reply timeout (seconds)")
	attempts = flag.Int("a", 3, "attempts before failing")
)

func main() {
	flag.Parse()

	c, err := client.New(*broker, *timeout, *attempts)
	if err != nil {
		panic(err)
	}
	c.Dial()
	defer c.Close()

	var argsAsBytes [][]byte
	for _, v := range flag.Args() {
		argsAsBytes = append(argsAsBytes, []byte(v))
	}
	m := mdp.CreateMessage(argsAsBytes...)

	reply, err := c.Send(*service, m)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	} else {
		for i, v := range reply {
			fmt.Printf("frame[%02d]: %s\n", i, v)
		}
	}
}
