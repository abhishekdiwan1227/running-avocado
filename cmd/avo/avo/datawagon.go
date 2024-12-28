package avo

import (
	"fmt"
	"time"

	"github.com/hashicorp/cronexpr"
	"gorm.io/gorm"
)

type DataWagon struct {
	db *gorm.DB
}

type TaskDefinition interface {
	ScriptTaskDefinition
}

func StartWagon() *DataWagon {
	return &DataWagon{db: GetConfig().DB}
}

func (dw *DataWagon) GetAllWorkPassengers() *[]Task {
	var works []Task
	dw.db.Joins("Schedule").Joins("Schedule.ScheduleDetail").Find(&works)
	return &works
}

func (dw *DataWagon) GetNextPassengers(till time.Time, workPassengers *chan struct {
	Task
	time.Time
}) {
	nowUtc := time.Now().UTC()
	tillUtc := till.UTC()
	works := dw.GetAllWorkPassengers()
	defer close(*workPassengers)
	for _, work := range *works {
		var next time.Time
		var interval time.Duration
		switch work.Schedule.ScheduleType {
		case CronSchedule:
			cron, err := cronexpr.Parse(*work.Schedule.ScheduleDetail.CronString)
			if err != nil {
				panic(err.Error())
			}

			next2 := cron.NextN(nowUtc, 2)
			next = next2[0]
			interval = next2[1].Sub(next)
		case Interval:
			interval = time.Duration(*work.Schedule.ScheduleDetail.Interval) * *work.Schedule.ScheduleDetail.IntervalKind
			if n := nowUtc.Add(interval); n.After(tillUtc) {
				appTickerValue := GetConfig().Ticker.TickerValue
				appTickerDuration := GetConfig().Ticker.TickerDuration
				appTickerInterval := time.Duration(appTickerValue) * appTickerDuration
				next = nowUtc.Add(interval % appTickerInterval)
			} else {
				next = n
			}

		default:
			continue
		}

		for i := next; i.Before(tillUtc); i = i.Add(interval) {
			fmt.Printf("Scheduled %s at %s\n", work.Name, i)
			*workPassengers <- struct {
				Task
				time.Time
			}{work, i}
		}
	}
}

func (dw *DataWagon) GetScriptDefinition(taskID uint) *ScriptTaskDefinition {
	definition := ScriptTaskDefinition{TaskID: taskID}
	dw.db.First(&definition)
	return &definition
}

func (dw *DataWagon) AddWork(work *Task) {
	dw.db.Create(work)
}
