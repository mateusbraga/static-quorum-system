// static-quorum-systemd runs a sample implementation of a static-quorum-system server.
//
// Most of the work is done at github.com/mateusbraga/static-quorum-system/pkg/server.
package main

import (
	"flag"
	"log"
	"runtime"

	"github.com/mateusbraga/static-quorum-system/pkg/server"

	"net/http"
	_ "net/http/pprof"
)

// Flags
var (
	bindAddr = flag.String("bind", "[::]:5000", "Set this process address")
)

func init() {
	// Make it parallel
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	flag.Parse()

	go func() {
		log.Println("Running pprof:", http.ListenAndServe("localhost:6060", nil))
	}()

	server.Run(*bindAddr)
}
