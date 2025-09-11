package main

import (
	"fmt"
	"time"
)

func Fork(id int, leftReq, rightReq chan string) {
	for {
		select {
		case <-leftReq:
			leftReq <- "ok"
			<-leftReq
			fmt.Println("Fork", id, "released by left philosopher")

		case <-rightReq:
			rightReq <- "ok"
			<-rightReq
			fmt.Println("Fork", id, "released by right philosopher")
		}
	}
}

func Philosopher(name string, leftReq, rightReq chan string) {
	for {
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

				}
			case <-time.After(1 * time.Second):
				fmt.Println("Could not find right fork. Realeasing left fork")
				leftReq <- "release"
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func main() {

	ch1 := make(chan string, 1)
	ch2 := make(chan string, 1)
	ch3 := make(chan string, 1)
	ch4 := make(chan string, 1)
	ch5 := make(chan string, 1)

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
