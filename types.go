package avo

import (
	"time"

	"gorm.io/gorm"
)

type Work struct {
	gorm.Model
	Name       string
	ScheduleID int
	Schedule   *Schedule
}

type ScheduleType int
type IntervalKind int

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
	Work Work
	fn   func(*Work)
}

func (w *Workload) Do() {
	w.fn(&w.Work)
}
