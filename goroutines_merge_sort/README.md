# Parallel Merge Sort vs Simple Merge Sort

This is a simple example of how to use goroutines to parallelize a merge sort algorithm.
We compare the performance of a simple merge sort algorithm with a parallel merge sort algorithm that uses goroutines.

## The Merge Sort Algorithm

The merge sort algorithm is a divide and conquer algorithm that recursively splits the input array into two halves, 
sorts each half, and then merges the two sorted halves into a single sorted array.

To speed up the merge sort algorithm, we use insertion sort for small subarrays (less than 12 elements). 

The implementation of the algorithm uses generics to allow sorting of any type of numbers.

```go
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
```

## The Merge Sort Algorithm with Goroutines

The parallel merge sort algorithm uses goroutines to sort the two halves of the input array in parallel.

To prevent the creation of too many goroutines, we use a threshold to determine when to use goroutines. 
If the size of the input array is less than the threshold, 
we use a simple merge sort algorithm instead of a parallel merge sort algorithm.

Here we use a threshold of 512 elements. We can benchmark the performance of the algorithm with
different thresholds to find the optimal threshold, but we will not do that in this example.

```go
// ParallelMerge Perform merge sort on a slice using goroutines
func ParallelMerge[T Number](items []T) []T {
	if len(items) < 2 {
		return items
	}

	// Use a simple merge sort algorithm if the size of the input array is less than the threshold
	if len(items) < 512 {
		return MergeSort(items)
	}

	// Create the wait group to wait for the goroutines to finish
	var wg sync.WaitGroup
	wg.Add(1)

	var middle = len(items) / 2  // Find the middle index of the input array
	var a []T                   // Create a slice to hold the first half of the input array
	go func() {                // Create a goroutine to sort the first half of the input array
		defer wg.Done()       // Decrement the wait group counter when the goroutine finishes
		a = ParallelMerge(items[:middle]) // Sort the first half of the input array
	}()
	var b = ParallelMerge(items[middle:]) // Sort the second half of the input array

	wg.Wait() // Wait for the goroutine to finish
	return merge(a, b) // Merge the two sorted halves
}
```

## Benchmarking the Merge Sort Algorithms

Now we can benchmark our merge sort algorithms to compare their performance.

To benchmark the algorithms, we use arrays of different sizes and measure the time it takes to sort the arrays.

To run the benchmark we use the following command:
```bash
go test -bench=. -benchtime 5s > benchmark.txt && benchstat benchmark.txt   
```

Benchstat is a tool that can be used to compare the results of benchmarks.

The results of the benchmark are as follows:

```bash
name                                   time/op
Mergesort/1000                     14.6µs ± 0%
Mergesort/10000                     473µs ± 0%
Mergesort/100000                   6.33ms ± 0%
Mergesort/1000000                  87.4ms ± 0%

MergesortWithGoroutines/1000       18.7µs ± 0%
MergesortWithGoroutines/10000       217µs ± 0%
MergesortWithGoroutines/100000     2.71ms ± 0%
MergesortWithGoroutines/1000000    29.0ms ± 0%
```

As we can see, the parallel merge sort algorithm is much faster than the simple merge sort algorithm for large arrays.

## Why is the Parallel Merge Sort Algorithm Faster? 

The parallel merge sort algorithm is faster because it uses goroutines to sort the two halves of the input array in parallel. 
The simple merge sort algorithm sorts the two halves of the input array sequentially.

As discussed in the previous article, the cost of creating a goroutine can be high. 
So we should use a parallel merge sort algorithm only when the size of the input array is large enough to justify the cost of creating goroutines. 

## Conclusion

In this example, we have seen how to use goroutines to parallelize a merge sort algorithm. We have also benchmarked the performance of the merge sort algorithms to compare their performance. 
Goroutines are a powerful tool that should be used only when the cost of creating goroutines is justified by the performance improvement. 

In this example, the code isn't more complex when using goroutines therefore it's worth using them.

## Code

The code for this article can be found on [GitHub](https://github.com/CorentinGS/go-teaching).

## License

This article and the corresponding code is licensed under the [ISC License](https://github.com/CorentinGS/go-teaching/blob/main/LICENSE).
If you want to distribute or cite this article and the corresponding code you can do it as long as you follow the ISC License and this steps:
- Cite the [GitHub repo](https://github.com/CorentinGS/go-teaching/)
- Include the [License](https://github.com/CorentinGS/go-teaching/blob/main/LICENSE)
- Email me or DM me on social networks to tell me how you reused my article so that I mention you on GitHub

## Contributing

If you want to contribute to the article or to the code itself, please feel free to open a pull request.

You can also [contact](https://corentings.vercel.app/links) me.

If you have any questions regarding this article, please open an issue.

## Disclaimer

I'm not a professional gopher and my publications shouldn't be taken as entirely true nor accurate/perfect.
This article is a personal work for further research and development purposes only.
I hope that it will be useful for others, but it shouldn't be seen as a guide. 

## References

- [Go](https://github.com/CorentinGS/Go/blob/master/sort/mergesort.go)