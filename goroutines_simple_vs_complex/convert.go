package goroutines_simple_vs_complex

import "sync"

func SimpleConvertPineApplesToSafety(pineapples []Pineapple) []SafePineApple {
	safePineApples := make([]SafePineApple, len(pineapples))

	for idx, pineapple := range pineapples {
		safePineApples[idx] = pineapple.ToSafePineApple()
	}

	return safePineApples
}

func GoroutinesConvertPineApplesToSafety(pineapples []Pineapple) []SafePineApple {
	// Create a slice to store the SafePineApples
	safePineApples := make([]SafePineApple, len(pineapples))

	// Split the offers into chunks
	chunks := [][]Pineapple{pineapples[:len(pineapples)/2], pineapples[len(pineapples)/2:]}

	mutex := sync.Mutex{}

	var wg sync.WaitGroup
	wg.Add(1)

	go func(chunk []Pineapple) {
		defer wg.Done()
		for idx, pineapple := range chunk {
			mutex.Lock()

			safePineApples[idx] = pineapple.ToSafePineApple()

			mutex.Unlock()
		}
	}(chunks[0])

	for idx, pineapple := range chunks[1] {
		mutex.Lock()
		safePineApples[idx+len(chunks[0])] = pineapple.ToSafePineApple()
		mutex.Unlock()
	}

	// Wait for all goroutines to finish
	wg.Wait()

	return safePineApples
}

func GoroutinesNoMutexConvertPineApplesToSafety(pineapples []Pineapple) []SafePineApple {
	// Create a slice to store the SafePineApples
	safePineApples := make([]SafePineApple, len(pineapples)/2, len(pineapples))
	safePineApples2 := make([]SafePineApple, len(pineapples)/2)

	var wg sync.WaitGroup
	wg.Add(1)

	go func(chunk []Pineapple) {
		defer wg.Done()
		for idx, pineapple := range chunk {
			safePineApples[idx] = pineapple.ToSafePineApple()
		}
	}(pineapples[:len(pineapples)/2])

	for idx, pineapple := range pineapples[len(pineapples)/2:] {
		safePineApples2[idx] = pineapple.ToSafePineApple()
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Group both pineapples
	safePineApples = append(safePineApples, safePineApples2...)

	return safePineApples
}
