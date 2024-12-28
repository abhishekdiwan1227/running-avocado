package avo

import (
	"os"
	"path/filepath"

	"gorm.io/gorm"
)

func CleanDatabase(db *gorm.DB) {
	homePathString := os.Getenv("HOME")
	projectDir := filepath.Join(homePathString, ".avo")

	os.RemoveAll(projectDir)
}
