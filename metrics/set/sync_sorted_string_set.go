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
	mu     sync.Mutex
	tables map[string]struct{}
}

func NewSyncSortedStringSet() *SyncSortedStringSet {
	return &SyncSortedStringSet{
		tables: make(map[string]struct{}),
	}
}

func (e *SyncSortedStringSet) Add(table string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.tables[table] = struct{}{}
}

func (e *SyncSortedStringSet) Get() string {
	e.mu.Lock()
	defer e.mu.Unlock()

	keys := make([]string, 0, len(e.tables))
	for k := range e.tables {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return strings.Join(keys, ",")
}
