package main

import (
	"testing"

	"github.com/afex/hystrix-go/hystrix/rolling"
)

func BenchmarkAddFailure(b *testing.B) {
	sw := NewSlidingWindow(10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sw.AddFailure(1)
	}
	b.StopTimer()
}

func BenchmarkAddSuccess(b *testing.B) {
	sw := NewSlidingWindow(10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sw.AddSuccess(1)
	}
	b.StopTimer()
}

func BenchmarkHystrixRollingNumberIncrement(b *testing.B) {
	n := rolling.NewNumber()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		n.Increment(1)
	}
}

func BenchmarkHystrixRollingNumberUpdateMax(b *testing.B) {
	n := rolling.NewNumber()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		n.UpdateMax(float64(i))
	}
}
