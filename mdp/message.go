package mdp

type Message [][]byte

func CreateMessage(frame []byte) Message {
	return [][]byte{frame}
}

func (m Message) AppendFrame(frame []byte) Message {
	m = append(m, frame)
	return m
}

func (m Message) AppendFrames(frames... []byte) Message {
	m = append(m, frames...)
	return m
}

func (m Message) PrependFrame(frame []byte) Message {
	m = append([][]byte{frame}, m...)
	return m
}

func (m Message) PrependFrames(frames... []byte) Message {
	m = append(frames, m...)
	return m
}
