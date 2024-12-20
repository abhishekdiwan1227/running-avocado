package main

import (
	"errors"
	"os"
	"path/filepath"

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

	avo.Init(db)
}
