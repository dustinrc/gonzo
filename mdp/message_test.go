package mdp

import "testing"

func TestCreateMessage(t *testing.T) {
	var (
		in  = []byte("foobar")
		out = [][]byte{[]byte("foobar")}
	)
	x := CreateMessage(in)
	for i, v := range x {
		if string(v) != string(out[i]) {
			t.Errorf("CreateMessage(%v) = %v, want %v", in, x, out)
		}
	}
}

func TestAppendFrame(t *testing.T) {
	var (
		origMsg     = CreateMessage([]byte("foobar"))
		newFrame    = []byte("raboof")
		newerFrame  = []byte("")
		newestFrame = []byte{0x0a, 0xa0}
		newMsg      = [][]byte{[]byte("foobar"), []byte("raboof"),
			[]byte(""), []byte{0x0a, 0xa0}}
	)
	x := origMsg.AppendFrame(newFrame)
	x = x.AppendFrame(newerFrame)
	x = x.AppendFrame(newestFrame)
	for i, v := range x {
		if string(v) != string(newMsg[i]) {
			t.Error("AppendFrame(%v) = %v, want %v", origMsg, x, newMsg)
		}
	}
}

func TestPrependFrame(t *testing.T) {
	var (
		origMsg     = CreateMessage([]byte("foobar"))
		newFrame    = []byte("raboof")
		newerFrame  = []byte("")
		newestFrame = []byte{0x0a, 0xa0}
		newMsg      = [][]byte{[]byte{0x0a, 0xa0}, []byte(""),
			[]byte("raboof"), []byte("foobar")}
	)
	x := origMsg.PrependFrame(newFrame)
	x = x.PrependFrame(newerFrame)
	x = x.PrependFrame(newestFrame)
	for i, v := range x {
		if string(v) != string(newMsg[i]) {
			t.Error("PrependFrame(%v) = %v, want %v", origMsg, x, newMsg)
		}
	}
}
