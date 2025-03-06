package enums

// list of system's write-ahead log types
const (
	WALStart  = "start"
	WALUpdate = "update"
	WALCommit = "commit"
	WALAbort  = "abort"
)
