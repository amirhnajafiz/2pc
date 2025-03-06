package packets

// list of packet types
const (
	PktDatabaseRequest int = iota + 1
	PktDatabasePrepare
	PktDatabaseCommit
	PktDatabaseAbort
	PktDatabaseBlock
	PktDatabaseUnblock
)

const (
	PktPaxosRequest int = iota + 100
	PktPaxosPrepare
	PktPaxosAccept
	PktPaxosAccepted
	PktPaxosCommit
	PktPaxosPing
	PktPaxosPong
	PktPaxosSync
)
