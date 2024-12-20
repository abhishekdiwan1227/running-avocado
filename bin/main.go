package main

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/abhishekdiwan1227/avo"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	homePathString := os.Getenv("HOME")
	projectDir := filepath.Join(homePathString, ".avo")
	dbPath := filepath.Join(projectDir, "avo.db")
	if _, err := os.Stat(dbPath); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(projectDir, os.ModePerm)
		if err != nil {
			panic(err.Error())
		}
	}
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	err = db.AutoMigrate(&avo.Work{})
	if err != nil {
		panic(err.Error())
	}

	db.Create(&avo.Work{
		Name: "Dummy Work",
		Schedule: &avo.Schedule{
			ScheduleDetail: &avo.ScheduleDetail{
				Interval:     func(i int) *int { return &i }(1000),
				IntervalKind: func(i time.Duration) *time.Duration { return &i }(time.Millisecond),
			},
			ScheduleType: avo.Interval,
		},
	})
	// create dummy data for work with schedule type cron

	avo.Init(db)
}
