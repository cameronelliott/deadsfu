package disrupt

import (
	"log"
	"sync"
	"sync/atomic"
	"unsafe"
)

type nolock struct{}

func (*nolock) Lock()   {}
func (*nolock) Unlock() {}

//don't use const for CacheLine size, it hits a linter bug
// SPMC
type Disrupt[T any] struct {
	next   int64
	_2     [64]byte
	len64  int64
	_3     [64]byte
	mask64 int64
	_4     [64]byte
	cond   sync.Cond
	_5     [64]byte
	buf    []T
	_1     [64]byte
}

func NewDisrupt[T any](n int64) *Disrupt[T] {

	if n <= 0 || (n&(n-1)) != 0 {
		log.Fatal("require  power of two")
	}

	buf := make([]T, n)

	return &Disrupt[T]{
		cond:   sync.Cond{L: &nolock{}},
		buf:    buf,
		len64:  int64(n),
		mask64: int64(n - 1),
	}
}

func (d *Disrupt[T]) NextIx() int64 {
	i := atomic.LoadInt64(&d.next)
	if i < 0 {
		i = -i
	}
	return i
}

func (d *Disrupt[T]) Close() {

	i := atomic.LoadInt64(&d.next)
	if i < 0 {
		log.Fatal("cannot close twice")
	}
	atomic.StoreInt64(&d.next, -i)

	d.cond.Broadcast()
	//d.cond.Signal()   //must uncomment signal below when using this
}

// not yet: for non-pointer, raw struct usage
func (d *Disrupt[T]) GetEmptyPtr(k int64) *T {
	log.Fatal("not yet")
	return nil
}

// not yet: for non-pointer, raw struct usage
func (d *Disrupt[T]) PutFilledPtr(k int64) {
	log.Fatal("not yet")
	return
}

func (d *Disrupt[T]) Put(v T) (index int64) {

	index = atomic.LoadInt64(&d.next)
	//	ix := i % d.len64
	ix := index & d.mask64

	if index < 0 {
		log.Fatal("closed")
	}

	if RaceEnabled {
		RaceAcquire(unsafe.Pointer(d))
		//RaceDisable()
	}

	d.buf[ix] = v

	if RaceEnabled {
		//RaceEnable()
		RaceReleaseMerge(unsafe.Pointer(d))
	}

	index++
	atomic.StoreInt64(&d.next, index)

	d.cond.Broadcast()
	//d.cond.Signal()   //must uncomment signal below when using this

	return index - 1 // return index of last stored element
}

func (d *Disrupt[T]) GetLast() (value T, next int64, more bool) {

	i := atomic.LoadInt64(&d.next)
	if i < 0 {
		i = -i
	}
	i--
	i = i & d.mask64
	return d.Get(i)

}

func (d *Disrupt[T]) Get(k int64) (value T, next int64, more bool) {

	//ix := k % d.len64

again:

	ix := k & d.mask64

	// if k >= i {
	// 	log.Fatal("invalid index too high")
	// }

	for j := atomic.LoadInt64(&d.next); k >= j; j = atomic.LoadInt64(&d.next) {
		if j < 0 { //closed?
			if k >= -j { // no values left
				var zeroval T
				return zeroval, k, false
			}
			break
		} else {
			d.cond.Wait()
			//d.cond.Signal() // wake any other waiters, when not using broadcast in Put
		}
	}

	if RaceEnabled {
		RaceAcquire(unsafe.Pointer(d))
		//RaceDisable()
	}

	val := d.buf[ix]

	if RaceEnabled {
		//RaceEnable()
		RaceReleaseMerge(unsafe.Pointer(d))
	}

	// did we grab stale data?

	j := atomic.LoadInt64(&d.next)

	if j < 0 { // if closed, fix sign
		j = -j
	}
	if k <= (j - d.len64) { // we read possibly overwritten data

		k++                              //discard bad data
		log.Println("disrupt data loss") // XXX better way to track someday
		goto again

		// val = zeroval
		// k = j - 1
	}

	k++

	return val, k, true

}
