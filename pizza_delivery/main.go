package main

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"
)

const (
	NumberOfDrivers = 4
	NumberOfCooks   = 3
)

type Pizza struct {
	Name         string
	deliveryTime time.Duration
}

var (
	deliveryChan = make(chan Pizza, 100)
	orderChan    = make(chan Pizza, 100)
	closeChan    = make(chan struct{})
	pizzasNames  = []string{"Pepperoni", "Cheese", "Veggie", "Hawaiian"}
)

func main() {
	log.Printf("Hello, world!")

	log.Printf("Welcome to the pizza delivery service!")

	ctx, cancelFunc := context.WithCancel(context.Background())

	// Create a wait group to wait for all the goroutines to finish
	var wg sync.WaitGroup

	wg.Add(NumberOfCooks)
	for i := 0; i < NumberOfCooks; i++ {
		go barista(&wg, ctx)
	}

	wg.Add(NumberOfDrivers)
	for i := 0; i < NumberOfDrivers; i++ {
		go driver(i, &wg, ctx)
	}

	wg.Add(1)

	go func(ctx context.Context) {
		for {
			select {
			case <-time.After(1 * time.Second):
				client()
			case <-ctx.Done():
				wg.Done()
				return
			}
		}
	}(ctx)

	time.Sleep(20 * time.Second)
	cancelFunc()
	ctx.Done()

	log.Printf("Restaurant is closing")

	wg.Wait()

	close(closeChan)
	close(orderChan)
	close(deliveryChan)

	log.Printf("All done!")

}

func driver(id int, wg *sync.WaitGroup, ctx context.Context) {
	counter := 0
	duration := 0 * time.Second
	log.Printf("Driver #%d is ready to deliver!", id)
	for {
		select {
		case pizza := <-deliveryChan:
			log.Printf("Driver #%d is delivering a %s pizza", id, pizza.Name)
			time.Sleep(pizza.deliveryTime)
			counter++
			duration += pizza.deliveryTime
			if counter == 3 {
				log.Printf("Driver #%d delivered %d pizzas in %s", id, counter, duration)
				counter = 0
				duration = 0 * time.Second
			}
		case <-ctx.Done():
			wg.Done()
			return
		}
	}
}

func barista(wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	for {
		select {
		case pizza := <-orderChan:
			log.Printf("Barista is cooking a %s pizza", pizza.Name)
			time.Sleep(3 * time.Second)
			deliveryChan <- pizza
		case <-ctx.Done():
			wg.Done()
			return

		}
	}
}

func client() {
	orderChan <- Pizza{
		Name:         getRandomPizzaName(),
		deliveryTime: getRamdomDeliveryTime(),
	}
}

func getRandomPizzaName() string {
	return pizzasNames[rand.Intn(len(pizzasNames))]
}

func getRamdomDeliveryTime() time.Duration {
	return time.Duration(rand.Intn(5)) * time.Second
}
