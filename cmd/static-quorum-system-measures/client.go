// Freestore_measures is a client that measures latency and throughput of static-quorum-system.
package main

import (
	"bufio"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/mateusbraga/gostat"
	"github.com/mateusbraga/static-quorum-system/pkg/client"
	"github.com/mateusbraga/static-quorum-system/pkg/view"
)

var (
	isWrite            = flag.Bool("write", false, "Client will measure write operations")
	size               = flag.Int("size", 1, "The size of the data being transfered")
	numberOfOperations = flag.Int("n", 1000, "Number of operations to perform (latency measurement)")
	measureLatency     = flag.Bool("latency", false, "Client will measure latency")
	measureThroughput  = flag.Bool("throughput", false, "Client will measure throughput")
	totalDuration      = flag.Duration("duration", 10*time.Second, "Duration to run operations (throughput measurement)")
	resultFile         = flag.String("o", "/proj/freestore/results_static.txt", "Result file filename")
)

var (
	latencies    []int64
	ops          int
	stopChan     <-chan time.Time
	systemClient *client.Client
)

func init() {
	// Make it parallel
	runtime.GOMAXPROCS(runtime.NumCPU())

	initialView := getInitialView()
	systemClient = client.New(initialView)
}

func main() {
	flag.Parse()

	stopChan = time.After(*totalDuration)

	if *measureLatency {
		latency()
	} else if *measureThroughput {
		throughput()
	} else {
		latencyAndThroughput()
	}
}

func latencyAndThroughput() {
	if *isWrite {
		log.Printf("Measuring throughput and latency of write operations with size %vB\n", *size)
	} else {
		log.Printf("Measuring throughput and latency of read operations with size %vB\n", *size)
	}

	data := createFakeData()

	startTime := time.Now()

	if *isWrite {
		for ops = 0; ops < *numberOfOperations; ops++ {
			timeBefore := time.Now()
			err := systemClient.Write(data)
			timeAfter := time.Now()
			if err != nil {
				log.Fatalln(err)
			}

			latencies = append(latencies, timeAfter.Sub(timeBefore).Nanoseconds())
		}
	} else {
		err := systemClient.Write(data)
		if err != nil {
			log.Fatalln("ERROR initial write:", err)
		}

		for ops = 0; ops < *numberOfOperations; ops++ {
			timeBefore := time.Now()
			_, err = systemClient.Read()
			timeAfter := time.Now()
			if err != nil {
				log.Fatalln(err)
			}

			latencies = append(latencies, timeAfter.Sub(timeBefore).Nanoseconds())
		}
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	gostat.TakeExtremes(latencies)
	latenciesMean := gostat.Mean(latencies)
	latenciesStandardDeviation := gostat.StandardDeviation(latencies)

	latenciesMeanDuration := time.Duration(int64(latenciesMean))
	latenciesStandardDeviationDuration := time.Duration(int64(latenciesStandardDeviation))

	opsPerSecond := float64(ops) / duration.Seconds()

	fmt.Printf("Result: latency %v (%v) - throughput %v [%v in %v]\n", latenciesMeanDuration, latenciesStandardDeviationDuration, int64(opsPerSecond), ops, duration.Seconds())
	saveResults(int64(latenciesMean), int64(latenciesStandardDeviation), int64(opsPerSecond), ops)
}

func saveResults(latenciesMean int64, latenciesStandardDeviation int64, opsPerSecond int64, opsTotal int) {
	var operation string

	if *isWrite {
		operation = "write"
	} else {
		operation = "read"
	}

	file, err := os.OpenFile(*resultFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	defer w.Flush()

	if _, err = w.Write([]byte(fmt.Sprintf("%v %v %v %v %v %v %v\n", latenciesMean, latenciesStandardDeviation, opsPerSecond, opsTotal, operation, *size, time.Now().Format(time.RFC3339)))); err != nil {
		log.Fatalln(err)
	}
}

func latency() {
	if *isWrite {
		log.Printf("Measuring latency of write operations with size %vB\n", *size)
	} else {
		log.Printf("Measuring latency of read operations with size %vB\n", *size)
	}

	data := createFakeData()

	if *isWrite {
		for ops = 0; ops < *numberOfOperations; ops++ {
			timeBefore := time.Now()
			err := systemClient.Write(data)
			timeAfter := time.Now()
			if err != nil {
				log.Fatalln(err)
			}

			latencies = append(latencies, timeAfter.Sub(timeBefore).Nanoseconds())
		}
	} else {
		err := systemClient.Write(data)
		if err != nil {
			log.Fatalln("Initial write:", err)
		}

		for ops = 0; ops < *numberOfOperations; ops++ {
			timeBefore := time.Now()
			_, err = systemClient.Read()
			timeAfter := time.Now()
			if err != nil {
				log.Fatalln(err)
			}

			latencies = append(latencies, timeAfter.Sub(timeBefore).Nanoseconds())
		}
	}

	gostat.TakeExtremes(latencies)
	latenciesMean := gostat.Mean(latencies)
	latenciesStandardDeviation := gostat.StandardDeviation(latencies)

	latenciesMeanDuration := time.Duration(int64(latenciesMean))
	latenciesStandardDeviationDuration := time.Duration(int64(latenciesStandardDeviation))

	fmt.Printf("Result: latency %v (%v) - %v ops\n", latenciesMeanDuration, latenciesStandardDeviationDuration, ops)
	saveResults(int64(latenciesMean), int64(latenciesStandardDeviation), 0, ops)
}

func throughput() {
	if *isWrite {
		log.Printf("Measuring throughput of write operations with size %vB for %v\n", *size, *totalDuration)
	} else {
		log.Printf("Measuring throughput of read operations with size %vB for %v\n", *size, *totalDuration)
	}

	data := createFakeData()

	if *isWrite {
		for ; ; ops++ {
			err := systemClient.Write(data)
			if err != nil {
				log.Fatalln(err)
			}

			select {
			case <-stopChan:
				goto RESULT
			default:
			}
		}
	} else {
		err := systemClient.Write(data)
		if err != nil {
			log.Fatalln("Initial write:", err)
		}

		for ; ; ops++ {
			_, err = systemClient.Read()
			if err != nil {
				log.Fatalln(err)
			}

			select {
			case <-stopChan:
				goto RESULT
			default:
			}
		}
	}

RESULT:
	opsPerSecond := float64(ops) / totalDuration.Seconds()

	fmt.Printf("Result: throughput %v [%v in %v]\n", int64(opsPerSecond), ops, totalDuration.Seconds())
	saveResults(0, 0, int64(opsPerSecond), ops)
}

func createFakeData() []byte {
	data := make([]byte, *size)

	n, err := io.ReadFull(rand.Reader, data)
	if n != len(data) || err != nil {
		log.Fatalln("error to generate data:", err)
	}
	return data
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