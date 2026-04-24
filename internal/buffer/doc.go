// Package buffer implements a fixed-capacity ring buffer used to retain
// recent port-scan result snapshots in memory.
//
// When the buffer reaches its configured capacity, the oldest entry is
// automatically evicted to make room for the newest one. All operations
// are safe for concurrent use.
package buffer
