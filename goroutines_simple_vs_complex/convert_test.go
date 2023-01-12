package goroutines_simple_vs_complex

import (
	"fmt"
	"testing"

	"github.com/go-faker/faker/v4"

	"time"
)

func Benchmark_SimpleConvertPineApplesToSafety(b *testing.B) {
	for _, n := range []int{500, 1000, 2000, 5000, 10000} {
		b.Run(fmt.Sprintf("Benchmark_SimpleConvertPineApplesToSafety-%d", n), func(b *testing.B) {
			pineApples := make([]Pineapple, n)

			var pine Pineapple

			for i := 0; i < n; i++ {
				_ = faker.FakeData(&pine)
				pine.Created = time.Now().AddDate(0, 0, -i)
				pine.ID = uint(i)
				pine.IsAlive = true
				pineApples[i] = pine
			}
			for i := 0; i < b.N; i++ {
				SimpleConvertPineApplesToSafety(pineApples)
			}
		})
	}
}

func Benchmark_GoroutinesConvertPineApplesToSafety(b *testing.B) {
	for _, n := range []int{500, 1000, 2000, 5000, 10000} {
		b.Run(fmt.Sprintf("Benchmark_GoroutinesConvertPineApplesToSafety-%d", n), func(b *testing.B) {
			pineApples := make([]Pineapple, n)

			var pine Pineapple

			for i := 0; i < n; i++ {
				_ = faker.FakeData(&pine)
				pine.Created = time.Now().AddDate(0, 0, -i)
				pine.ID = uint(i)
				pine.IsAlive = true
				pineApples[i] = pine
			}
			for i := 0; i < b.N; i++ {
				GoroutinesConvertPineApplesToSafety(pineApples)
			}
		})
	}
}

func Benchmark_NoMutexGoroutinesConvertPineApplesToSafety(b *testing.B) {
	for _, n := range []int{500, 1000, 2000, 5000, 10000} {
		b.Run(fmt.Sprintf("Benchmark_NoMutexGoroutinesConvertPineApplesToSafety-%d", n), func(b *testing.B) {
			pineApples := make([]Pineapple, n)

			var pine Pineapple

			for i := 0; i < n; i++ {
				_ = faker.FakeData(&pine)
				pine.Created = time.Now().AddDate(0, 0, -i)
				pine.ID = uint(i)
				pine.IsAlive = true
				pineApples[i] = pine
			}
			for i := 0; i < b.N; i++ {
				GoroutinesNoMutexConvertPineApplesToSafety(pineApples)
			}
		})
	}
}

func Test_GoroutinesConvertPineApplesToSafety(t *testing.T) {
	t.Run("TestGoroutinesConvertPineApplesToSafety", func(t *testing.T) {
		pineApples := make([]Pineapple, 10000)

		var pine Pineapple

		for i := 0; i < 10000; i++ {
			_ = faker.FakeData(&pine)
			pine.Created = time.Now().AddDate(0, 0, -i)
			pine.ID = uint(i)
			pine.IsAlive = true
			pineApples[i] = pine
		}

		got := GoroutinesConvertPineApplesToSafety(pineApples)
		// Check order by Created field
		for i := 0; i < len(got)-1; i++ {
			if got[i].ID > got[i+1].ID {
				t.Errorf("GoroutinesConvertPineApplesToSafety() = %v, want %v", got, pineApples)
			}
		}
	})
}

func Test_GouroutinesNoMutexConvertPineApplesToSafety(t *testing.T) {
	t.Run("TestGouroutinesNoMutexConvertPineApplesToSafety", func(t *testing.T) {
		pineApples := make([]Pineapple, 10000)

		var pine Pineapple

		for i := 0; i < 10000; i++ {
			_ = faker.FakeData(&pine)
			pine.Created = time.Now().AddDate(0, 0, -i)
			pine.ID = uint(i)
			pine.IsAlive = true
			pineApples[i] = pine
		}

		got := GoroutinesNoMutexConvertPineApplesToSafety(pineApples)
		// Check order by Created field
		for i := 0; i < len(got)-1; i++ {
			if got[i].ID > got[i+1].ID {
				t.Errorf("TestGouroutinesNoMutexConvertPineApplesToSafety() = %v, want %v", got, pineApples)
			}
		}
	})
}

func Test_SimpleConvertPineApplesToSafety(t *testing.T) {
	t.Run("TestSimpleConvertPineApplesToSafety", func(t *testing.T) {
		pineApples := make([]Pineapple, 10000)

		var pine Pineapple

		for i := 0; i < 10000; i++ {
			_ = faker.FakeData(&pine)
			pine.Created = time.Now().AddDate(0, 0, -i)
			pine.ID = uint(i)
			pine.IsAlive = true
			pineApples[i] = pine
		}

		got := SimpleConvertPineApplesToSafety(pineApples)
		// Check order by Created field
		for i := 0; i < len(got)-1; i++ {
			if got[i].ID > got[i+1].ID {
				t.Errorf("TestSimpleConvertPineApplesToSafety() = %v, want %v", got, pineApples)
			}
		}
	})
}
