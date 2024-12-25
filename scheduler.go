package avo

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup
var workloads chan Workload = make(chan Workload)

func StartQueue(works *chan struct {
	Task
	time.Time
}) {
	wg.Add(1)
	go letsThis()
	schedule(works)
	wg.Wait()
	close(workloads)
}

func schedule(works *chan struct {
	Task
	time.Time
}) {
	for work := range *works {
		now := time.Now().UTC()
		fmt.Println("Now", now)
		delay := work.Time.Sub(now)

		fmt.Printf("%s is scheduled to run in %f\n", work.Name, delay.Seconds())
		scheduleWork(delay, work.Task)
	}
}

func scheduleWork(delay time.Duration, work Task) {
	if delay > 0 {
		fn := func(currentWork *Task) {
			wg.Add(1)
			time.AfterFunc(delay, func() {
				fmt.Printf("[%s] Running %s\n", time.Now().UTC(), currentWork.Name)
				wg.Done()
			})
		}
		workloads <- Workload{work, fn}
	}
}

func letsThis() {
	defer wg.Done()
	for {
		workload, ok := <-workloads
		if !ok {
			return
		}
		wg.Add(1)
		go func() {
			workload.Do()
			wg.Done()
		}()
	}
}
