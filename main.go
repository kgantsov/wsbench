package main

import (
	"flag"
)

var (
	logLevel    string // Log level
	name        string // Name of the benchmark
	serverURL   string // URL of the WebSocket server
	numClients  int    // Number of concurrent WebSocket clients
	numMessages int    // Number of messages each client will send
)

func main() {
	flag.StringVar(&name, "name", "WebSocket Benchmark", "Name of the benchmark")
	flag.StringVar(&serverURL, "url", "ws://localhost:8080/ws", "WebSocket server URL")
	flag.IntVar(&numClients, "clients", 100, "Number of concurrent WebSocket clients")
	flag.IntVar(&numMessages, "messages", 1000, "Number of messages each client will send")

	flag.Parse()

	bench := NewWSBench(name, numClients, numMessages)
	bench.run()
}
