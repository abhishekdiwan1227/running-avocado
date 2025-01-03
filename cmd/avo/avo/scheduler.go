package avo

import (
	"log"
	"sync"
	"time"

	avocado "github.com/abhishekdiwan1227/running-avocado/lib"
)

var (
	wg        sync.WaitGroup
	workloads chan avocado.Workload = make(chan avocado.Workload)
)

func StartQueue(works *chan struct {
	avocado.Task
	time.Time
},
) {
	wg.Add(1)
	go letsThis()
	schedule(works)
	wg.Wait()
	close(workloads)
}

func schedule(works *chan struct {
	avocado.Task
	time.Time
},
) {
	for work := range *works {
		now := time.Now().UTC()
		delay := work.Time.Sub(now).Round(time.Second)

		scheduleWork(delay, work.Task)
	}
}

func scheduleWork(delay time.Duration, work avocado.Task) {
	if delay >= 0 {
		func(currentWork *avocado.Task) {
			wg.Add(1)
			time.AfterFunc(delay, func() {
				log.Printf("Running %s\n", currentWork.Name)
				runner := avocado.CreateRunner(currentWork)
				workloads <- avocado.Workload{Task: work, Runner: runner}
				wg.Done()
			})
		}(&work)
	}
}

func letsThis() {
	defer wg.Done()
	for {
		workload, ok := <-workloads
		if !ok {
			return
		}

		job := avocado.GetConfig().Wagon.AddNewJob(workload.Task.ID)

		wg.Add(1)
		go func() {
			result := workload.Do()
			avocado.GetConfig().Wagon.CompleteJobEntry(job, *&result.ReturnCode)
			wg.Done()
		}()
	}
}
