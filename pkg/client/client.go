/*
Package client implements a static-quorum-system client
*/
package client

import (
	"log"

	"github.com/mateusbraga/static-quorum-system/pkg/view"
)

type Client struct {
	view *view.View
}

// New returns a new Client with initialView.
func New(initialView *view.View) *Client {
	return &Client{initialView}
}

func (thisClient Client) View() *view.View {
	return thisClient.view
}

// Write v to the system's register.
func (thisClient *Client) Write(v interface{}) error {
	readValue, err := readQuorum(thisClient.View())
	if err != nil {
		// Special cases:
		//  diffResultsErr: can be ignored
		if err == diffResultsErr {
			// Do nothing - we will write a new value anyway
		} else {
			return err
		}
	}

	writeMsg := RegisterMsg{}
	writeMsg.Value = v
	writeMsg.Timestamp = readValue.Timestamp + 1

	err = writeQuorum(thisClient.View(), writeMsg)
	if err != nil {
		return err
	}

	return nil
}

// Read executes the quorum read protocol.
func (thisClient *Client) Read() (interface{}, error) {
	readMsg, err := readQuorum(thisClient.View())
	if err != nil {
		// Expected: diffResultsErr (will write most current value to view).
		if err == diffResultsErr {
			log.Println("Found divergence: Going to 2nd phase of read protocol")
			return thisClient.read2ndPhase(thisClient.View(), readMsg)
		} else {
			return 0, err
		}
	}

	return readMsg.Value, nil
}

func (thisClient *Client) read2ndPhase(destinationView *view.View, readMsg RegisterMsg) (interface{}, error) {
	err := writeQuorum(destinationView, readMsg)
	if err != nil {
		return 0, err
	}
	return readMsg.Value, nil
}

type RegisterMsg struct {
	Value     interface{} // Value of the register
	Timestamp int         // Timestamp of the register

	Err error
}
