package server

import (
	"net/rpc"
	"sync"

	"github.com/mateusbraga/static-quorum-system/pkg/view"
)

var register Value

//  ---------- RPC Requests -------------
type ClientRequest int

func (r *ClientRequest) Read(clientView *view.View, reply *Value) error {
	register.mu.RLock()
	defer register.mu.RUnlock()

	reply.Value = register.Value
	reply.Timestamp = register.Timestamp

	return nil
}

func (r *ClientRequest) Write(value Value, reply *Value) error {
	register.mu.Lock()
	defer register.mu.Unlock()

	// Two writes with the same timestamp -> give preference to first one. This makes the Write operation idempotent and still read/write coherent.
	if value.Timestamp > register.Timestamp {
		register.Value = value.Value
		register.Timestamp = value.Timestamp
	}

	return nil
}

// --------- Init ---------
func init() {
	register.Value = nil
	register.Timestamp = 0
}

func init() {
	rpc.Register(new(ClientRequest))
}

// --------- Types ---------
type Value struct {
	Value     interface{}
	Timestamp int

	mu sync.RWMutex
}
