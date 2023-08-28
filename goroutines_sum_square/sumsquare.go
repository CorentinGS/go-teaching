package goroutines_sum_square

import (
	"math/rand"
	"reflect"
	"runtime"
	"sync"
	"unsafe"
)

func simpleParallelSumSquare(items []int) int {
	if len(items) <= 10000 { // Threshold for small slices
		return simpleSumSquare(items) // Use the simpleSumSquare function
	}

	const chunkSize = 10000

	// Divide the items into chunks
	chunks := make([][]int, 0)
	for i := 0; i < len(items); i += chunkSize {
		end := i + chunkSize // end index for the chunk
		if end > len(items) {
			end = len(items) // last chunk may be smaller than chunkSize
		}
		chunks = append(chunks, items[i:end]) // append the chunk to the chunks slice
	}

	// Create a goroutine for each chunk
	wg := sync.WaitGroup{}
	resultChan := make(chan int, len(chunks)) // channel for receiving results

	for _, chunk := range chunks { // iterate over the chunks
		wg.Add(1)              // increment the wait group counter
		go func(chunk []int) { // create a goroutine
			resultChan <- simpleSumSquare(chunk) // send the result to the result channel
			wg.Done()                            // decrement the wait group counter when the goroutine finishes
		}(chunk) // pass the chunk to the goroutine
	}

	wg.Wait()         // Wait for all goroutines to finish
	close(resultChan) // close the result channel

	// Sum the results
	total := 0
	for partialSum := range resultChan {
		total += partialSum
	}

	return total // return the total sum
}

func optimizedParallelSumSquare(items []int) int {
	if len(items) <= 10000 {
		return simpleSumSquare(items)
	}

	const chunkSize = 10000

	// Divide the items into chunks without creating a slice of slices
	// Instead, we create a slice of start and end indices
	chunkIndices := make([]struct{ start, end int }, 0)
	for i := 0; i < len(items); i += chunkSize {
		end := i + chunkSize
		if end > len(items) {
			end = len(items)
		}
		chunkIndices = append(chunkIndices, struct{ start, end int }{i, end})
	}

	wg := sync.WaitGroup{}
	resultChan := make(chan int, len(chunkIndices))

	for _, indices := range chunkIndices {
		wg.Add(1)
		go func(start, end int) {
			resultChan <- simpleSumSquare(items[start:end])
			wg.Done()
		}(indices.start, indices.end)
	}

	wg.Wait()
	close(resultChan)

	total := 0
	for partialSum := range resultChan {
		total += partialSum
	}

	return total
}

func simpleSumSquare(items []int) int {
	total := 0 // total sum
	for i := 0; i < len(items); i++ {
		total += items[i] * items[i] // square the item and add it to the total
	}
	return total // return the total sum
}

func sumSquare(items []int) int {
	number := make(chan int)   // channel for sending numbers
	response := make(chan int) // channel for receiving responses

	var wg sync.WaitGroup // wait group for waiting for all goroutines to finish

	total := 0 // total sum

	// Create a goroutine for each item in the slice
	for _, item := range items {
		wg.Add(1)           // increment the wait group counter
		go func(item int) { // create a goroutine
			defer wg.Done()    // decrement the wait group counter when the goroutine finishes
			sum1 := <-number   // receive a number from the number channel
			sum1 = sum1 * sum1 // square the number
			response <- sum1   // send the result to the response channel
		}(item) // pass the item to the goroutine
		number <- item      // send the item to the number channel
		total += <-response // receive the result from the response channel
	}

	defer close(number)   // close the number channel
	defer close(response) // close the response channel

	wg.Wait() // wait for all goroutines to finish

	return total // return the total sum
}

func parallelSumSquare(items []int) int {
	if len(items) <= 10000 { // Threshold for small slices
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
