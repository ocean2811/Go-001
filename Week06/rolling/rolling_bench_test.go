package rolling

import (
	"testing"

	"github.com/afex/hystrix-go/hystrix/rolling"
)

const BucketNum = 10

func BenchmarkRollingNumberIncrement(b *testing.B) {
	n := NewNumber(BucketNum)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		n.Increment(1)
	}
}

func BenchmarkRollingNumberUpdateMax(b *testing.B) {
	n := NewNumber(BucketNum)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		n.UpdateMax(float64(i))
	}
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
