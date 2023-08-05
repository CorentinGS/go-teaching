package goroutines_sum_square

import (
	"fmt"
	"reflect"
	"testing"
)

func hundredFirstIntegers() []int {
	hundredFirstIntegers := make([]int, 100)
	for i := 0; i < 100; i++ {
		hundredFirstIntegers[i] = i
	}
	return hundredFirstIntegers
}

func thousandFirstIntegers() []int {
	thousandFirstIntegers := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		thousandFirstIntegers[i] = i
	}
	return thousandFirstIntegers
}

func testFramework(t *testing.T, sumFunction func([]int) int) {
	sortTests := []struct {
		input    []int
		expected int
		name     string
	}{
		//10 first integers
		{
			input:    []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			expected: 285,
			name:     "10 first integers",
		},
		// 100 first integers
		{
			input:    hundredFirstIntegers(),
			expected: 328350,
			name:     "100 first integers",
		},
		// 1000 first integers
		{
			input:    thousandFirstIntegers(),
			expected: 332833500,
			name:     "1000 first integers",
		},
	}
	for _, test := range sortTests {
		t.Run(test.name, func(t *testing.T) {
			actual := sumFunction(test.input)
			sorted := reflect.DeepEqual(actual, test.expected)
			if !sorted {
				t.Errorf("test %s failed", test.name)
				t.Errorf("actual %v expected %v", actual, test.expected)
			}
		})
	}
}

func benchmarkFramework(b *testing.B, sumFunction func([]int) int) {
	sizes := [][]int{RandomArray(100, 0, 100),
		RandomArray(1000, 0, 1000),
		RandomArray(10000, 0, 10000),
		RandomArray(100000, 0, 100000),
		RandomArray(1000000, 0, 1000000),
	}
	b.ResetTimer()
	for _, size := range sizes {
		b.Run(fmt.Sprintf("%d", len(size)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				sumFunction(size)
			}
		})
	}
}

func BenchmarkSumSquare(b *testing.B) {
	benchmarkFramework(b, sumSquare)
}

func BenchmarkSimpleSumSquare(b *testing.B) {
	benchmarkFramework(b, simpleSumSquare)
}

func TestSumSquare(t *testing.T) {
	testFramework(t, sumSquare)
}

func TestSimpleSumSquare(t *testing.T) {
	testFramework(t, simpleSumSquare)
}

func TestSumSquareWithGoroutines(t *testing.T) {
	testFramework(t, parallelSumSquare)
}

func BenchmarkSumSquareWithGoroutines(b *testing.B) {
	benchmarkFramework(b, parallelSumSquare)
}

func TestSimpleParallelSumSquare(t *testing.T) {
	testFramework(t, simpleParallelSumSquare)
}

func BenchmarkSimpleParallelSumSquare(b *testing.B) {
	benchmarkFramework(b, simpleParallelSumSquare)
}

func TestUnsafeParallelSumSquare(t *testing.T) {
	testFramework(t, unsafeParallelSumSquare)
}

func BenchmarkUnsafeParallelSumSquare(b *testing.B) {
	benchmarkFramework(b, unsafeParallelSumSquare)
}
