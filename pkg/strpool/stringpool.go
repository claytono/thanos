package strpool

import "sync"

// StringPool is used to keep a shared reference to strings that are expected
// to be shared across an instance of the thanos store.  This is mostly used to
// reduce the size of index cache references.
type StringPool struct {
	pool           *sync.Map
	bytesStored    int64
	bytesRetrieved int64
}

// NewStringPool seturns a new StringPool instance
func NewStringPool() *StringPool {
	return &StringPool{pool: &sync.Map{}}
}

// GetCachedString returns a new copy of the string that is shared
func (sc *StringPool) GetCachedString(s string) (cs string, ok bool) {
	iface, ok := sc.pool.LoadOrStore(s, s)
	sc.bytesRetrieved += int64(len(s))
	if !ok {
		sc.bytesStored += int64(len(s))
	}
	return iface.(string), ok
}
