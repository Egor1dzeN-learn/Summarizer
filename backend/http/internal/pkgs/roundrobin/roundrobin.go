package roundrobin

import (
	"sync/atomic"
)

type RoundRobin[T any] struct {
	elements []*T
	next     uint32
}

func New[T any](elements ...*T) RoundRobin[T] {
	return RoundRobin[T]{
		elements: elements,
		next:     0,
	}
}

func (r *RoundRobin[T]) Next() *T {
	n := atomic.AddUint32(&r.next, 1)
	return r.elements[(int(n)-1)%len(r.elements)]
}
