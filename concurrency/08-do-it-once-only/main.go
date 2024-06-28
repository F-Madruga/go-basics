package main

import (
	"fmt"
	"math/rand"
	"sync"
)

var missionCompleted bool

func main() {
	numberOfNinjas := 100

	var wg sync.WaitGroup
	wg.Add(numberOfNinjas)

	var once sync.Once

	for i := 0; i < numberOfNinjas; i++ {
		go func() {
			if foundTreasure() {
				once.Do(markMissionCompleted)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	checkMissionCompletion()
}

func checkMissionCompletion() {
	if missionCompleted {
		fmt.Println("Mission is now completed.")
	} else {
		fmt.Println("Mission was a failure.")
	}
}

func markMissionCompleted() {
	missionCompleted = true
	fmt.Println("Marking mission completed.")
}

func foundTreasure() bool {
	return rand.Intn(10) == 0
}
