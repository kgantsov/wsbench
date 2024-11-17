package main

import (
	"sort"
	"time"
)

type stats struct {
	avg   time.Duration
	total time.Duration
	min   time.Duration
	p50   time.Duration
	p90   time.Duration
	p95   time.Duration
	p99   time.Duration
	max   time.Duration
}

func calcStats(latencies []time.Duration) stats {
	if len(latencies) == 0 {
		return stats{}
	}

	avgLatency := time.Duration(0)
	totalLatency := time.Duration(0)
	for _, latency := range latencies {
		totalLatency += latency
	}
	avgLatency = totalLatency / time.Duration(len(latencies))

	sort.Slice(latencies, func(i, j int) bool {
		return latencies[i] < latencies[j]
	})

	p50 := latencies[len(latencies)/2]
	p90 := latencies[len(latencies)*90/100]
	p95 := latencies[len(latencies)*95/100]
	p99 := latencies[len(latencies)*99/100]

	return stats{
		avg:   avgLatency,
		total: totalLatency,
		min:   latencies[0],
		p50:   p50,
		p90:   p90,
		p95:   p95,
		p99:   p99,
		max:   latencies[len(latencies)-1],
	}
}
