package main2

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

const (
	NB_BABOONS = 10
)

var (
	mutexTurn      = new(sync.Mutex)
	countSemaphore = semaphore.NewWeighted(int64(50))
	mutexWest      = new(sync.Mutex)
	mutexEast      = new(sync.Mutex)
	mutexRope      = new(sync.Mutex)
	counterEast    = 0
	counterWest    = 0
	ctx            = context.Background()
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	wg := new(sync.WaitGroup)

	for i := 1; i <= NB_BABOONS; i++ {
		time.Sleep(1 * time.Second)
		weight := 0
		direction := ""
		if rand.Intn(5) == 0 {
			direction = "east"
		} else {
			direction = "west"
		}

		if rand.Intn(2) == 0 {
			weight = 10
		} else {
			weight = 20
		}

		wg.Add(1)
		go baboon(direction, weight, i, wg)
	}
	wg.Wait()
}

func baboon(direction string, weight int, nb int, wg *sync.WaitGroup) {
	fmt.Printf("Baboon %d created direction %s and weight %d\n", nb, direction, weight)
	var mutexDirection *sync.Mutex
	var counter *int

	if direction == "east" {
		mutexDirection = mutexEast
		counter = &counterEast
	} else if direction == "west" {
		mutexDirection = mutexWest
		counter = &counterWest
	} else {
		fmt.Printf("Error direction")
	}

	mutexTurn.Lock()
	mutexDirection.Lock()
	*counter++
	if *counter == 1 {
		mutexRope.Lock()
	}
	mutexDirection.Unlock()
	mutexTurn.Unlock()

	countSemaphore.Acquire(ctx, int64(weight))
	fmt.Printf("->Baboon %d crossing direction %s (+ %d)\n", nb, direction, weight)
	time.Sleep(6 * time.Second)
	fmt.Printf("<-Baboon %d finished crossing and arrived to side %s (- %d)\n", nb, direction, weight)
	countSemaphore.Release(int64(weight))

	mutexDirection.Lock()
	*counter--
	if *counter == 0 {
		mutexRope.Unlock()
	}
	mutexDirection.Unlock()

	wg.Done()
}
