package packets

// list of packet types
const (
	PktRequest int = iota + 1
	PktPrepare
	PktCommit
	PktAbort
	PktPaxosAccept
	PktPaxosAccepted
	PktPaxosCommit
)
