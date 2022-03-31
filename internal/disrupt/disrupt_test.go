package disrupt

import (
	"math"

	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var A, B int64

var X, Y uint64

func BenchmarkLoadInt64(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		k := int64(0)
		for p.Next() {
			k += atomic.LoadInt64(&A)
		}
		atomic.StoreInt64(&B, k)
	})
}

func BenchmarkLoadUInt64(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		k := uint64(0)
		for p.Next() {
			k += atomic.LoadUint64(&X)
		}
		atomic.StoreUint64(&Y, k)
	})
}

func RoundUpPowerOf2(n int64) (power int64) {
	for power = int64(1); power < n; power *= 2 {
	}
	return
}
func BenchmarkDisrupt2pow30(b *testing.B) {
	benchmarkDisrupt(b, int64(math.Pow(2, 30)))
}
func BenchmarkDisrupt2pow18(b *testing.B) {
	benchmarkDisrupt(b, int64(math.Pow(2, 18)))
}

func BenchmarkDisruptNotConcurrent(b *testing.B) {
	benchmarkDisrupt(b, 0)
}

func benchmarkDisrupt(b *testing.B, bufsize int64) {

	notConcurrent := false // run get/put concurrent or serially

	if bufsize == 0 {
		notConcurrent = true
		bufsize = RoundUpPowerOf2(int64(b.N))

		// if bufsize<4096{
		// 	bufsize=4096
		// }
	}
	b.Log("benchmarkDisrupt", "b.n", b.N, "bufsize", bufsize, "notConcurrent", notConcurrent)

	w := sync.WaitGroup{}

	w.Add(1)

	a := NewDisrupt[int64](bufsize, -5)

	go func() {

		//runtime.LockOSThread()
		time.Sleep(time.Millisecond)

		for i := int64(0); i < int64(b.N); i++ {
			a.Put(i)
		}
		a.Close()

		w.Done()

	}()

	if notConcurrent {
		w.Wait()
	}

	b.ResetTimer()

	//runtime.LockOSThread()
	for ix := int64(0); ix < int64(b.N); ix++ {
		val, nextix, isopen := a.Get(ix, false)
		if ix != val || ix+1 != nextix || isopen != true {
			require.Equal(b, ix, val)
			require.Equal(b, ix+1, nextix)
			if b.N == 1 {
				require.Equal(b, false, isopen) // special case, buffer overrun
			} else {
				require.Equal(b, true, isopen) //normal case
			}

		}
	}

	if !notConcurrent {
		w.Wait()
	}

	b.StopTimer()

	var zero int64
	v, ix, isopen := a.Get(int64(b.N), false)
	require.Equal(b, zero, v)
	require.Equal(b, int64(b.N), ix)
	require.Equal(b, false, isopen)

}

func TestDisrupt(t *testing.T) {
	testDisrupt(t, int64(math.Pow(2, 20)), 1e6) // if bufsize too small, overrun and test will fail
}

func testDisrupt(t *testing.T, bufsize int64, NN int64) {

	w := sync.WaitGroup{}

	w.Add(1)

	a := NewDisrupt[int64](bufsize, -5)

	go func() {

		//runtime.LockOSThread()
		time.Sleep(time.Millisecond)

		for i := int64(0); i < int64(NN); i++ {
			a.Put(i)
		}
		a.Close()

		w.Done()

	}()

	//runtime.LockOSThread()
	for ix := int64(0); ix < int64(NN); ix++ {
		val, nextix, isopen := a.Get(ix, false)
		if ix != val || ix+1 != nextix || isopen != true {
			require.Equal(t, ix, val)
			require.Equal(t, ix+1, nextix)
			require.Equal(t, true, isopen)
		}
	}

	w.Wait()

	var zero int64
	v, ix, isopen := a.Get(int64(NN), false)
	require.Equal(t, zero, v)
	require.Equal(t, int64(NN), ix)
	require.Equal(t, false, isopen)

}

var C, D [1500]byte

func Benchmark1500Copy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		C = D

	}
}

func TestPutGet1(t *testing.T) {

	w := sync.WaitGroup{}

	w.Add(1)
	a := NewDisrupt[int64](256, -5)

	go func() {

		t0 := time.Now()

		v, i, ok := a.Get(0, false)

		assert.Equal(t, int64(99), v)
		assert.Equal(t, int64(1), i)
		assert.Equal(t, true, ok)
		dur := time.Since(t0)
		assert.GreaterOrEqual(t, dur, time.Millisecond, ok)

		w.Done()

	}()

	time.Sleep(time.Millisecond)
	a.Put(99)

	w.Wait()

}
