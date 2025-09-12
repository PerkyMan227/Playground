package main

import (
	"fmt"
	"math/rand/v2"
	"time"
)

var waiter = make(chan struct{}, 4)

func Fork(id int, leftReq, rightReq chan string) {
	fmt.Println("Fork", id, "is ready.")
	select {}
}

func Philosopher(name string, id int, leftReq, rightReq chan string) {
	meals := 0
	for meals < 3 {

		fmt.Println(name, "is thinking...")

		time.Sleep(time.Second * time.Duration(rand.IntN(3)))
		waiter <- struct{}{}

		leftReq <- "want"
		if <-leftReq == "ok" {
			fmt.Println(name, "picked up LEFT fork")

			rightReq <- "want"
			select {
			case reply := <-rightReq:
				if reply == "ok" {
					fmt.Println(name, "picked up right fork.")
					fmt.Println(name, " IS EATING")
					time.Sleep(2 * time.Second)

					leftReq <- "release"
					rightReq <- "release"
					fmt.Println(name, "released both forkse. Thinking...")
					meals++

				}
			case <-time.After(1 * time.Second):
				fmt.Println("Could not find right fork. Realeasing left fork")
				leftReq <- "release"
			}
		}
		<-waiter
		time.Sleep(4 * time.Second)
	}
	fmt.Println(name, "Has eaten 3 big meals, and leaves the table")
}

func main() {

	ch1 := make(chan string, 2)
	ch2 := make(chan string, 2)
	ch3 := make(chan string, 2)
	ch4 := make(chan string, 2)
	ch5 := make(chan string, 2)

	go Philosopher("Kant", ch1, ch2)
	go Philosopher("Karl Marx", ch2, ch3)
	go Philosopher("Søren Pape", ch3, ch4)
	go Philosopher("Kong Fuzi", ch4, ch5)
	go Philosopher("Aristoteles", ch5, ch1)

	go Fork(1, ch1, ch2)
	go Fork(2, ch2, ch3)
	go Fork(3, ch3, ch4)
	go Fork(4, ch4, ch5)
	go Fork(5, ch5, ch1)

	select {}

}
