package main

import (
	"fmt"
	"time"
)

func main() {
	smokeSignal := make(chan bool)
	evilNinja := "Tommy"
	go attack(evilNinja, smokeSignal)
	// Receive message - variable on the right side of the arrow
	fmt.Println(<-smokeSignal)
}

func attack(target string, attacked chan bool) {
	fmt.Println("Throwing ninja stars at", target)
	time.Sleep(time.Second)
	// Send message to channel - variable on the left side of the arrow
	attacked <- true
}
