package main

import "time"

func sleepUntilNextMinute() {
	timeToSleep := time.Duration(60 - time.Now().Second())
	time.Sleep(timeToSleep * time.Second)
}
