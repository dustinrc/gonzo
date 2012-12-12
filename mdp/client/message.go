package client

type Message [][]byte

func CreateMessage(frame []byte) Message {
	return [][]byte{frame}
}

func (m Message) AddFrame(frame []byte) Message {
	m = append(m, frame)
	return m
}
