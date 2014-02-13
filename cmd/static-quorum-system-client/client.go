/*
static-quorum-system-client is a sample implementation of a static-quorum-system client.
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/mateusbraga/static-quorum-system/pkg/client"
	"github.com/mateusbraga/static-quorum-system/pkg/view"
)

var (
	nTotal = flag.Uint64("n", math.MaxUint64, "number of times to perform a read and write operation")
)

func main() {
	flag.Parse()

	initialView := getInitialView()
	quorumClient := client.New(initialView)

	var finalValue interface{}
	var err error
	for i := uint64(0); i < *nTotal; i++ {
		startRead := time.Now()
		finalValue, err = quorumClient.Read()
		endRead := time.Now()
		if err != nil {
			log.Fatalln(err)
		}

		startWrite := time.Now()
		err = quorumClient.Write(finalValue)
		endWrite := time.Now()
		if err != nil {
			log.Fatalln(err)
		}

		if i%1000 == 0 {
			fmt.Printf("%v: Read %v (%v)-> Write (%v)\n", i, finalValue, endRead.Sub(startRead), endWrite.Sub(startWrite))
		} else {
			fmt.Printf(".")
		}
	}
}

func getInitialView() *view.View {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalln(err)
	}

	var processes []view.Process
	switch {
	case strings.Contains(hostname, "node-"): // emulab.net
		processes = append(processes, view.Process{"10.1.1.2:5000"})
		processes = append(processes, view.Process{"10.1.1.3:5000"})
		processes = append(processes, view.Process{"10.1.1.4:5000"})
	default:
		processes = append(processes, view.Process{"[::]:5000"})
		processes = append(processes, view.Process{"[::]:5001"})
		processes = append(processes, view.Process{"[::]:5002"})
	}

	return view.NewWithProcesses(processes...)
}
