package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func main() {
	count := 0
	total := 0
	
	var mu sync.Mutex
	cond := sync.NewCond(&mu)

	for i := 0; i < 10; i++ {
		go func(person int) {
			vote := requestVote(person)
			mu.Lock()
			defer mu.Unlock()
			if vote {
				count += 1
			}
			total += 1

			/*
				Broadcast wakes up any sleeping thread 
				General Syntax for worker is acquire the lock, perform the operation
				broadcase and then release the lock
			*/
			cond.Broadcast()
		}(i)
	}
	mu.Lock()
	for count < 5 && total != 10 {
		/* 
			Waits for some thread to Broadcase, upon which this main thread
			re-acquires the locks and then continues with the execution of the loop
			General syntax for main thread is acquire the lock, while condition is false, wait, 
			then do something and then release the lock
			This wait will prevent wastage of CPU cycle. If did not have this, CPU
			will waste time looping over the for loop 
		*/
		cond.Wait()
	}
	if count >= 5 {
		fmt.Println("Got 5+ votes")
	} else {
		fmt.Println("Lost")
	}
	mu.Unlock()
}

func requestVote(person int) bool {
    vote := rand.Intn(2) == 1
	fmt.Printf("Person %d voted %t\n", person, vote)
	return vote
}