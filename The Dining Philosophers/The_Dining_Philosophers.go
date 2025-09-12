package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"sync"
	"time"
)

// waiter channel acts as a semaphore to limit the number of philosophers attempting to eat, to prevent deadlock
var waiter = make(chan struct{}, 4)

func Fork(id int, leftPhilosopherReq, rightPhilosopherReq chan string) {
	fmt.Printf("Fork %d is ready.\n", id)
	for {
		select {
		case msg := <-leftPhilosopherReq:
			if msg == "want" {
				leftPhilosopherReq <- "ok"
				fmt.Printf("Fork %d picked up by LEFT side philosopher.\n", id)
				<-leftPhilosopherReq
				fmt.Printf("Fork %d released by LEFT side philosopher.\n", id)
			}
		case msg := <-rightPhilosopherReq:
			if msg == "want" {
				rightPhilosopherReq <- "ok"
				fmt.Printf("Fork %d picked up by RIGHT side philosopher.\n", id)
				<-rightPhilosopherReq
				fmt.Printf("Fork %d released by RIGHT side philosopher.\n", id)
			}
		}
	}
}

// Uses a "left for some, right for others" fork strategy, to prevent deadlock.
func Philosopher(name string, id int, leftForkComm, rightForkComm chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	meals := 0
	for meals < 3 {
		fmt.Printf("%s (P%d) is thinking...\n", name, id)
		time.Sleep(time.Second * time.Duration(rand.IntN(3)))

		waiter <- struct{}{} //Philosopher requests permission to sit and eat. Works like a ticketing system.

		//Philosophers 1-4 pick up LEFT then RIGHT.
		//Philosopher 5 picks up RIGHT then LEFT.
		if id != 5 {
			fmt.Printf("%s (P%d) wants to pick up LEFT fork...\n", name, id)
			leftForkComm <- "want"
			_ = <-leftForkComm //Acnowledgment
			fmt.Printf("%s (P%d) picked up LEFT fork.\n", name, id)

			fmt.Printf("%s (P%d) wants to pick up RIGHT fork...\n", name, id)
			rightForkComm <- "want"
			_ = <-rightForkComm //Acnowledgment
			fmt.Printf("%s (P%d) picked up RIGHT fork.\n", name, id)
		} else {
			fmt.Printf("%s (P%d) wants to pick up RIGHT fork...\n", name, id)
			rightForkComm <- "want"
			_ = <-rightForkComm //Acnowledgment
			fmt.Printf("%s (P%d) picked up RIGHT fork.\n", name, id)

			fmt.Printf("%s (P%d) wants to pick up LEFT fork...\n", name, id)
			leftForkComm <- "want"
			_ = <-leftForkComm //Acnowledgment
			fmt.Printf("%s (P%d) picked up LEFT fork.\n", name, id)
		}

		fmt.Printf("%s (P%d) IS EATING!!!\n", name, id)
		//Philosopher is tired from his meal, he sleeps a bit...
		time.Sleep(time.Second * time.Duration(rand.IntN(3)+2))

		meals++

		fmt.Printf("%s (P%d) is releasing forks.\n", name, id)
		leftForkComm <- "release"
		rightForkComm <- "release"

		<-waiter

		fmt.Printf("%s (P%d) has eaten %f meal(s).\n", name, id, meals)
		time.Sleep(time.Second * 4)
	}
	fmt.Printf("%s (P%d) has eaten 3 big meals and leaves the table.\n", name, id)
}

func main() {

	var wg sync.WaitGroup

	forkSides := make([][2]chan string, 5)

	for i := 0; i < 5; i++ {
		forkSides[i][0] = make(chan string, 2)
		forkSides[i][1] = make(chan string, 2)
	}

	for i := 0; i < 5; i++ {
		go Fork(i+1, forkSides[i][0], forkSides[i][1])
	}

	wg.Add(5)
	go Philosopher("Kant", 1, forkSides[0][1], forkSides[1][0], &wg)
	go Philosopher("Karl Marx", 2, forkSides[1][1], forkSides[2][0], &wg)
	go Philosopher("Søren Pape", 3, forkSides[2][1], forkSides[3][0], &wg)
	go Philosopher("Kong Fuzi", 4, forkSides[3][1], forkSides[4][0], &wg)
	go Philosopher("Aristoteles", 5, forkSides[4][1], forkSides[0][0], &wg)

	wg.Wait()

	fmt.Println("All Philosophers has eaten, exiting program")

	os.Exit(0)
	select {}
}
