package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	gettingReadyForMission()
	gettingReadyForMissionWithCond()
}

var ready bool

func gettingReadyForMissionWithCond() {
	cond := sync.NewCond(&sync.Mutex{})
	go gettingReadyWithCond(cond)
	workIntervals := 0
	cond.L.Lock()
	for !ready {
		workIntervals++
		cond.Wait()
	}
	cond.L.Unlock()
	fmt.Printf("We are now ready! After %d work intervals.\n", workIntervals)
}

func gettingReadyWithCond(cond *sync.Cond) {
	sleep()
	ready = true
	cond.Signal()
}

func gettingReadyForMission() {
	go gettingReady()
	workIntervals := 0
	for !ready {
		workIntervals++
	}
	fmt.Printf("We are now ready! After %d work intervals.\n", workIntervals)
}

func gettingReady() {
	sleep()
	ready = true
}

func sleep() {
	someTime := time.Duration(1+rand.Intn(5)) * time.Second
	time.Sleep(someTime)
}
