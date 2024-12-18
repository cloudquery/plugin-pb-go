package set

import (
	"sort"
	"strings"
	"sync"
)

// SyncSortedStringSet is a concurrency-safe sorted set of strings.
//
// - Use `Add` to add a string to the set.
// - Use `Get` to get a comma-separated, sorted string of all the strings in the set.
type SyncSortedStringSet struct {
	mu      sync.Mutex
	strings map[string]struct{}
}

func NewSyncSortedStringSet() *SyncSortedStringSet {
	return &SyncSortedStringSet{
		strings: make(map[string]struct{}),
	}
}

// Add adds a string to the set.
func (e *SyncSortedStringSet) Add(str string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.strings[str] = struct{}{}
}

// Get returns a comma-separated, sorted string of all the strings in the set.
func (e *SyncSortedStringSet) Get() string {
	if e == nil {
		return ""
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	keys := make([]string, 0, len(e.strings))
	for k := range e.strings {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return strings.Join(keys, ",")
}
