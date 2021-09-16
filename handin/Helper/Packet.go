package Helper

type Packet struct {
	Type string
	Msg  string
	MessagesSent map[string]bool
}

var PacketType = struct {
	UPDATE string
	PULL string
	PULL_REPLY string
} {
	UPDATE: "UPDATE",
	PULL: "PULL",
	PULL_REPLY: "PULL_REPLY",
}
