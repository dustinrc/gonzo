package main

import (
	"flag"
	"fmt"
	"github.com/dustinrc/gonzo"
	"github.com/dustinrc/gonzo/mdp"
	"os"
)

var (
	broker = flag.String("b", "tcp://127.0.0.1:5555", "broker connection point")
	service = flag.String("s", "echo", "service provided")
)

func echo(request gonzo.Message) (reply gonzo.Message) {
	reply = gonzo.CreateMessage(request...)
	reply = reply.Prepend([]byte("echoing..."))
	return
}

func main() {
	flag.Parse()

	w, err := mdp.NewWorker(*broker, *service)
	if err != nil {
		panic(err)
	}
	err = w.Dial()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer w.Close()

	w.Ready()
	defer w.Disconnect()

	w.Listen(echo)
}
