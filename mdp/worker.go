package mdp

import (
	"fmt"
	zmq "github.com/alecthomas/gozmq"
	"github.com/dustinrc/gonzo"
	"time"
)

type RequestHandler func(gonzo.Message) gonzo.Message

type Worker interface {
	Dial() error
	Listen(RequestHandler)
	Close()
}

type worker struct {
	conn    *gonzo.Connection
	url     string
	service string
}

func NewWorker(brokerURL string, service string) (Worker, error) {
	newWorker := worker{url: brokerURL, service: service}
	return &newWorker, nil
}

func (w *worker) Dial() error {
	if w.conn != nil {
		return nil
	}

	conn, err := gonzo.NewConnection(w.url, zmq.DEALER)
	if err != nil {
		return err
	}

	w.conn = conn
	return nil
}

func (w *worker) Listen(rh RequestHandler) {
	rq := make(chan gonzo.Message, 1)
	missed := 0
	go w.listen(rq)
	w.ready()
	for {
		select {
		case m := <-rq:
			switch m[2][0] {
			default:
				w.reconnect()
			case HEARTBEAT:
				fmt.Println("HEARTBEAT Received")
				w.heartbeat()
			case REQUEST:
				fmt.Println("REQUEST Received")
				replyBody := rh(m[5:])
				w.reply(replyBody, m[3])
			}
		case <-time.After(3 * time.Second):
			w.heartbeat()
			if missed++; missed >= 3 {
				w.reconnect()
				missed = 0
			}
		}
	}
}

func (w *worker) Close() {
	w.disconnect()
	w.conn.Close()
	w.conn = nil
}

func CreateWorkerMessage(command byte) gonzo.Message {
	return gonzo.CreateMessage([]byte(""), []byte(WV01), []byte{command})
}

func (w *worker) ready() {
	m := CreateWorkerMessage(READY)
	m = m.Append([]byte(w.service))
	w.conn.Send(m, 0.0)
}

func (w *worker) reply(replyBody gonzo.Message, addr []byte) {
	m := CreateWorkerMessage(REPLY)
	m = m.Append(addr, []byte(""))
	m = m.Append(replyBody...)
	w.conn.Send(m, 0.0)
}

func (w *worker) listen(requests chan gonzo.Message) {
	for {
		m, _ := w.conn.Recv(-1)
		requests <- m
	}
}

func (w *worker) heartbeat() {
	m := CreateWorkerMessage(HEARTBEAT)
	w.conn.Send(m, 0.0)
	fmt.Println("HEARTBEAT Sent")
}

func (w *worker) disconnect() {
	m := CreateWorkerMessage(DISCONNECT)
	w.conn.Send(m, 0.0)
}

func (w *worker) reconnect() {
	fmt.Println("disconnect and reconnect...")
	w.disconnect()
	if err := w.Dial(); err != nil {
		panic(err)
	}
	w.ready()
}
