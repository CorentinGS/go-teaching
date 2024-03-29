# Goroutine vs Simple function

This is a simple example of why goroutines might be overkill for some tasks and less efficient than a simple function.

### Structures 

We got a simple structure that contains sensitive information that we don't want to be exposed to the outside world.
Therefore, we created a second structure that hides the sensitive information and only exposes the information we want to be public.

```go
// Pineapple is a struct that represents a database object with sensitive data that should be hidden
type Pineapple struct {
	Paro       string `faker:"name"`
	Turkey     string `faker:"name"`
	Banana     string `faker:"name"`
	Age        int    `faker:"number"`
	Size       int    `faker:"number"`
	IsAlive    bool
	ID         uint
	SecretCode []byte
	Created    time.Time
	Updated    time.Time
}

// SafePineApple is a struct that represents a Pineapple object without sensitive data
type SafePineApple struct {
    Paro    string
    Turkey  string
    Banana  string
    IsAlive bool
    Age     int
    ID      uint
}
```

### Conversion 

We need to convert our Pineapple object to a SafePineApple object.
We can do this by creating a method on the Pineapple struct that returns a SafePineApple object.

```go
// ToSafePineApple converts a Pineapple object to a SafePineApple object
func (p *Pineapple) ToSafePineApple() SafePineApple {
	return SafePineApple{
		Paro:    p.Paro,
		Turkey:  p.Turkey,
		Banana:  p.Banana,
		IsAlive: p.IsAlive,
		Age:     p.Age,
		ID:      p.ID,
	}
}
```

## Use case

In our use case we have an array of Pineapple objects coming from our database that we want to convert to SafePineApple objects and store them in a new array.
The order of the objects in the array should be the same as the original array as it has already been sorted by a sql query.


### Simple function

We can do this by creating a simple function that takes an array of Pineapple objects and returns an array of SafePineApple objects.

```go
// SimpleConvertPineApplesToSafety converts an array of Pineapple objects to an array of SafePineApple objects
func SimpleConvertPineApplesToSafety(pineapples []Pineapple) []SafePineApple {
	safePineApples := make([]SafePineApple, len(pineapples))

	for idx, pineapple := range pineapples {
		safePineApples[idx] = pineapple.ToSafePineApple()
	}

	return safePineApples
}
```

This function is very simple and easy to understand. We loop through the array of Pineapple objects and convert them to SafePineApple objects. 

### Goroutine without mutex

We can do this by using goroutines to work on the array concurrently and store the results in a new array.

```go
func GoroutinesNoMutexConvertPineApplesToSafety(pineapples []Pineapple) []SafePineApple {
	// Create a slice to store the SafePineApples
	safePineApples := make([]SafePineApple, len(pineapples)/2, len(pineapples))
	safePineApples2 := make([]SafePineApple, len(pineapples)/2)

	var wg sync.WaitGroup // Create a WaitGroup to wait for all goroutines to finish
	wg.Add(1)            // Add 1 to the WaitGroup

	// Create a goroutine to convert the first half of the Pineapple objects
	go func(chunk []Pineapple) {
		defer wg.Done() // Decrement the WaitGroup when the goroutine is done
		for idx, pineapple := range chunk { // Loop through the chunk of Pineapple objects
			safePineApples[idx] = pineapple.ToSafePineApple() // Convert the Pineapple object to a SafePineApple object
		}
	}(pineapples[:len(pineapples)/2]) // Pass the first half of the Pineapple objects to the goroutine

	// Convert the second half of the Pineapple objects in the main thread 
	for idx, pineapple := range pineapples[len(pineapples)/2:] {
		safePineApples2[idx] = pineapple.ToSafePineApple()
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Group both pineapples
	safePineApples = append(safePineApples, safePineApples2...)

	// Return the SafePineApples
	return safePineApples
}
```

This function is a bit more complex than the simple function. We use a goroutine to convert the first half while the other half is handled by the main thread.
We use a WaitGroup to wait for the goroutine to finish before returning the results. 

### Goroutine with mutex

We can also add mutexes to the goroutine to make it thread safe.

```go
func GoroutinesConvertPineApplesToSafety(pineapples []Pineapple) []SafePineApple {
	// Create a slice to store the SafePineApples
	safePineApples := make([]SafePineApple, len(pineapples))

	// Split the offers into chunks
	chunks := [][]Pineapple{pineapples[:len(pineapples)/2], pineapples[len(pineapples)/2:]}
	
	mutex := sync.Mutex{} // Create a mutex to lock the slice when writing to it

	var wg sync.WaitGroup // Create a WaitGroup to wait for all goroutines to finish
	wg.Add(1)           // Add 1 to the WaitGroup

	// Create a goroutine to convert the first half of the Pineapple objects
	go func(chunk []Pineapple) {
		defer wg.Done() // Decrement the WaitGroup when the goroutine is done
		for idx, pineapple := range chunk { // Loop through the chunk of Pineapple objects
			mutex.Lock() // Lock the mutex
			safePineApples[idx] = pineapple.ToSafePineApple() // Convert the Pineapple object to a SafePineApple object
			mutex.Unlock()  // Unlock the mutex
		}
	}(chunks[0]) // Pass the first half of the Pineapple objects to the goroutine

	// Convert the second half of the Pineapple objects in the main thread
	for idx, pineapple := range chunks[1] {
		mutex.Lock() // Lock the mutex
		safePineApples[idx+len(chunks[0])] = pineapple.ToSafePineApple() // Convert the Pineapple object to a SafePineApple object
		mutex.Unlock() // Unlock the mutex
	}

	// Wait for all goroutines to finish
	wg.Wait()

	return safePineApples
}
```

This function make use of the mutex to lock the slice when writing to it. 
This makes sure that the goroutine and the main thread don't write to the same index at the same time but instead wait for the other to finish.


## Benchmark

Now that we have our functions we can benchmark them to see which one is the fastest. 
To benchmark our functions we run them with arrays of different sizes and see which one is the fastest.

Our benchmark function looks like this:

```go
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
```

To run the benchmark we use the following command:
```bash
go test -bench=. -benchtime 5s > benchmark.txt && benchstat benchmark.txt   
```

Benchstat is a tool that can be used to compare the results of benchmarks. 

The results of the benchmark are as follows:

```bash
name                                                                        time/op
Benchmark_SimpleConvertPineApplesToSafety-500-32                          15.8µs ± 0%
Benchmark_SimpleConvertPineApplesToSafety-1000-32                         32.0µs ± 0%
Benchmark_SimpleConvertPineApplesToSafety-2000-32                         66.5µs ± 0%
Benchmark_SimpleConvertPineApplesToSafety-5000-32                          193µs ± 0%
Benchmark_SimpleConvertPineApplesToSafety-10000-32                         465µs ± 0%
Benchmark_GoroutinesConvertPineApplesToSafety-500-32                      23.5µs ± 0%
Benchmark_GoroutinesConvertPineApplesToSafety-1000-32                     46.2µs ± 0%
Benchmark_GoroutinesConvertPineApplesToSafety-2000-32                     87.7µs ± 0%
Benchmark_GoroutinesConvertPineApplesToSafety-5000-32                      242µs ± 0%
Benchmark_GoroutinesConvertPineApplesToSafety-10000-32                     507µs ± 0%
Benchmark_NoMutexGoroutinesConvertPineApplesToSafety-500-32               28.3µs ± 0%
Benchmark_NoMutexGoroutinesConvertPineApplesToSafety-1000-32              48.7µs ± 0%
Benchmark_NoMutexGoroutinesConvertPineApplesToSafety-2000-32               105µs ± 0%
Benchmark_NoMutexGoroutinesConvertPineApplesToSafety-5000-32               257µs ± 0%
Benchmark_NoMutexGoroutinesConvertPineApplesToSafety-10000-32              533µs ± 0%
```

As you can see the simple function is the fastest. 

## Why ? 

The reason why the simple function is the fastest is that the process of converting the Pineapple objects to SafePineApple objects is very fast.
The time it takes to create the goroutines and wait for them to finish is longer than the time it takes to convert the Pineapple objects to SafePineApple objects.
Furthermore, in our goroutines implementation we have to convert the Pineapple objects then lock the mutex, write to the slice and unlock the mutex. 
This is a lot of overhead for a very simple task. 

## Conclusion

Don't use goroutines when you don't need them. That may seem obvious, but it's easy to forget when you're trying to optimize your code. 
In this example the overhead of creating the goroutines and waiting for them to finish is longer than the time it takes to convert the Pineapple objects to SafePineApple objects.

Goroutines are great for tasks that take a long time to complete and can be done in parallel. 
I would suggest to write the simplest code possible and then benchmark it to see if you can improve it using goroutines instead. 

Moreover, simple code is easier to read and maintain than complex code, that's why writing complex code might not be necessary if performance is not an issue.
I'd prefer to have a simple function that takes a few milliseconds longer to complete than a complex function.

## Code

The code for this article can be found on [GitHub](https://github.com/CorentinGS/go-teaching/tree/main/goroutines_simple_vs_complex).

## License

This article and the corresponding code is licensed under the [ISC License](https://github.com/CorentinGS/go-teaching/blob/main/LICENSE). 
If you want to distribute or cite this article and the corresponding code you can do it as long as you follow the ISC License and this steps:
  - Cite the [GitHub repo](https://github.com/CorentinGS/go-teaching/tree/main/goroutines_simple_vs_complex)
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






