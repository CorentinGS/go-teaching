package goroutines_merge_sort_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/corentings/goTeaching/goroutines_merge_sort"
)

func testFramework(t *testing.T, sortingFunction func([]int) []int) {
	sortTests := []struct {
		input    []int
		expected []int
		name     string
	}{
		//Sorted slice
		{
			input:    []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			expected: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			name:     "Sorted Unsigned",
		},
		//Reversed slice
		{
			input:    []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
			expected: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			name:     "Reversed Unsigned",
		},
		//Sorted slice
		{
			input:    []int{-10, -9, -8, -7, -6, -5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			expected: []int{-10, -9, -8, -7, -6, -5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			name:     "Sorted Signed",
		},
		//Reversed slice
		{
			input:    []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0, -1, -2, -3, -4, -5, -6, -7, -8, -9, -10},
			expected: []int{-10, -9, -8, -7, -6, -5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			name:     "Reversed Signed",
		},
		//Reversed slice, even length
		{
			input:    []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1, -1, -2, -3, -4, -5, -6, -7, -8, -9, -10},
			expected: []int{-10, -9, -8, -7, -6, -5, -4, -3, -2, -1, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			name:     "Reversed Signed #2",
		},
		//Random order with repetitions
		{
			input:    []int{-5, 7, 4, -2, 6, 5, 8, 3, 2, -7, -1, 0, -3, 9, -6, -4, 10, 9, 1, -8, -9, -10},
			expected: []int{-10, -9, -8, -7, -6, -5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 9, 10},
			name:     "Random order Signed",
		},
		//Single-entry slice
		{
			input:    []int{1},
			expected: []int{1},
			name:     "Singleton",
		},
		// Empty slice
		{
			input:    []int{},
			expected: []int{},
			name:     "Empty Slice",
		},
	}
	for _, test := range sortTests {
		t.Run(test.name, func(t *testing.T) {
			actual := sortingFunction(test.input)
			sorted := reflect.DeepEqual(actual, test.expected)
			if !sorted {
				t.Errorf("test %s failed", test.name)
				t.Errorf("actual %v expected %v", actual, test.expected)
			}
		})
	}
}

func TestInsertionsort(t *testing.T) {
	testFramework(t, goroutines_merge_sort.Insertionsort[int])
}

func TestMergesort(t *testing.T) {
	testFramework(t, goroutines_merge_sort.MergeSort[int])
}

func TestMergesortWithGoroutines(t *testing.T) {
	testFramework(t, goroutines_merge_sort.ParallelMerge[int])
}

func benchmarkFramework(b *testing.B, sortingFunction func([]int) []int) {
	sizes := [][]int{goroutines_merge_sort.RandomArray(100, 0, 100),
		goroutines_merge_sort.RandomArray(1000, 0, 1000),
		goroutines_merge_sort.RandomArray(10000, 0, 10000),
		goroutines_merge_sort.RandomArray(100000, 0, 100000),
		goroutines_merge_sort.RandomArray(1000000, 0, 1000000),
	}
	b.ResetTimer()
	for _, size := range sizes {
		b.Run(fmt.Sprintf("%d", len(size)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				sortingFunction(size)
			}
		})
	}
}

func BenchmarkMergesort(b *testing.B) {
	benchmarkFramework(b, goroutines_merge_sort.MergeSort[int])
}

func BenchmarkMergesortWithGoroutines(b *testing.B) {
	benchmarkFramework(b, goroutines_merge_sort.ParallelMerge[int])
}
