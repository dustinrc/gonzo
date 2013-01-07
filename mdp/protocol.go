package mdp

const (
	CV01 = "MDPC01"
	WV01 = "MDPW01"
)

const (
	_ byte = iota
	READY
	REQUEST
	REPLY
	HEARTBEAT
	DISCONNECT
)
