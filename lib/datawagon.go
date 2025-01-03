package avocado

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/cronexpr"
	"gorm.io/gorm"
)

type DataWagon struct {
	DB *gorm.DB
}

type TaskDefinition interface {
	ScriptTaskDefinition
}

func StartWagon(db *gorm.DB) *DataWagon {
	return &DataWagon{DB: db}
}

func (dw *DataWagon) GetAllWorkPassengers() *[]Task {
	var works []Task
	dw.DB.Joins("Schedule").Joins("Schedule.ScheduleDetail").Find(&works)
	return &works
}

func (dw *DataWagon) GetNextPassengers(till time.Time, workPassengers *chan struct {
	Task
	time.Time
},
) {
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
			log.Printf("Scheduled %s at %s\n", work.Name, i)
			*workPassengers <- struct {
				Task
				time.Time
			}{work, i}
		}
	}
}

func (dw *DataWagon) GetScriptDefinition(taskID uint) *ScriptTaskDefinition {
	definition := ScriptTaskDefinition{TaskID: taskID}
	dw.DB.First(&definition)
	return &definition
}

func (dw *DataWagon) AddScriptDefinition(definition *ScriptTaskDefinition, task *Task) {
	taskType := Local
	dw.DB.Create(definition)

	task.TaskDefinitionID = &definition.ID
	task.TaskDefinitionType = &taskType

	dw.DB.Save(task)
}

func (dw *DataWagon) AddTask(work *Task) {
	dw.DB.Create(work)
}

func (dw *DataWagon) AddNewJob(taskID uint) *Job {
	job := &Job{
		TaskID:    taskID,
		StartedAt: time.Now().UTC(),
		JobID:     uuid.New(),
	}
	dw.DB.Create(job)
	return job
}

func (dw *DataWagon) MigrateDatabase() {
	err := dw.DB.AutoMigrate(&Task{}, &ScriptTaskDefinition{}, &Job{})
	if err != nil {
		panic(err.Error())
	}
}

func (dw *DataWagon) CompleteJobEntry(job *Job, returnCode int) {
	endTime := time.Now().UTC()
	job.EndedAt = &endTime
	job.ReturnCode = &returnCode

	dw.DB.Save(job)
}
