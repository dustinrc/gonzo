package client

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

func TestAddFrame(t *testing.T) {
	var (
		origMsg     = CreateMessage([]byte("foobar"))
		newFrame    = []byte("raboof")
		newerFrame  = []byte("")
		newestFrame = []byte{0x0a, 0xa0}
		newMsg      = [][]byte{[]byte("foobar"), []byte("raboof"),
			[]byte(""), []byte{0x0a, 0xa0}}
	)
	x := origMsg.AddFrame(newFrame)
	x = x.AddFrame(newerFrame)
	x = x.AddFrame(newestFrame)
	for i, v := range x {
		if string(v) != string(newMsg[i]) {
			t.Error("AddFrame(%v) = %v, want %v", origMsg, x, newMsg)
		}
	}
}
