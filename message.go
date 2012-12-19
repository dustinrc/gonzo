package gonzo

type Message [][]byte

func CreateMessage(frames ...[]byte) Message {
	return append([][]byte{}, frames...)
}

func (m Message) Append(frames ...[]byte) Message {
	return append(m, frames...)
}

func (m Message) Prepend(frames ...[]byte) Message {
	return append(frames, m...)
}
