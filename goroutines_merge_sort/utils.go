package goroutines_merge_sort

import "math/rand"

func Insertionsort[N Number](array []N) []N {
	for i := 1; i < len(array); i++ {
		for j := i; j > 0 && array[j] < array[j-1]; j-- {
			swap(&array, j, j-1)
		}
	}
	return array
}

func min[T Number](values ...T) T {
	var min = values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
	}
	return min
}

func swap[T any](array *[]T, i int, j int) {
	var tmp = (*array)[i]
	(*array)[i] = (*array)[j]
	(*array)[j] = tmp
}

func random(min int, max int) int {
	return min + rand.Intn(int(max-min))
}

func RandomArray(size int, min int, max int) []int {
	var array = make([]int, size)
	for i := 0; i < size; i++ {
		array[i] = random(min, max)
	}
	return array
}
