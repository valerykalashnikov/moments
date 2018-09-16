package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/valerykalashnikov/moments"
)

func main() {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	rateBackupPath := os.Getenv("RATE_BACKUP_PATH")

	if rateBackupPath == "" {
		rateBackupPath = "moments.bak"
		var _, err = os.Stat(rateBackupPath)

		// create rate backup file if not exists
		if os.IsNotExist(err) {
			var file, err = os.Create(rateBackupPath)
			if err != nil {
				log.Fatalf("Unable to create file for moments backups, %s", err)
			}
			file.Close()
		}
	}

	f, err := os.OpenFile(rateBackupPath, os.O_RDWR, os.ModeAppend)
	if err != nil {
		log.Fatalf("Unable to open file to backup moments, %s", err)
	}
	defer f.Close()

	var counter *moments.MomentsCounter

	counter, err = moments.NewMomentsCounterFrom(f)
	if err != nil {
		if err == moments.ErrNoSavedMoments {
			counter = moments.NewMomentsCounter(1 * time.Minute)
		} else {
			log.Fatalf("Unable to backup file, %s", err)
		}

	}

	handleSigterm(func() {
		if _, err := f.Seek(0, 0); err != nil {
			log.Fatalf("Unable to reset rate file before backing up, %s", err)
		}
		counter.Save(f)
	})

	fmt.Println("Listening on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), NewAppHandler(counter)))

}

func handleSigterm(handleExit func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		handleExit()
		os.Exit(1)
	}()
}
