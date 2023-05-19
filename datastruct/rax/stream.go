package rax

import (
	"time"
)

// streamID key-value
type streamID struct {
	Key   string //Unix time in milliseconds
	Value string //Sequence number
}

type Stream struct {
	MsgList Rax      //Stream list
	Length  int      //Stream length
	LastID  streamID //Zero if there are yet no items.
	Cgroups Rax      //Consumers group list
}

// StreamCG Consumer group
type StreamCG struct {
	LastID    streamID //Last delivered (not acknowledged) ID for this group.
	Pel       Rax      //Pending entries list.
	Consumers Rax      //Consumer list
}

type StreamConsumer struct {
	Name       string    //Consumer name
	ActiveTime time.Time //Last time this consumer was active.
	Pel        Rax       //Pending entries list.
}
type StreamNACK struct {
	consumer      StreamConsumer //The consumer this message was delivered toin the last delivery.
	DeliveryTime  time.Time      //Last time this message was delivered.
	deliveryCount int            //Number of times this message was delivered
}
