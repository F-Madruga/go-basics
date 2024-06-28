package main

import (
	"fmt"
	"math/rand"
)

func main() {
	channel := make(chan string)

	go throwingNinjaStar(channel)
	for message := range channel {
		fmt.Println(message)
	}
}

func throwingNinjaStar(channel chan string) {
	numRounds := 3
	for i := 0; i < numRounds; i++ {
		score := rand.Intn(10)
		channel <- fmt.Sprintf("You scored %d", score)
	}
	close(channel)
}
