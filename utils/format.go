package utils

import (
	"fmt"
	"time"
)

func DurationSince(since time.Time) string {
	ms := time.Since(since).Milliseconds()

	if ms < 1000 {
		return fmt.Sprintf("%.2fms", float64(ms))
	}

	if ms < 60000 {
		s := float64(ms) / 1000
		return fmt.Sprintf("%.2fs", s)
	}

	m := float64(ms) / 60000
	return fmt.Sprintf("%.2fm", m)
}
