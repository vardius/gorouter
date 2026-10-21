//go:build race

// Package race reports whether the race detector is on. Under it, sync.Pool drops items at random,
// so allocation counts mean nothing and the allocation tests skip.
package race

// Enabled is true when built with -race.
const Enabled = true
