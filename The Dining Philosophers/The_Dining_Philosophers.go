package main

import (
	"fmt"
	"time"
)

func Philosopher(name string, leftFork Fork, rightFork Fork) {
	<-leftFork
	<-rightFork

	fmt.Println(name, "Time to eat!")
	time.Sleep(time.Second)

	fmt.Println("Done eating")
	time.Sleep(time.Second)

	rightFork.Unlock()
	leftFork.Unlock()

}

func main() {
	/*
		bobTurn := make(chan bool)
		mannyTurn := make(chan bool)
	*/

	fork1 := &Fork{rightPhilosopher: "Aristoteles", leftPhilosopher: "Platon"}
	fork2 := &Fork{rightPhilosopher: "Platon", leftPhilosopher: "Kant"}
	fork3 := &Fork{rightPhilosopher: "Kant", leftPhilosopher: "Karl Marx"}
	fork4 := &Fork{rightPhilosopher: "Karl Marx", leftPhilosopher: "Kong Fuzi"}
	fork5 := &Fork{rightPhilosopher: "Kong Fuzi", leftPhilosopher: "Aristoteles"}

	go Philosopher("Aristoteles", fork1, fork5)
	go Philosopher("Platon", fork2, fork1)
	go Philosopher("Kant", fork3, fork2)
	go Philosopher("Karl Marx", fork4, fork3)
	go Philosopher("Kong Fuzi", fork5, fork4)

	select {}

}

/*
package main

import (
	"fmt"
	"time"
)

// Worker represents Bob or Manny
func worker(name string, myTurn <-chan bool, otherTurn chan<- bool) {
	for {
		// Wait until it's my turn
		<-myTurn

		// Building a board
		fmt.Println(name, "is building a board...")
		time.Sleep(time.Second)

		// Resting
		fmt.Println(name, "is resting...")
		time.Sleep(time.Second)

		// Signal to the other worker that it's their turn
		otherTurn <- true
	}
}

func main() {
	// Channels for turn-taking
	bobTurn := make(chan bool)
	mannyTurn := make(chan bool)

	// Start workers
	go worker("Bob", bobTurn, mannyTurn)
	go worker("Manny", mannyTurn, bobTurn)

	// Kick things off: Bob starts first
	bobTurn <- true

	// Keep program alive
	select {}
}

package main

import (
	"fmt"
	"sync"
	"time"
)

const N = 10000000

var balance = 0

var arbiter sync.Mutex

func worker() {
	for i := 0; i < N; i++ {
		arbiter.Lock()
		balance++
		arbiter.Unlock()
	}
	fmt.Println("Done")
}

func main() {

	go worker()
	go worker()

	for {
		time.Sleep(1000 * time.Millisecond)
		fmt.Println(balance)
	}
}
*/
