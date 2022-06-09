package mutexmap

import (
	"math/rand"
	"runtime"
	"sync"
	"testing"
)

// Benchmarks taken from: https://medium.com/@deckarep/the-new-kid-in-town-gos-sync-map-de24a6bf7c2c

// Store result outside the benchmark to not let the compiler optimize the benchmark away
var globalResult int
var globalResultChan = make(chan int, 100)

func nrand(n int) []int {
	ints := make([]int, n)
	for i := range ints {
		ints[i] = rand.Int()
	}
	return ints
}

func populateMap(n int, m *MutexMap[int, int]) []int {
	ints := nrand(n)
	for _, v := range ints {
		m.Store(v, v)
	}
	return ints
}

func populateSyncMap(n int, s *sync.Map) []int {
	ints := nrand(n)
	for _, v := range ints {
		s.Store(v, v)
	}
	return ints
}

func BenchmarkStoreMutex(b *testing.B) {
	values := nrand(b.N)
	m := NewMutexMap[int, int]()

	b.ResetTimer()
	for _, v := range values {
		m.Store(v, v)
	}
}

func BenchmarkStoreSync(b *testing.B) {
	values := nrand(b.N)
	var s sync.Map

	b.ResetTimer()
	for _, v := range values {
		s.Store(v, v)
	}
}

func BenchmarkDeleteMutex(b *testing.B) {
	values := nrand(b.N)
	m := NewMutexMap[int, int]()

	for _, v := range values {
		m.Store(v, v)
	}

	b.ResetTimer()
	for _, v := range values {
		m.Delete(v)
	}
}

func BenchmarkDeleteSync(b *testing.B) {
	values := nrand(b.N)
	var s sync.Map

	for _, v := range values {
		s.Store(v, v)
	}

	b.ResetTimer()
	for _, v := range values {
		s.Delete(v)
	}
}

func BenchmarkLoadMutexFound(b *testing.B) {
	values := nrand(b.N)
	m := NewMutexMap[int, int]()

	for _, v := range values {
		m.Store(v, v)
	}
	localResult := 0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r, ok := m.Load(values[i])
		if ok {
			localResult = r // Use output to prevent optimization
		}
	}
	globalResult = localResult
}

func BenchmarkLoadSyncFound(b *testing.B) {
	values := nrand(b.N)
	var s sync.Map

	for _, v := range values {
		s.Store(v, v)
	}
	localResult := 0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r, ok := s.Load(values[i])
		if ok {
			localResult = r.(int)
		}
	}
	globalResult = localResult
}

func BenchmarkLoadMutexNotFound(b *testing.B) {
	values := nrand(b.N)
	m := NewMutexMap[int, int]()

	for _, v := range values {
		m.Store(v, v)
	}
	localResult := 0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r, ok := m.Load(i)
		if ok {
			localResult = r
		}
	}
	globalResult = localResult
}

func BenchmarkLoadSyncNotFound(b *testing.B) {
	values := nrand(b.N)
	var s sync.Map

	for _, v := range values {
		s.Store(v, v)
	}
	localResult := 0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r, ok := s.Load(i)
		if ok {
			localResult = r.(int)
		}
	}
	globalResult = localResult
}

func BenchmarkMutexStableKeys1(b *testing.B) {
	benchmarkMutexStableKeys(b, 1)
}
func BenchmarkMutexStableKeys2(b *testing.B) {
	benchmarkMutexStableKeys(b, 2)
}
func BenchmarkMutexStableKeys3(b *testing.B) {
	benchmarkMutexStableKeys(b, 3)
}
func BenchmarkMutexStableKeys4(b *testing.B) {
	benchmarkMutexStableKeys(b, 4)
}

func benchmarkMutexStableKeys(b *testing.B, numCPUs int) {
	runtime.GOMAXPROCS(numCPUs)

	m := NewMutexMap[int, int]()
	populateMap(b.N, m)

	var wg sync.WaitGroup
	wg.Add(numCPUs)

	globalResultChan = make(chan int, numCPUs)

	b.ResetTimer()
	for cpu := 0; cpu < numCPUs; cpu++ {
		go func(n int) {
			localResult := 0
			for i := 0; i < n; i++ {
				r, ok := m.Load(5)
				if ok {
					localResult = r
				}
			}
			globalResultChan <- localResult
			wg.Done()
		}(b.N)
	}

	wg.Wait()
}

func BenchmarkSyncStableKeys1(b *testing.B) {
	benchmarkSyncStableKeys(b, 1)
}
func BenchmarkSyncStableKeys2(b *testing.B) {
	benchmarkSyncStableKeys(b, 2)
}
func BenchmarkSyncStableKeys3(b *testing.B) {
	benchmarkSyncStableKeys(b, 3)
}
func BenchmarkSyncStableKeys4(b *testing.B) {
	benchmarkSyncStableKeys(b, 4)
}

func benchmarkSyncStableKeys(b *testing.B, numCPUs int) {
	runtime.GOMAXPROCS(numCPUs)

	var s sync.Map
	populateSyncMap(b.N, &s)

	var wg sync.WaitGroup
	wg.Add(numCPUs)

	globalResultChan = make(chan int, numCPUs)

	b.ResetTimer()
	for cpu := 0; cpu < numCPUs; cpu++ {
		go func(n int) {
			localResult := 0
			for i := 0; i < n; i++ {
				r, ok := s.Load(5)
				if ok {
					localResult = r.(int)
				}
			}
			globalResultChan <- localResult
			wg.Done()
		}(b.N)
	}

	wg.Wait()
}

func BenchmarkMutexStableKeysFound1(b *testing.B) {
	benchmarkMutexStableKeysFound(b, 1)
}
func BenchmarkMutexStableKeysFound2(b *testing.B) {
	benchmarkMutexStableKeysFound(b, 2)
}
func BenchmarkMutexStableKeysFound3(b *testing.B) {
	benchmarkMutexStableKeysFound(b, 3)
}
func BenchmarkMutexStableKeysFound4(b *testing.B) {
	benchmarkMutexStableKeysFound(b, 4)
}

func benchmarkMutexStableKeysFound(b *testing.B, numCPUs int) {
	runtime.GOMAXPROCS(numCPUs)

	m := NewMutexMap[int, int]()
	values := populateMap(b.N, m)

	var wg sync.WaitGroup
	wg.Add(numCPUs)

	globalResultChan = make(chan int, numCPUs)

	b.ResetTimer()
	for cpu := 0; cpu < numCPUs; cpu++ {
		go func(n int) {
			localResult := 0
			for i := 0; i < n; i++ {
				r, ok := m.Load(values[i])
				if ok {
					localResult = r
				}
			}
			globalResultChan <- localResult
			wg.Done()
		}(b.N)
	}

	wg.Wait()
}

func BenchmarkSyncStableKeysFound1(b *testing.B) {
	benchmarkSyncStableKeysFound(b, 1)
}
func BenchmarkSyncStableKeysFound2(b *testing.B) {
	benchmarkSyncStableKeysFound(b, 2)
}
func BenchmarkSyncStableKeysFound3(b *testing.B) {
	benchmarkSyncStableKeysFound(b, 3)
}
func BenchmarkSyncStableKeysFound4(b *testing.B) {
	benchmarkSyncStableKeysFound(b, 4)
}

func benchmarkSyncStableKeysFound(b *testing.B, numCPUs int) {
	runtime.GOMAXPROCS(numCPUs)

	var s sync.Map
	values := populateSyncMap(b.N, &s)

	var wg sync.WaitGroup
	wg.Add(numCPUs)

	globalResultChan = make(chan int, numCPUs)

	b.ResetTimer()
	for cpu := 0; cpu < numCPUs; cpu++ {
		go func(n int) {
			localResult := 0
			for i := 0; i < n; i++ {
				r, ok := s.Load(values[i])
				if ok {
					localResult = r.(int)
				}
			}
			globalResultChan <- localResult
			wg.Done()
		}(b.N)
	}

	wg.Wait()
}
