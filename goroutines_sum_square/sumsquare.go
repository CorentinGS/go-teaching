package goroutines_sum_square

import (
	"math/rand"
	"reflect"
	"runtime"
	"sync"
	"unsafe"
)

func simpleParallelSumSquare(items []int) int {
	const chunkSize = 10000

	// Divide the items into chunks
	chunks := make([][]int, 0)
	for i := 0; i < len(items); i += chunkSize {
		end := i + chunkSize
		if end > len(items) {
			end = len(items)
		}
		chunks = append(chunks, items[i:end])
	}

	// Create a goroutine for each chunk
	wg := sync.WaitGroup{}
	resultChan := make(chan int, len(chunks))

	for _, chunk := range chunks {
		wg.Add(1)
		go func(chunk []int) {
			resultChan <- simpleSumSquare(chunk)
			wg.Done()
		}(chunk)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(resultChan)

	// Sum the results
	total := 0
	for partialSum := range resultChan {
		total += partialSum
	}

	return total
}

func simpleSumSquare(items []int) int {
	total := 0
	for i := 0; i < len(items); i++ {
		total += items[i] * items[i]
	}
	return total
}

func sumSquare(items []int) int {
	number := make(chan int)
	response := make(chan int)

	var wg sync.WaitGroup

	total := 0

	for _, item := range items {
		wg.Add(1)
		go func(item int) {
			defer wg.Done()
			sum1 := <-number
			sum1 = sum1 * sum1
			response <- sum1
		}(item)
		number <- item
		total += <-response
	}

	defer close(number)
	defer close(response)

	wg.Wait()

	return total
}

func parallelSumSquare(items []int) int {
	if len(items) < 10000 { // Threshold for small slices
		return simpleSumSquare(items)
	}

	totalCPU := runtime.NumCPU()
	chunkSize := (len(items) + totalCPU - 1) / totalCPU
	resultChan := make(chan int, totalCPU)
	wg := sync.WaitGroup{}

	for i := 0; i < totalCPU; i++ {
		start := i * chunkSize
		end := (i + 1) * chunkSize
		if end > len(items) {
			end = len(items)
		}

		wg.Add(1)
		go func(start, end int) {
			total := 0
			for i := start; i < end; i++ {
				item := items[i]
				total += item * item
			}
			resultChan <- total
			wg.Done()
		}(start, end)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	total := 0
	for partialSum := range resultChan {
		total += partialSum
	}

	return total
}

func unsafeSlice(items []int, start, end int) []int {
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&items))
	header.Data = uintptr(unsafe.Pointer(&items[0]))
	header.Len = end - start
	header.Cap = end - start

	return *(*[]int)(unsafe.Pointer(&header))
}

func unsafeParallelSumSquare(items []int) int {
	if len(items) < 10000 { // Threshold for small slices
		return simpleSumSquare(items)
	}

	totalCPU := runtime.NumCPU()
	chunkSize := (len(items) + totalCPU - 1) / totalCPU
	resultChan := make(chan int, totalCPU)
	wg := sync.WaitGroup{}

	for i := 0; i < totalCPU; i++ {
		start := i * chunkSize
		end := (i + 1) * chunkSize
		if end > len(items) {
			end = len(items)
		}

		wg.Add(1)
		go func(start, end int) {
			partialItems := unsafeSlice(items, start, end)
			partialSum := 0
			for _, item := range partialItems {
				partialSum += item * item
			}
			resultChan <- partialSum
			wg.Done()
		}(start, end)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	total := 0
	for partialSum := range resultChan {
		total += partialSum
	}

	return total
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
