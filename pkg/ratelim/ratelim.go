// Package ratelim provides utility functions to prevent alert spam.
package ratelim

import (
	"fmt"
	"hash"
	"hash/fnv"
	"sync"
	"time"
)

// Hash64 for mock testing.
type Hash64 interface {
	Write(p []byte) (n int, err error)
	Sum64() uint64
}

// New64a returns Hash64 function.
func New64a() hash.Hash64 {
	return fnv.New64a()
}

// Hash the given message into a deterministic 64-bit FNV-1a number
// if optional hashFunc not provided.
func Hash(msg string, hashFunc Hash64) (uint64, error) {
	if hashFunc == nil {
		hashFunc = New64a()
	}

	_, err := hashFunc.Write([]byte(msg))
	if err != nil {
		return 0, fmt.Errorf("error hashing message: %w", err)
	}

	return hashFunc.Sum64(), nil
}

// MessageByTS provides a thread-safe map optimized to reduce garbage collection
// and memory overhead by avoiding string allocations.
//
// Use the Hash function to generate keys and TimeNow for compatible timestamp values.
type MessageByTS struct {
	mu sync.RWMutex
	m  map[uint64]int64
}

// NewMap returns an optimized map for storing messages' timestamps.
func NewMap() *MessageByTS {
	m := make(map[uint64]int64)

	return &MessageByTS{
		m:  m,
		mu: sync.RWMutex{},
	}
}

// Add method uses mutex to update the map safely.
//
// Use only TimeNow to get a compatible timestamp.
func (m *MessageByTS) Add(hash uint64, ts int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.m[hash] = ts
}

// AddNow calls TimeNow method to store a compatible timestamp for the given hash.
func (m *MessageByTS) AddNow(hash uint64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.m[hash] = TimeNow()
}

// Get method uses read lock for concurrent safe reads.
func (m *MessageByTS) Get(hash uint64) (int64, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ts, ok := m.m[hash]

	return ts, ok
}

// Remove method safely deletes hash from the map.
func (m *MessageByTS) Remove(hash uint64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.m, hash)
}

// TimeNow returns the current time in UnixMilli format,
// ensuring consistent interpretation of the time value.
func TimeNow() int64 {
	return time.Now().UnixMilli()
}

// IsPastDue determines whether the provided UnixMilli timestamp is older than the specified duration.
func IsPastDue(timestamp int64, duration time.Duration) bool {
	return time.Since(time.UnixMilli(timestamp)) > duration
}
