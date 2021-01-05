package main

import (
	"fmt"
	"math/rand"
	"time"

	"go.uber.org/atomic"
)

func main() {
	// wantFailureRate := 12 //12%
	// wantFailureRate := 50 //50%
	wantFailureRate := 99 //99%

	rand.Seed(time.Now().Unix())

	sw := NewSlidingWindow(10)
	defer sw.Stop()

	go func() {
		for range time.Tick(time.Millisecond * 500) {
			fmt.Printf("%d: FailureRate=%.4f\n", time.Now().Unix(), sw.CalcFailureRate())
		}
	}()

	for range time.Tick(time.Millisecond * 50) {
		if rand.Intn(100) < wantFailureRate {
			//Failure
			sw.AddFailure(1)
		} else {
			//Success
			sw.AddSuccess(1)
		}
	}
}

//===================Bucket==============

// Bucket struct contain Success/Failure
type Bucket [2]atomic.Int32

const (
	// BucketSuccessIndex define
	BucketSuccessIndex = iota
	// BucketFailureIndex define
	BucketFailureIndex
)

// AddSuccess add Success counter
func (b *Bucket) AddSuccess(n int) {
	b[BucketSuccessIndex].Add(int32(n))
}

// AddFailure add Failure counter
func (b *Bucket) AddFailure(n int) {
	b[BucketFailureIndex].Add(int32(n))
}

// GetAll get Success and Failure counter
func (b *Bucket) GetAll() (success, failure int) {
	return int(b[BucketSuccessIndex].Load()), int(b[BucketFailureIndex].Load())
}

// Reset Bucket count to zero
func (b *Bucket) Reset() {
	b[0].Store(0)
	b[1].Store(0)
}

//===================Bucket==============

//===================Sliding Window==============

// SlidingWindow struct
type SlidingWindow struct {
	bucket []Bucket
	bCap   int64

	currentBucket atomic.Value //*Bucket
	currentIndex  int64

	switchTicker *time.Ticker
	stop         bool
}

// NewSlidingWindow return a new SlidingWindow with bucketNum
func NewSlidingWindow(bucketNum int) *SlidingWindow {
	sw := &SlidingWindow{
		bucket: make([]Bucket, bucketNum),
		bCap:   int64(bucketNum),
	}

	sw.currentBucket.Store(&sw.bucket[0])
	sw.currentIndex = 0

	sw.switchTicker = time.NewTicker(time.Millisecond * 200) //200ms
	go func() {
		for t := range sw.switchTicker.C {
			if sw.stop {
				return
			}

			if i := t.Unix() % sw.bCap; i != sw.currentIndex {
				sw.bucket[i].Reset()
				sw.currentIndex = i
				sw.currentBucket.Store(&sw.bucket[sw.currentIndex])
			}
		}
	}()

	return sw
}

// Stop destroy resource
func (sw *SlidingWindow) Stop() {
	sw.stop = true
}

// AddSuccess add n to success counter
func (sw *SlidingWindow) AddSuccess(n int) {
	sw.currentBucket.Load().(*Bucket).AddSuccess(n)
}

// AddFailure add n to failure counter
func (sw *SlidingWindow) AddFailure(n int) {
	sw.currentBucket.Load().(*Bucket).AddFailure(n)
}

// CalcFailureRate calculate current Failure Rate
func (sw *SlidingWindow) CalcFailureRate() float64 {
	success, failure := 0, 0
	for _, b := range sw.bucket {
		s, f := b.GetAll()
		success += s
		failure += f
	}

	return float64(failure) / float64(success+failure)
}

//===================Sliding Window==============
