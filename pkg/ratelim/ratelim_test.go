package ratelim_test

import (
	"crypto/rand"
	"fmt"
	"testing"
	"time"

	"github.com/Sengoku11/go-monorepo/pkg/ratelim"
)

// Hash is deterministic, it will consistently produce the same hash value for the same input string.
func ExampleHash() {
	hash, err := ratelim.Hash("Hello 42 world!")

	fmt.Printf("%d, %v\n", hash, err)

	// Output:
	// 8654963907013259302, <nil>
}

func ExampleMessageByTS_Add() {
	mbt := ratelim.NewMap()
	hash, _ := ratelim.Hash("Hello world!")
	now := ratelim.TimeNow()

	mbt.Add(hash, now)

	ts, exist := mbt.Get(hash)
	fmt.Printf("%t, %t\n", ts == now, exist)
	// Output:
	// true, true
}

func ExampleMessageByTS_Get() {
	mbt := ratelim.NewMap()
	hash, _ := ratelim.Hash("Hello world!")

	mbt.Add(hash, 1740925588000)

	ts, exist := mbt.Get(hash)
	fmt.Printf("%d, %t\n", ts, exist)
	// Output:
	// 1740925588000, true
}

func ExampleMessageByTS_AddNow() {
	mbt := ratelim.NewMap()
	hash, _ := ratelim.Hash("Hello world!")

	mbt.AddNow(hash)

	ts, exist := mbt.Get(hash)
	fmt.Printf("%t, %t\n", time.Now().UnixMilli() >= ts, exist)
	// Output:
	// true, true
}

func ExampleMessageByTS_Remove() {
	mbt := ratelim.NewMap()
	hash, _ := ratelim.Hash("Hello world!")

	mbt.AddNow(hash)
	mbt.Remove(hash)

	ts, exist := mbt.Get(hash)
	fmt.Printf("%d, %t\n", ts, exist)
	// Output:
	// 0, false
}

func ExampleIsPastDue() {
	now := ratelim.TimeNow()
	duration := time.Millisecond

	before := ratelim.IsPastDue(now, duration)
	time.Sleep(duration)
	after := ratelim.IsPastDue(now, duration)

	fmt.Printf("%t, %t", before, after)
	// Output:
	// false, true
}

func TestHash(t *testing.T) {
	t.Parallel()

	const iterations = 1_000
	for i := range iterations {
		txt := rand.Text()

		hash1, err := ratelim.Hash(txt)
		if err != nil {
			t.Fatalf("Iteration %d: unexpected error hashing text %q: %v", i, txt, err)
		}

		hash2, err := ratelim.Hash(txt)
		if err != nil {
			t.Fatalf("Iteration %d: unexpected error on second hash call for text %q: %v", i, txt, err)
		}

		if hash1 != hash2 {
			t.Fatalf("Iteration %d: non-deterministic hash for text %q: got %v and %v", i, txt, hash1, hash2)
		}
	}
}

func TestNewMap(t *testing.T) {
	t.Parallel()

	if m := ratelim.NewMap(); m == nil {
		t.Errorf("return nil map")
	}
}

func TestMessageByTS(t *testing.T) {
	t.Parallel()

	size := 1_000
	expectedMap := make(map[uint64]int64, size)

	for i := range size {
		timestamp := ratelim.TimeNow() + int64(i)
		txt := rand.Text()

		hash, err := ratelim.Hash(txt)
		if err != nil {
			t.Fatalf("Iteration %d: unexpected error hashing text %q: %v", i, txt, err)
		}

		expectedMap[hash] = timestamp
	}

	testMap := ratelim.NewMap()

	for hash, timestamp := range expectedMap {
		// Run concurrently to test mutual exclusivity.
		t.Run(fmt.Sprintf("Add-%d", hash), func(t *testing.T) {
			t.Parallel()

			testMap.Add(hash, timestamp)

			// Must use read lock and return stored values.
			if receivedTS, exist := testMap.Get(hash); !exist {
				t.Fatalf("cannot fetch stored hash %d", hash)
			} else if receivedTS != timestamp {
				t.Fatalf("stored timestamp is different from the original: %d and %d", receivedTS, timestamp)
			}

			// Remove the key and verify removal.
			testMap.Remove(hash)

			if _, exist := testMap.Get(hash); exist {
				t.Fatalf("key %d still exists after Remove", hash)
			}
		})
	}
}
