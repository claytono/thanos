package strpool

import "sync"

// StringPool is used to keep a shared reference to strings that are expected
// to be shared across an instance of the thanos store.  This is mostly used to
// reduce the size of index cache references.
type StringPool struct {
	m              *sync.Mutex
	pool           map[string]string
	bytesStored    int64
	bytesRetrieved int64
}

// NewStringPool seturns a new StringPool instance
func NewStringPool() *StringPool {
	return &StringPool{
		m:    &sync.Mutex{},
		pool: map[string]string{},
	}
}

// Lock locks the pool for read or write
func (sc *StringPool) Lock() {
	sc.m.Lock()
}

// Unlock unlocks the pool for read or write
func (sc *StringPool) Unlock() {
	sc.m.Unlock()
}

// GetCachedStringWithLock returns a copy of the string that is shared if it is
// already in the pool, otherwise it returns the original string and stores it
// in the pool.  It locks the datastructure for read/write so it can be called
// across threads with no special measures.
func (sc *StringPool) GetCachedStringWithLock(s string) (cs string, ok bool) {
	sc.Lock()
	defer sc.Unlock()
	sc.bytesRetrieved += int64(len(s))

	cs, ok = sc.pool[s]
	if ok {
		return cs, ok
	}

	sc.bytesRetrieved += int64(len(s))
	sc.pool[s] = s
	return s, ok
}

// GetCachedString returns a copy of the string that is shared if it is already
// in the pool, otherwise it returns the original string and stores it in the
// pool.  It doesn't lock the datastructure and you should call Lock() and
// Unlock() yourself to ensure thread safety.
func (sc *StringPool) GetCachedString(s string) (cs string, ok bool) {
	sc.bytesRetrieved += int64(len(s))

	cs, ok = sc.pool[s]
	if ok {
		return cs, ok
	}

	sc.bytesRetrieved += int64(len(s))
	sc.pool[s] = s
	return s, ok
}
