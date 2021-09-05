package Helper

type Packet struct {
	Type string
	Msg  string
	MessagesSent map[string]bool
}

