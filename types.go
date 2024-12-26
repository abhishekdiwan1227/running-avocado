package avo

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Name               string
	ScheduleID         int
	Schedule           *Schedule
	Active             bool                `gorm:"default:true"`
	TaskDefinitionID   *uint               `gorm:"index:idx_task_definition"`
	TaskDefinitionType *TaskDefinitionType `gorm:"index:idx_task_definition"`
}

type ScheduleType int

const (
	Interval ScheduleType = iota
	CronSchedule
)

func (ScheduleType *ScheduleType) String() string {
	switch *ScheduleType {
	case Interval:
		return "Interval"
	case CronSchedule:
		return "Cron"
	default:
		return "Invalid"
	}
}

type Schedule struct {
	gorm.Model
	ScheduleDetailID int
	ScheduleDetail   *ScheduleDetail
	ScheduleType     ScheduleType
}

type ScheduleDetail struct {
	gorm.Model
	CronString   *string
	Interval     *int
	IntervalKind *time.Duration
}

type Workload struct {
	Task Task
	fn   func(*Task)
}

func (w *Workload) Do() { w.fn(&w.Task) }

type TaskDefinitionType int

const (
	Local TaskDefinitionType = iota
)

type ScriptTaskDefinition struct {
	gorm.Model
	TaskId    uint
	Path      string
	Arguments *string
}
