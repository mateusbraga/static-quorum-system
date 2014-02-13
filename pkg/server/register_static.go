// +build staticRegister

package server

import (
	"net/rpc"
	"sync"

	"github.com/mateusbraga/freestore/pkg/view"
)

var register Value

//  ---------- Interface -------------
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

	if value.Timestamp > register.Timestamp {
		register.Value = value.Value
		register.Timestamp = value.Timestamp
	}

	return nil
}

// --------- Bootstrapping ---------
func init() {
	register.mu.Lock() // The register starts locked
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

	View *view.View
	Err  error

	mu sync.RWMutex
}