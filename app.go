package avo

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type App struct {
	DB     *gorm.DB
	Ticker AppTicker
	Config AppConfig
	Wagon  *DataWagon
}

type AppTicker struct {
	TickerValue    int
	TickerDuration time.Duration
}

type AppConfig struct {
	AppDirectoryPath string
}

var app *App = &App{
	Ticker: AppTicker{
		TickerValue:    5,
		TickerDuration: time.Second,
	},
}

var once sync.Once

func Start() {
	once.Do(initializeApp)
}

func initializeApp() {
	homePathString := os.Getenv("HOME")
	projectDir := filepath.Join(homePathString, ".avo")
	app.Config.AppDirectoryPath = projectDir
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
	app.DB = db

	err = db.AutoMigrate(&Task{}, &ScriptTaskDefinition{})
	if err != nil {
		panic(err.Error())
	}

	app.Wagon = StartWagon()

	task := &Task{
		Name: "TestTask",
		Schedule: &Schedule{
			ScheduleType: Interval,
			ScheduleDetail: &ScheduleDetail{
				Interval:     func(interval int) *int { return &interval }(4),
				IntervalKind: func(intervalKind time.Duration) *time.Duration { return &intervalKind }(time.Second),
			},
		},
	}
	db.Create(task)

	def := &ScriptTaskDefinition{
		Path: path.Join(homePathString, "test.sh"),
	}
	db.Create(def)
	task.TaskDefinitionID = &def.ID
	task.TaskDefinitionType = func(defType TaskDefinitionType) *TaskDefinitionType { return &defType }(Local)

	db.Save(task)

}

func GetConfig() *App {
	return app
}
