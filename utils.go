package main

import (
	"fmt"
	"math/rand"
)

func random(min int, max int) int {
	return rand.Intn(max-min) + min
}

var sizes = []string{"B", "kB", "MB", "GB", "TB", "PB", "EB"}

func formatSize(s float64, base float64) string {
	unitsLimit := len(sizes)
	i := 0
	for s >= base && i < unitsLimit {
		s = s / base
		i++
	}

	f := "%.0f %s"
	if i > 1 {
		f = "%.2f %s"
	}

	return fmt.Sprintf(f, s, sizes[i])
}
