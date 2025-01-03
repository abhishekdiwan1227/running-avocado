package main

import (
	"log"
	"sync"
	"time"

	"github.com/abhishekdiwan1227/running-avocado/cmd/avo/avo"
	"github.com/abhishekdiwan1227/running-avocado/lib"
)

var (
	wagonTicker *time.Ticker = time.NewTicker(avocado.GetConfig().Ticker.TickerDuration * time.Duration(avocado.GetConfig().Ticker.TickerValue))
	wgi         sync.WaitGroup
)

func main() {
	avocado.Start()
	wgi.Add(1)
	startWorkReminder()
	wgi.Wait()
}

func startWorkReminder() {
	for {
		tick := (<-wagonTicker.C).UTC()
		log.Print("Checking for work")
		next := tick.Add(5 * time.Second)

		works := make(chan struct {
			avocado.Task
			time.Time
		})
		wgi.Add(1)
		go avo.StartQueue(&works)
		wagon := avocado.GetConfig().Wagon
		wagon.GetNextPassengers(next, &works)
		wgi.Done()
	}
}
