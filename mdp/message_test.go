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

func TestAppendFrame(t *testing.T) {
	origMsg := CreateMessage(in)
	x := origMsg.AppendFrame(frame1)
	x = x.AppendFrame(frame2)
	x = x.AppendFrame(frame3)
	for i, v := range x {
		if string(v) != string(newMsg[i]) {
			t.Error("AppendFrame(%v) = %v, want %v", origMsg, x, newMsg)
		}
	}
}

func TestAppendFrames(t *testing.T) {
	origMsg := CreateMessage(in)
	x := origMsg.AppendFrames(frame1, frame2, frame3)
	for i, v := range x {
		if string(v) != string(newMsg[i]) {
			t.Error("AppendFrames(%v) = %v, want %v", origMsg, x, newMsg)
			break
		}
	}
}

func TestPrependFrame(t *testing.T) {
	origMsg := CreateMessage(in)
	x := origMsg.PrependFrame(frame1)
	x = x.PrependFrame(frame2)
	x = x.PrependFrame(frame3)
	for i, v := range x {
		if string(v) != string(newRevMsg[i]) {
			t.Error("PrependFrame(%v) = %v, want %v", origMsg, x, newRevMsg)
		}
	}
}

func TestPrependFrames(t *testing.T) {
	origMsg := CreateMessage(in)
	x := origMsg.PrependFrames(frame3, frame2, frame1)
	for i, v := range x {
		if string(v) != string(newRevMsg[i]) {
			t.Error("PrependFrames(%v) = %v, want %v", origMsg, x, newRevMsg)
			break
		}
	}
}
