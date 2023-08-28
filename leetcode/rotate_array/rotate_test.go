package rotate_array

import (
	"fmt"
	"reflect"
	"testing"
)

func largeArray() []int {
	var array = make([]int, 100)
	for i := 0; i < 100; i++ {
		array[i] = i
	}

	return array
}

func largeArrayReversed() []int {
	var array = make([]int, 100)
	for i := 1; i < 100; i++ {
		array[i] = i - 1
	}

	array[0] = 99

	return array
}

func testFramework(t *testing.T, rotateFunction func(nums []int, k int)) {
	reverseTests := []struct {
		input    []int
		k        int
		expected []int
		name     string
	}{
		{
			input:    []int{1, 2, 3, 4, 5, 6, 7},
			k:        3,
			expected: []int{5, 6, 7, 1, 2, 3, 4},
			name:     "7 elements, k=3",
		},
		{
			input:    []int{1, 2, 3, 4, 5, 6, 7},
			k:        10,
			expected: []int{5, 6, 7, 1, 2, 3, 4},
			name:     "7 elements, k=10",
		},
		{
			input:    []int{1, 2, 3, 4, 5, 6, 7},
			k:        0,
			expected: []int{1, 2, 3, 4, 5, 6, 7},
			name:     "7 elements, k=0",
		},
		{
			input:    []int{1, 2, 3, 4, 5, 6, 7},
			k:        1,
			expected: []int{7, 1, 2, 3, 4, 5, 6},
			name:     "7 elements, k=1",
		},
		{
			// large array
			input:    largeArray(),
			k:        100000,
			expected: largeArray(),
			name:     "large array, k=100000",
		},
		{
			// large array
			input:    largeArray(),
			k:        0,
			expected: largeArray(),
			name:     "large array, k=0",
		},
		{
			// large array
			input:    largeArray(),
			k:        1,
			expected: largeArrayReversed(),
			name:     "large array, k=1",
		},
	}

	for _, test := range reverseTests {
		t.Run(test.name, func(t *testing.T) {
			rotateFunction(test.input, test.k)
			if !reflect.DeepEqual(test.input, test.expected) {
				t.Errorf("test %s failed", test.name)
				t.Errorf("actual %v expected %v", test.input, test.expected)
			}
		})
	}
}

func TestRotateCopy(t *testing.T) {
	testFramework(t, rotateCopy)
}

func TestRotateInPlace(t *testing.T) {
	testFramework(t, rotateInPlace)
}

func TestRotateGnGn(t *testing.T) {
	testFramework(t, rotateGnGn)
}

func benchmarkFramework(b *testing.B, rotateFunction func(nums []int, k int)) {
	arrays := [][]int{RandomArray(100, 0, 100),
		RandomArray(1000, 0, 1000),
		RandomArray(10000, 0, 10000),
		RandomArray(100000, 0, 100000),
	}
	b.ResetTimer()
	for _, array := range arrays {
		b.Run(fmt.Sprintf("%d", len(array)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				rotateFunction(array, random(0, len(array)))
			}
		})
	}
}

func BenchmarkRotateCopy(b *testing.B) {
	benchmarkFramework(b, rotateCopy)
}

func BenchmarkRotateInPlace(b *testing.B) {
	benchmarkFramework(b, rotateInPlace)
}

func BenchmarkRotateGnGn(b *testing.B) {
	benchmarkFramework(b, rotateGnGn)
}
