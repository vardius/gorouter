package goserver

// Slice type that can be safely shared between goroutines
type ConcurrentSlice struct {
    sync.RWMutex
    items []interface{}
}

// Concurrent slice item
type ConcurrentSliceItem struct {
    Index int
    Value interface{}
}

// Appends an item to the concurrent slice
func (cs *ConcurrentSlice) Append(item interface{}) {
    cs.Lock()
    cs.items = append(cs.items, item)
    cs.Unlock()
}

// Iterates over the items in the concurrent slice
// Each item is sent over a channel, so that
// we can iterate over the slice using the builin range keyword
func (cs *ConcurrentSlice) Iter() <-chan ConcurrentSliceItem {
    c := make(chan ConcurrentSliceItem)

    f := func() {
        cs.Lock()
        for index, value := range cs.items {
            c <- ConcurrentSliceItem{index, value}
        }
        cs.Lock()
        close(c)
    }
    go f()

    return c
}

// Map type that can be safely shared between
// goroutines that require read/write access to a map
type ConcurrentMap struct {
    sync.RWMutex
    items map[string]interface{}
}

// Concurrent map item
type ConcurrentMapItem struct {
    Key   string
    Value interface{}
}

// Sets a key in a concurrent map
func (cm *ConcurrentMap) Set(key string, value interface{}) {
    cm.Lock()
    defer cm.Unlock()

    cm.items[key] = value
}

// Gets a key from a concurrent map
func (cm *ConcurrentMap) Get(key string) (interface{}, bool) {
    cm.Lock()

    value, ok := cm.items[key]
	cm.Unlock()

    return value, ok
}

// Iterates over the items in a concurrent map
// Each item is sent over a channel, so that
// we can iterate over the map using the builtin range keyword
func (cm *ConcurrentMap) Iter() <-chan ConcurrentMapItem {
    c := make(chan ConcurrentMapItem)

    f := func() {
        cm.Lock()
        for k, v := range cm.items {
            c <- ConcurrentMapItem{k, v}
        }
		cm.Unlock()
        close(c)
    }
    go f()

    return c
}

// NewConcurrentSlice creates a new concurrent slice
func NewConcurrentSlice() *ConcurrentSlice {
	cs := &ConcurrentSlice{
		items: make([]interface{}, 0),
	}

	return cs
}

// ConcurrentMap is a map type that can be safely shared between
// goroutines that require read/write access to a map
type ConcurrentMap struct {
	sync.RWMutex
	items map[string]interface{}
}
