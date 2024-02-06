package block_stm

import "sync/atomic"

// StoreMin implements a compare-and-swap operation that stores the minimum of the current value and the given value.
func StoreMin(a *atomic.Uint64, b uint64) {
	for {
		old := a.Load()
		if old <= b {
			return
		}
		if a.CompareAndSwap(old, b) {
			return
		}
	}
}

func DecreaseAtomic(a *atomic.Uint64) {
	a.Add(^uint64(0))
}

func IncreaseAtomic(a *atomic.Uint64) {
	a.Add(1)
}
