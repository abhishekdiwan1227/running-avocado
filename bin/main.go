package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/abhishekdiwan1227/avo"
)

var wagonTicker *time.Ticker = time.NewTicker(avo.GetConfig().Ticker.TickerDuration * time.Duration(avo.GetConfig().Ticker.TickerValue))
var wgi sync.WaitGroup

func main() {
	avo.Start()
	wgi.Add(1)
	startWorkReminder()
	wgi.Wait()
}

func startWorkReminder() {
	for {
		tick := (<-wagonTicker.C).UTC()
		fmt.Printf("[%s] Checking for work\n", tick)
		next := tick.Add(5 * time.Second)

		works := make(chan struct {
			avo.Task
			time.Time
		})
		wgi.Add(1)
		go avo.StartQueue(&works)
		wagon := avo.GetConfig().Wagon
		wagon.GetNextPassengers(next, &works)
		wgi.Done()
	}
}
