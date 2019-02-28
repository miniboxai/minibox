package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Database struct {
	*gorm.DB
}

func Open(dialect string, args ...interface{}) (*Database, error) {
	gdb, err := gorm.Open(dialect, args...)
	db := &Database{
		DB: gdb,
	}
	return db, err
}
