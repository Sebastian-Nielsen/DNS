package Helper

type Packet struct {
	Type string
	Msg  string
	MessagesSent map[string]bool
	PeersInArrivalOrderValues []string
}

var PacketType = struct {
	BROADCASTED_KNOWN_LISTENER_PORT string
	BROADCAST_LISTENER_PORT string
	LISTENER_PORT string
	BROADCAST_MSG string
	PULL          string
	PULL_REPLY string
} {
	BROADCASTED_KNOWN_LISTENER_PORT: "BROADCASTED_KNOWN_LISTENER_PORT",
	BROADCAST_LISTENER_PORT: "BROADCAST_LISTENER_PORT",
	LISTENER_PORT:           "LISTENER_PORT",
	BROADCAST_MSG:           "BROADCAST_MSG",
	PULL:                    "PULL",
	PULL_REPLY:              "PULL_REPLY",
}
