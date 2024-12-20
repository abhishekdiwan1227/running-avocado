package avo

import (
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/cronexpr"
	"gorm.io/gorm"
)

var wg sync.WaitGroup
var workloads chan Workload = make(chan Workload)

func Init(db *gorm.DB) {
	var works []Work
	db.Joins("Schedule").Joins("Schedule.ScheduleDetail").Find(&works)

	letsThis(workloads)
	wg.Add(1)
	for _, work := range works {
		switch work.Schedule.ScheduleType {
		case CronSchedule:
			cron, err := cronexpr.Parse(*work.Schedule.ScheduleDetail.CronString)
			if err != nil {
				panic(err.Error())
			}
			next := cron.Next(time.Now())
			delay := time.Until(next)
			fmt.Printf("%s(Type: %s) scheduled at %s\n", work.Name, work.Schedule.ScheduleType.String(), next.String())

			if delay > 0 {
				scheduleWork(delay, workloads, work)
			}
		case Interval:
			next := time.Now().Add(time.Duration(*work.Schedule.ScheduleDetail.Interval) * *work.Schedule.ScheduleDetail.IntervalKind)
			delay := time.Until(next)
			fmt.Printf("%s(Type: %s) scheduled at %s\n", work.Name, work.Schedule.ScheduleType.String(), next.String())

			if delay > 0 {
				scheduleWork(delay, workloads, work)
			}

		default: // do nothing
		}
	}

	close(workloads)
	wg.Wait()
}

func scheduleWork(delay time.Duration, workloads chan Workload, work Work) {
	if delay > 0 {
		fn := func(currentWork *Work) {
			wg.Add(1)
			time.AfterFunc(delay, func() {
				fmt.Println("Running", currentWork.Name)
				wg.Done()
			})
		}
		workloads <- Workload{work, fn}
	}
}

func letsThis(workloads chan Workload) {
	go func() {
		defer wg.Done()
		for {
			workload, ok := <-workloads
			if !ok {
				return
			}
			workload.Do()
		}
	}()
}
