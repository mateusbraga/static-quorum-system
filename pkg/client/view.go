package client

import (
	"errors"
	"fmt"
	"github.com/mateusbraga/static-quorum-system/pkg/comm"
	"github.com/mateusbraga/static-quorum-system/pkg/view"
	"log"
)

// GetCurrentView asks processes for the its current view and returns it.
func GetCurrentView(processes ...view.Process) (*view.View, error) {
	for _, loopProcess := range processes {
		var receivedView *view.View
		err := comm.SendRPCRequest(loopProcess, "RegisterService.GetCurrentView", 0, &receivedView)
		if err != nil {
			log.Println(err)
			continue
		}

		return receivedView, nil
	}
	return nil, errors.New(fmt.Sprintf("Failed to get current view from any of the processes: %v", processes))
}
