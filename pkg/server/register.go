//TODO make it a key value storage

package server

import (
	"log"
	"net/rpc"
	"sync"

	"github.com/mateusbraga/static-quorum-system/pkg/view"
)

var register Value

//  ---------- RPC Requests -------------

type RegisterService int

func (r *RegisterService) Read(clientViewRef view.ViewRef, reply *Value) error {
	register.mu.RLock()
	defer register.mu.RUnlock()

	reply.Value = register.Value
	reply.Timestamp = register.Timestamp

	return nil
}

func (r *RegisterService) Write(value Value, reply *Value) error {
	register.mu.Lock()
	defer register.mu.Unlock()

	// Two writes with the same timestamp -> give preference to first one. This makes the Write operation idempotent and still read/write coherent.
	if value.Timestamp > register.Timestamp {
		register.Value = value.Value
		register.Timestamp = value.Timestamp
	}

	return nil
}

func (r *RegisterService) GetCurrentView(value int, reply *view.View) error {
	*reply = *currentView.View()
	log.Println("Done GetCurrentView request")
	return nil
}

// --------- Init ---------

func init() {
	register.mu.Lock() // The register starts locked
	register.Value = nil
	register.Timestamp = 0
}

func init() {
	rpc.Register(new(RegisterService))
}

// --------- Types ---------

type Value struct {
	Value     interface{}
	Timestamp int
	mu        sync.RWMutex
}
