package packets

// list of packet types
const (
	PktDatabaseRequest int = iota + 1
	PktDatabasePrepare
	PktDatabaseCommit
	PktDatabaseAbort
)

const (
	PktPaxosRequest int = iota + 100
	PktPaxosAccept
	PktPaxosAccepted
	PktPaxosCommit
)
