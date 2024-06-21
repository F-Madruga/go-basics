package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	defer func() {
		fmt.Println(time.Since(start))
	}()

	var smokeSignals []chan bool
	evilNinjas := []string{"Tommy", "Johnny", "Bobby", "Andy"}

	for i, evilNinja := range evilNinjas {
		smokeSignals = append(smokeSignals, make(chan bool))
		go attack(evilNinja, smokeSignals[i])
	}

	for _, smokeSignal := range smokeSignals {
		<- smokeSignal
	}
}

func attack(target string, attacked chan bool)  {
	fmt.Println("Throwing ninja stars at", target)
	time.Sleep(time.Second)
	attacked <- true
}
