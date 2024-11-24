package csm

import (
	"github.com/F24-CSE535/2pc/cluster/internal/memory"
	"github.com/F24-CSE535/2pc/cluster/pkg/packets"
)

// dispatcher is a simple go-routine that accepts user requests into a channel and sends
// them to CSMs when its not busy with consensus protocol.
type Dispatcher struct {
	memory *memory.SharedMemory
	input  chan *packets.Packet
	output chan *packets.Packet
	notify chan bool
}

// new dispatcher returns a dispatcher instance.
func NewDispatcher(input, output chan *packets.Packet, mem *memory.SharedMemory) *Dispatcher {
	// input channel is the channel that gRPC level methods publish in
	// output channel is the channel of CSMs
	instance := Dispatcher{
		memory: mem,
		input:  input,
		output: output,
		notify: make(chan bool),
	}

	// start the dispatcher inside a go-routine
	go instance.start()

	return &instance
}

// on start, the dispatcher gets messages from its input channel, publishs them inside output channel
// and waits for a notify signal.
func (d *Dispatcher) start() {
	for {
		// capture packets
		pkt := <-d.input

		// drop the message if node is not leader
		if d.memory.GetNodeName() != d.memory.GetLeader() || d.memory.GetBlockStatus() {
			continue
		}

		// export to the output
		d.output <- pkt

		// wait until the handlers sends a notify response
		<-d.notify
	}
}

// GetNotifyChannel returns the notification channel of dispatcher.
func (d *Dispatcher) GetNotifyChannel() chan bool {
	return d.notify
}
