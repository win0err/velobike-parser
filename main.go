package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/win0err/velobike-parser/database"
	"github.com/win0err/velobike-parser/helpers"
	"github.com/win0err/velobike-parser/parkings"
)

var wg = &sync.WaitGroup{}
var isInterrupted = false

func init() {
	db, err := database.GetConnection()
	if err != nil {
		log.Fatalln("[FATAL] Unable get initial DB connection:", err)
	}
	defer db.Close()

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalln("[FATAL] Unable to migrate DB:", err)
	}
}

func main() {
	go interruptHandler()

	for !isInterrupted {
		wg.Add(1)
		response, err := parkings.Get()
		if err != nil {
			log.Println("[ERROR] Unable to get parkings data:", err)

			log.Println("[INFO] Retry in 5 seconds...")
			time.Sleep(5 * time.Second)

			wg.Done()
			continue
		}

		go processResponse(response)
		helpers.SleepUntilNextMinute()
	}
}

func interruptHandler() {
	go func() {
		chSignal := make(chan os.Signal, 1)
		signal.Notify(chSignal, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

		for {
			select {
			case sig := <-chSignal:
				log.Println("[INFO] Received", sig)
				isInterrupted = true

				log.Println("[INFO] Shutting down...")
				wg.Wait()

				os.Exit(0)
			default:
			}
		}
	}()
}
