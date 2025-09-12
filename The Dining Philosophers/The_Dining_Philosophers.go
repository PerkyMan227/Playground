package main

import (
	"fmt"
	"math/rand/v2"
	"time"
)

// waiter channel acts as a semaphore to limit the number of philosophers attempting to eat
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

// Uses an asymmetric fork acquisition strategy for one philosopher to prevent deadlock.
func Philosopher(name string, id int, leftForkComm, rightForkComm chan string) {
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

	forkSides := make([][2]chan string, 5)

	for i := 0; i < 5; i++ {
		forkSides[i][0] = make(chan string, 2)
		forkSides[i][1] = make(chan string, 2)
	}

	for i := 0; i < 5; i++ {
		go Fork(i+1, forkSides[i][0], forkSides[i][1])
	}

	go Philosopher("Kant", 1, forkSides[0][1], forkSides[1][0])        // P1: Left F1 (right side of F1), Right F2 (left side of F2)
	go Philosopher("Karl Marx", 2, forkSides[1][1], forkSides[2][0])   // P2: Left F2 (right side of F2), Right F3 (left side of F3)
	go Philosopher("Søren Pape", 3, forkSides[2][1], forkSides[3][0])  // P3: Left F3 (right side of F3), Right F4 (left side of F4)
	go Philosopher("Kong Fuzi", 4, forkSides[3][1], forkSides[4][0])   // P4: Left F4 (right side of F4), Right F5 (left side of F5)
	go Philosopher("Aristoteles", 5, forkSides[4][1], forkSides[0][0]) // P5: Left F5 (right side of F5), Right F1 (left side of F1)

	// Keep the main goroutine alive
	select {}
}

/*
Explanation of Deadlock Prevention:

The system does not deadlock due to two primary mechanisms:

1.  **The Waiter Semaphore (Capacity of N-1):**
    The `waiter` channel acts as a semaphore, ensuring that at most `N-1` (where N=5 philosophers, so 4) philosophers can be in the "attempting to eat" state (i.e., trying to acquire forks) at any given time.
    *   **Why this prevents deadlock:** If all N philosophers could sit at the table and simultaneously pick up their first fork, they would each be holding one fork and waiting for another fork that is held by their neighbor. This creates a classic circular dependency, leading to a deadlock. By limiting the number of philosophers who can attempt to eat to N-1, we guarantee that at least one philosopher will always be able to pick up *both* forks.
    *   **How it works:** If 4 philosophers pick up their left forks, the 5th philosopher cannot proceed until one of the first 4 finishes eating and releases the `waiter` token. The 4 philosophers will attempt to pick up their right forks. Since there are only 4 active, there will always be at least one free right fork, allowing at least one philosopher to complete their meal, release both forks, and the `waiter` token. This breaks the cycle.

2.  **Asymmetric Fork Acquisition (Breaking Symmetry):**
    One philosopher (in this case, "Aristoteles" with ID 5) is programmed to pick up their forks in a different order than the others.
    *   **Normal philosophers (IDs 1-4):** Pick up the **LEFT** fork first, then the **RIGHT** fork.
    *   **Asymmetric philosopher (ID 5):** Picks up the **RIGHT** fork first, then the **LEFT** fork.
    *   **Why this prevents deadlock:** This strategy ensures that a circular wait condition, where every philosopher is holding their first fork and waiting for their second fork (which is held by the next philosopher in the circle), cannot form. By having one philosopher acquire the right fork first, they essentially create a "break" in the potential chain of dependency. If all others are waiting for their right fork (held by the next person), the asymmetric philosopher might be able to acquire their right fork because it's the "first" they're trying to get, and then their left. This strategy, combined with the waiter, provides robust deadlock prevention.
*/
