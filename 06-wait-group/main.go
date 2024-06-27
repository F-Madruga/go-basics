package main

import (
	"fmt"
	"sync"
)

func main() {
	var beeper sync.WaitGroup
	evilNinjas := []string{"Tommy", "Johnny", "Bobby"}
	beeper.Add(len(evilNinjas))
	for _, evilNinja := range evilNinjas {
		go attack(evilNinja, &beeper)
	}
	beeper.Wait()
	fmt.Println("Mission completed")
}

func attack(evilNinja string, beeper *sync.WaitGroup) {
	fmt.Println("Attacked evil ninja:", evilNinja)
	beeper.Done()
}
