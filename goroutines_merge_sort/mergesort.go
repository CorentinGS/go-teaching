package goroutines_merge_sort

import "sync"

const K = 32

func merge[T Number](a []T, b []T) []T {

	var r = make([]T, len(a)+len(b))
	var i = 0
	var j = 0

	for i < len(a) && j < len(b) {

		if a[i] <= b[j] {
			r[i+j] = a[i]
			i++
		} else {
			r[i+j] = b[j]
			j++
		}

	}

	for i < len(a) {
		r[i+j] = a[i]
		i++
	}
	for j < len(b) {
		r[i+j] = b[j]
		j++
	}

	return r

}

func MergeSort[T Number](items []T) []T {
	size := len(items)
	if size < 2 {
		return items
	}

	if size < K {
		return Insertionsort(items)
	}

	middle := size / 2
	var a = MergeSort(items[:middle])
	var b = MergeSort(items[middle:])

	return merge(a, b)
}

// ParallelMerge Perform merge sort on a slice using goroutines
func ParallelMerge[T Number](items []T) []T {
	if len(items) < 2 {
		return items
	}

	if len(items) < 512 {
		return MergeSort(items)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	var middle = len(items) / 2
	var a []T
	go func() {
		defer wg.Done()
		a = ParallelMerge(items[:middle])
	}()
	var b = ParallelMerge(items[middle:])

	wg.Wait()
	return merge(a, b)
}
