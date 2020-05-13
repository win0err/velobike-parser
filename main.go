package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
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

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetConnMaxLifetime(time.Hour)

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalln("[FATAL] Unable to migrate DB:", err)
	}
}

func main() {
	go interruptHandler()

	mode := ""
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}

	switch mode {
	case "import":
		fromFile := len(os.Args) > 2

		var reader io.Reader
		if fromFile {
			file, _ := os.Open(os.Args[2])
			reader = bufio.NewReader(file)
		} else {
			reader = os.Stdin
		}

		if err := importData(reader); err != nil {
			panic(err)
		}

	case "export":
		if len(os.Args) < 3 {
			fmt.Printf(
				"Usage: \n"+
					"%s export -from=\"2006-01-02 15:04 MST\" -to=\"2006-01-02 15:04 MST\"\n"+
					"%s export -all\n",
				os.Args[0],
				os.Args[0],
			)
			os.Exit(1)
		}
		exportCmd := flag.NewFlagSet("export", flag.ExitOnError)
		all := exportCmd.Bool("all", false, "all")
		from := exportCmd.String("from", "", "from")
		to := exportCmd.String("to", "", "to")

		exportCmd.Parse(os.Args[2:])

		data, err := exportData(*all, *from, *to)
		if err == nil {
			os.Stdout.Write(data)
		} else {
			panic(err)
		}

	default:
		errorHandler := func(err error) {
			log.Println("[ERROR] Unable to get parkings data:", err)

			log.Println("[INFO] Retry in 5 seconds...")
			time.Sleep(5 * time.Second)
		}

		for !isInterrupted {
			wg.Add(1)

			response, err := parkings.Get()
			if err != nil {
				errorHandler(err)
				wg.Done()
				continue
			}

			states := parkings.ToStates(*response)
			if len(states) == 0 {
				errorHandler(err)
				wg.Done()
				continue
			}

			go processResponse(states)
			helpers.SleepUntilNextMinute()
		}
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
