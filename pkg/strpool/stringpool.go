package strpool

import "sync"

// StringPool is used to keep a shared reference to strings that are expected
// to be shared across an instance of the thanos store.  This is mostly used to
// reduce the size of index cache references.
type StringPool struct {
	rwm            *sync.RWMutex
	pool           map[string]string
	bytesStored    int64
	bytesRetrieved int64
}

// NewStringPool seturns a new StringPool instance
func NewStringPool() *StringPool {
	return &StringPool{
		rwm:  &sync.RWMutex{},
		pool: map[string]string{},
	}
}

// GetCachedString returns a copy of the string that is shared if it is already
// in the pool, otherwise it returns the original string and stores it in the
// pool.
func (sc *StringPool) GetCachedString(s string) (cs string, ok bool) {
	sc.bytesRetrieved += int64(len(s))

	sc.rwm.RLock()
	cs, ok = sc.pool[s]
	sc.rwm.RUnlock()
	if ok {
		return cs, ok
	}

	sc.bytesRetrieved += int64(len(s))
	sc.rwm.Lock()
	sc.pool[s] = s
	sc.rwm.Unlock()
	return s, ok
}
