package rolling

import (
	"time"

	"go.uber.org/atomic"
)

type numberBucket struct {
	atomic.Float64
}

func (b *numberBucket) Reset() {
	b.Store(0)
}

func (b *numberBucket) Increment(i float64) {
	b.Add(i)
}

func (b *numberBucket) UpdateMax(n float64) {
	for {
		old := b.Load()
		if n > old {
			if b.CAS(old, n) {
				return
			}
		} else {
			return
		}
	}
}

// Number tracks a numberBucket over a bounded number of
// time buckets. Currently the buckets are one second long and only the last N seconds are kept.
type Number struct {
	buckets []numberBucket
	bCap    int64

	currentBucket atomic.Value //*Bucket
	currentIndex  int64

	switchTicker *time.Ticker
	stop         bool
}

// NewNumber initializes a RollingNumber struct.
func NewNumber(N int) *Number {
	num := &Number{
		buckets: make([]numberBucket, N),
		bCap:    int64(N),
	}

	num.currentBucket.Store(&num.buckets[0])
	num.currentIndex = 0

	num.switchTicker = time.NewTicker(time.Millisecond * 200) //200ms
	go func() {
		for t := range num.switchTicker.C {
			if num.stop {
				return
			}

			if i := t.Unix() % num.bCap; i != num.currentIndex {
				num.buckets[i].Reset()
				num.currentIndex = i
				num.currentBucket.Store(&num.buckets[num.currentIndex])
			}
		}
	}()

	return num
}

// Stop destroy resource
func (num *Number) Stop() {
	num.stop = true
}

// Increment increments the number in current timeBucket.
func (num *Number) Increment(i float64) {
	if i == 0 {
		return
	}

	num.currentBucket.Load().(*numberBucket).Increment(i)
}

// UpdateMax updates the maximum value in the current bucket.
func (num *Number) UpdateMax(n float64) {
	num.currentBucket.Load().(*numberBucket).UpdateMax(n)
}

// Sum sums the values over the buckets in the last N seconds.
func (num *Number) Sum() float64 {
	sum := float64(0)

	for _, bucket := range num.buckets {
		sum += bucket.Load()
	}

	return sum
}

// Max returns the maximum value seen in the last N seconds.
func (num *Number) Max() float64 {
	var max float64

	for _, bucket := range num.buckets {
		if n := bucket.Load(); n > max {
			max = n
		}
	}

	return max
}

// Avg returns the average value seen in the last N seconds.
func (num *Number) Avg() float64 {
	return num.Sum() / float64(num.bCap)
}
