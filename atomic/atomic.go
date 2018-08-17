package atomic

import (
	"sync/atomic"
)

var (
	default_atomic_uint64 uint64 = 0
)

func GetDefUInt64() uint64 {
	return atomic.LoadUint64(&default_atomic_uint64)
}

func SetDefUInt64(val uint64) {
	atomic.StoreUint64(&default_atomic_uint64, val)
}

func ResetDefUInt64() {
	atomic.StoreUint64(&default_atomic_uint64, 0)
}

func IncrDefUInt64() {
	atomic.AddUint64(&default_atomic_uint64, 1)
}

func DecrDefUInt64() {
	atomic.AddUint64(&default_atomic_uint64, ^uint64(0))
}

func AddDefUInt64(delta uint64) {
	atomic.AddUint64(&default_atomic_uint64, delta)
}

func SubtractDefUInt64(delta uint64) {
	atomic.AddUint64(&default_atomic_uint64, ^uint64(delta-1))
}

func GetUInt64(addr *uint64) uint64 {
	return atomic.LoadUint64(addr)
}

func SetUInt64(addr *uint64, val uint64) {
	atomic.StoreUint64(addr, val)
}

func ResetUInt64(addr *uint64) {
	atomic.StoreUint64(addr, 0)
}

func IncrUInt64(addr *uint64) {
	atomic.AddUint64(addr, 1)
}

func DecrUInt64(addr *uint64) {
	atomic.AddUint64(addr, ^uint64(0))
}

func AddUInt64(addr *uint64, delta uint64) {
	atomic.AddUint64(addr, delta)
}

func SubtractUInt64(addr *uint64, delta uint64) {
	atomic.AddUint64(addr, ^uint64(delta-1))
}
