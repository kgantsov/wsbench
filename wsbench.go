package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type wsbench struct {
	name        string
	numClients  int
	numMessages int

	latencies      []time.Duration
	latenciesMutex sync.Mutex

	connLatencies      []time.Duration
	connLatenciesMutex sync.Mutex

	bytesSent     int
	bytesReceived int

	benchmarkDuration time.Duration
	stats             stats
	connStats         stats
}

func NewWSBench(name string, numClients, numMessages int) *wsbench {
	return &wsbench{
		name:               name,
		numClients:         numClients,
		numMessages:        numMessages,
		latencies:          make([]time.Duration, 0, numClients*numMessages),
		latenciesMutex:     sync.Mutex{},
		connLatencies:      make([]time.Duration, 0, numClients*numMessages),
		connLatenciesMutex: sync.Mutex{},
	}
}

func (w *wsbench) run() {
	var wg sync.WaitGroup
	wg.Add(numClients)

	messagePayload := strings.Repeat("a", 1024)

	started := time.Now()

	for i := 0; i < numClients; i++ {
		go func(clientID int) {
			defer wg.Done()

			// Add a random jitter to the client to avoid thundering herd problem
			delay := random(5, 20)
			time.Sleep(time.Duration(delay) * time.Millisecond)

			connStart := time.Now()
			conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
			if err != nil {
				fmt.Printf("Client %d: could not connect to server: %v\n", clientID, err)
				return
			}
			defer conn.Close()

			w.connLatenciesMutex.Lock()
			w.connLatencies = append(w.connLatencies, time.Since(connStart))
			w.connLatenciesMutex.Unlock()

			for j := 0; j < numMessages; j++ {
				start := time.Now()

				messageContent := fmt.Sprintf("Client: %d message_id: %d payload: %s", clientID, j, messagePayload)

				if err := conn.WriteMessage(websocket.TextMessage, []byte(messageContent)); err != nil {
					fmt.Printf("Client %d: error writing message: %v\n", clientID, err)
					return
				}

				_, msg, err := conn.ReadMessage()
				if err != nil {
					fmt.Printf("Client %d: error reading message: %v\n", clientID, err)
					return
				}

				elapsed := time.Since(start)

				w.latenciesMutex.Lock()
				w.latencies = append(w.latencies, elapsed)
				w.bytesSent += len(messageContent)
				w.bytesReceived += len(msg)
				w.latenciesMutex.Unlock()
			}
		}(i)
	}

	wg.Wait()

	w.benchmarkDuration = time.Since(started)

	w.stats = calcStats(w.latencies)
	w.connStats = calcStats(w.connLatencies)

	w.printStats()
}

func (w *wsbench) printStats() {
	expectedMessagers := w.numClients * w.numMessages
	fmt.Printf("Benchmark: %s\n", name)
	fmt.Printf(
		"clients..............................: %d\n",
		numClients,
	)
	fmt.Printf(
		"messages_per_client..................: %d\n",
		numMessages,
	)
	fmt.Printf(
		"data_sent............................: %s %s/s\n",
		formatSize(float64(w.bytesSent), 1024.0),
		formatSize(float64(w.bytesSent)/w.benchmarkDuration.Seconds(), 1024.0),
	)
	fmt.Printf(
		"data_received........................: %s %s/s\n",
		formatSize(float64(w.bytesReceived), 1024.0),
		formatSize(float64(w.bytesReceived)/w.benchmarkDuration.Seconds(), 1024.0),
	)
	fmt.Printf(
		"successes............................: %.2f%% ✓%d ✗%d\n",
		(float64(len(w.latencies))/float64(expectedMessagers))*100,
		len(w.latencies),
		(expectedMessagers)-len(w.latencies),
	)
	fmt.Printf(
		"fails................................: %.2f%% ✓%d ✗%d\n",
		(float64(expectedMessagers-len(w.latencies))/float64(expectedMessagers))*100,
		(expectedMessagers)-len(w.latencies),
		len(w.latencies),
	)
	fmt.Printf(
		"latency_connecting...................: sum=%s avg=%s min=%s med=%s max=%s p(90)=%s p(95)=%s p(99)=%s\n",
		w.connStats.total,
		w.connStats.avg,
		w.connStats.min,
		w.connStats.p50,
		w.connStats.max,
		w.connStats.p90,
		w.connStats.p95,
		w.connStats.p99,
	)
	fmt.Printf(
		"latency..............................: sum=%s avg=%s min=%s med=%s max=%s p(90)=%s p(95)=%s p(99)=%s\n",
		w.stats.total,
		w.stats.avg,
		w.stats.min,
		w.stats.p50,
		w.stats.max,
		w.stats.p90,
		w.stats.p95,
		w.stats.p99,
	)
	fmt.Printf(
		"throughput...........................: %.2f/s\n",
		float64(len(w.latencies))/(w.stats.total.Seconds()/float64(numClients)),
	)
	fmt.Printf(
		"throughput_with_connecting...........: %.2f/s\n",
		float64(len(w.latencies))/float64(w.benchmarkDuration.Seconds()),
	)
	fmt.Printf(
		"duration_benchmark...................: %s\n",
		w.benchmarkDuration,
	)
}
