package models

import (
	"errors"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/spf13/viper"
	"minibox.ai/minibox/pkg/utils"
)

type Database struct {
	*gorm.DB
}

type ModelByInt struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

type Model struct {
	ID        string `gorm:"primary_key" auto:"uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

var db *gorm.DB

func SetDatabase(_db *gorm.DB) error {
	db = _db
	return nil
}

func Open(dialect, database string) error {
	_db, err := gorm.Open(dialect, database)
	if err != nil {
		return err
	}
	db = _db
	return nil
}

func Close() {
	if db != nil {
		db.Close()
	}
}

func LogMode(on bool) {
	db.LogMode(on)
}

func LoadConfig() error {
	dbconfig := viper.Sub("database")
	if dbconfig == nil {
		return errors.New("missing database section in minibox config")
	}
	adapter := dbconfig.GetString("adapter")
	switch adapter {
	case "sqlite3":
		database := dbconfig.GetString("database")
		return Open(adapter, database)
	case "mysql", "postgresql":
		url := dbconfig.GetString("url")
		if utils.Empty(url) {
			url = buildConnectString(dbconfig)
		}
		return Open(adapter, url)
	default:
		return errors.New("invalid database adpater")
	}
}

func GetDB() *Database {
	return &Database{db}
}

func buildConnectString(sub *viper.Viper) string {
	var (
		url      string
		database = sub.GetString("database")
		host     = sub.GetString("host")
		port     = sub.GetInt("port")
		user     = sub.GetString("username")
		pass     = sub.GetString("password")
		// pool     = sub.GetInt("pool")
	)
	if !utils.Empty(user) {
		url += user
		if !utils.Empty(pass) {
			url += ":" + pass
		}
		url += "@"
	}

	if !utils.Empty(host) {
		url += host
	}

	if port > 0 {
		url += strconv.Itoa(port)
	}

	if !utils.Empty(database) {
		url += "/" + database
	}

	// if pool > 0 {
	// 	url += strconv.Itao(port)
	// }

	return url
}
