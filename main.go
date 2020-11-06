package main

import (
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/avast/retry-go"
	"github.com/robfig/cron/v3"

	"github.com/win0err/velobike-parser/database"
	"github.com/win0err/velobike-parser/helpers"
	"github.com/win0err/velobike-parser/parkings"
	"github.com/win0err/velobike-parser/savers"
)

var c = cron.New(cron.WithSeconds())
var wg = &sync.WaitGroup{}
var noMigrate = flag.Bool("no-migrate", false, "disable automigration of database schema on start")

func init() {
	flag.Parse()
}

func main() {
	if !*noMigrate {
		if err := database.AutoMigrate(); err != nil {
			helpers.Log.Fatal("automigration failed:", err)
		}

		helpers.Log.Info("automigration successfully completed")
	}

	schedule, err := cron.ParseStandard(helpers.Config.ParseInterval)
	if err != nil {
		helpers.Log.Fatal("wrong parse interval:", err)
	}

	nextRun := schedule.Next(time.Now())
	retryDuration := schedule.Next(nextRun).Sub(nextRun) / 4

	// todo: sleep until run. e.g.
	// if we want to parse every 15 minutes,
	// we need to invoke the function in 00, 15, 30, 45 minutes (not 05, 20, 35, 50)
	helpers.Log.Info("starting in", nextRun.Sub(time.Now()))

	job := func() {
		wg.Add(1)
		defer wg.Done()

		err := retry.Do(
			parse,
			retry.Attempts(3),
			retry.OnRetry(func(n uint, err error) { helpers.Log.Infof("retrying (%d attempt)...\n", n + 1) }),
			retry.Delay(retryDuration),
		)

		if err != nil {
			helpers.Log.Critical("error while retrying:", err)
		}
	}

	c.Schedule(schedule, cron.FuncJob(job))

	go onInterrupt(func() { c.Stop() })

	c.Run()
	wg.Wait()

	database.Connection.DB().Close()

	helpers.Log.Info("shutting down...")
}


func parse() error {
	req := parkings.NewRequest()

	if err := req.Get(); err != nil {
		helpers.Log.Warning("unable to get parkings data:", err)

		return err
	}

	if err := req.Parse(); err != nil {
		helpers.Log.Warning("unable to parse parkings data:", err)

		return err
	}

	states := parkings.ToStates(*req.ParsedResponse)

	if len(states) > 0 {
		currentTime := states[0].Time

		if err := savers.ToDb(states); err == nil {
			helpers.Log.Info("data successfully saved for", currentTime)
		} else {
			helpers.Log.Warning("unable to save to database:", err)

			if err := savers.ToJson(states); err == nil {
				helpers.Log.Infof("data backed up for %s\n", currentTime)
			} else {
				helpers.Log.Warning("error while saving to JSON:", err)
			}

			return err
		}
	}

	return nil
}

func onInterrupt(cmd func()) {
	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	for {
		select {
		case sig := <-chSignal:
			helpers.Log.Info("received", sig)
			cmd()
		default:
		}
	}
}