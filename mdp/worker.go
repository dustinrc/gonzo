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
	Ready()
	Reply(gonzo.Message, []byte)
	Heartbeat()
	Disconnect()
	Close()
}

type worker struct {
	conn *gonzo.Connection
	url string
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

func (w *worker) Close() {
	w.conn.Close()
	w.conn = nil
}

func CreateWorkerMessage(command byte) gonzo.Message {
	return gonzo.CreateMessage([]byte(""), []byte(WV01), []byte{command})
}

func (w *worker) Ready() {
	m := CreateWorkerMessage(READY)
	m = m.Append([]byte(w.service))
	w.conn.Send(m, 0.0)
}

func (w *worker) Reply(replyBody gonzo.Message, addr []byte) {
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

func (w *worker) Listen(rh RequestHandler) {
	rq := make(chan gonzo.Message, 1)
	missed := 0
	go w.listen(rq)
	for {
		select {
		case m := <-rq:
			switch m[2][0] {
			default:
				w.Reconnect()
			case HEARTBEAT:
				fmt.Println("HEARTBEAT Received")
				w.Heartbeat()
			case REQUEST:
				fmt.Println("REQUEST Received")
				replyBody := rh(m[5:])
				w.Reply(replyBody, m[3])
			}
		case <-time.After(3 * time.Second):
			w.Heartbeat()
			if missed++; missed >= 3 {
				w.Reconnect()
				missed = 0
			}
		}
	}
}

func (w *worker) Heartbeat() {
	m := CreateWorkerMessage(HEARTBEAT)
	w.conn.Send(m, 0.0)
	fmt.Println("HEARTBEAT Sent")
}

func (w *worker) Disconnect() {
	m := CreateWorkerMessage(DISCONNECT)
	w.conn.Send(m, 0.0)
}

func (w *worker) Reconnect() {
	fmt.Println("disconnect and reconnect...")
	w.Disconnect()
	if err := w.Dial(); err != nil { panic(err) }
	w.Ready()
}
