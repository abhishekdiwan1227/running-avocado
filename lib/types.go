package avocado

import (
	"io"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	ScheduleType       int
	TaskDefinitionType int
)

const (
	Interval ScheduleType = iota
	CronSchedule
)

const (
	Local TaskDefinitionType = iota
)

type Task struct {
	gorm.Model
	Name               string              `json:"name"`
	ScheduleID         int                 `json:"-"`
	Schedule           *Schedule           `json:"schedule"`
	Active             bool                `gorm:"default:true" json:"active"`
	TaskDefinitionID   *uint               `gorm:"index:idx_task_definition" json:"taskDefinitionId"`
	TaskDefinitionType *TaskDefinitionType `gorm:"index:idx_task_definition" json:"taskDefinitionType"`
}

type Schedule struct {
	gorm.Model
	ScheduleDetailID int             `json:"-"`
	ScheduleDetail   *ScheduleDetail `json:"scheduleDetails"`
	ScheduleType     ScheduleType    `json:"type"`
}

type ScheduleDetail struct {
	gorm.Model
	CronString   *string        `json:"cron"`
	Interval     *int           `json:"interval"`
	IntervalKind *time.Duration `json:"intervalKind"`
}

type Workload struct {
	Task   Task
	Runner Runner
}

type WorkloadResult struct {
	ReturnCode   int
	StdoutWriter *io.Writer
	StderrWriter *io.Writer
	CompletedAt  time.Time
}

type ScriptTaskDefinition struct {
	gorm.Model
	TaskID     uint    `json:"taskId"`
	Command    string  `json:"command"`
	Entrypoint string  `json:"entrypoint"`
	Arguments  *string `json:"args"`
}

type Job struct {
	gorm.Model
	TaskID     uint       `json:"taskId"`
	StartedAt  time.Time  `json:"startedAt"`
	EndedAt    *time.Time `json:"endedAt"`
	ReturnCode *int       `json:"returnCode"`
	JobID      uuid.UUID  `json:"jobId"`
}

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

func (w *Workload) Do() *WorkloadResult { return w.Runner.Run(&w.Task) }
