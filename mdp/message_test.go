package mdp

import "testing"

var (
	in     = []byte("foobar")
	out    = [][]byte{[]byte("foobar")}
	frame1 = []byte("raboof")
	frame2 = []byte("")
	frame3 = []byte{0x0a, 0xa0}
	newMsg = [][]byte{[]byte("foobar"), []byte("raboof"),
		[]byte(""), []byte{0x0a, 0xa0}}
	newRevMsg = [][]byte{[]byte{0x0a, 0xa0}, []byte(""),
		[]byte("raboof"), []byte("foobar")}
)

func TestCreateMessage(t *testing.T) {
	x := CreateMessage(in)
	for i, v := range x {
		if string(v) != string(out[i]) {
			t.Errorf("CreateMessage(%v) = %v, want %v", in, x, out)
		}
	}
}

func TestCreateMessageFromMultiple(t *testing.T) {
	x := CreateMessage(in, frame1, frame2, frame3)
	for i, v := range x {
		if string(v) != string(newMsg[i]) {
			t.Errorf("CreateMessage(%v) = %v, want %v", in, x, newMsg)
		}
	}
}

func TestAppend(t *testing.T) {
	origMsg := CreateMessage(in)
	x := origMsg.Append(frame1, frame2, frame3)
	for i, v := range x {
		if string(v) != string(newMsg[i]) {
			t.Error("Append(%v) = %v, want %v", origMsg, x, newMsg)
			break
		}
	}
}

func TestPrepend(t *testing.T) {
	origMsg := CreateMessage(in)
	x := origMsg.Prepend(frame3, frame2, frame1)
	for i, v := range x {
		if string(v) != string(newRevMsg[i]) {
			t.Error("Prepend(%v) = %v, want %v", origMsg, x, newRevMsg)
			break
		}
	}
}
