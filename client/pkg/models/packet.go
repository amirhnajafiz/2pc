package models

// Packet is a wrapper for transfering gRPC messages to the processor.
type Packet struct {
	Label   string
	Payload interface{}
}
