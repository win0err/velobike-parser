package helpers

import "time"

func SleepUntilNextMinute() {
	timeToSleep := time.Duration(60 - time.Now().Second())
	time.Sleep(timeToSleep * time.Second)
}

// GetCurrentTime returns time accurate to minutes
// Example: 2020-04-15 16:39:00.000000000
func GetCurrentTime() time.Time {
	return time.Now().Truncate(time.Minute)
}