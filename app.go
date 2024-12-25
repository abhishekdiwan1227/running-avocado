package avo

import (
	"errors"
	"os"
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
}

func GetConfig() *App {
	return app
}
