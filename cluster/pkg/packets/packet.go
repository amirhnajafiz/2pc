package packets

// Packet is a wrapper for transfering gRPC messages to state machines.
type Packet struct {
	Label   int
	Payload interface{}
}
