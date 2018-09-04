package gorm_cmd

import (
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"

	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type DBConfig struct {
	Driver string
	Conn   string
}

func OpenDB() *gorm.DB {
	var err error

	dbConfig := loadDBConfig()
	db, err := gorm.Open(dbConfig.Driver, dbConfig.Conn)
	if err != nil {
		panic(err)
	}

	//Create table version if not exists
	if err := initVersion(db); err != nil {
		log.Fatal("Can't create version table:", err)
	}

	return db
}

func loadDBConfig() *DBConfig {
	godotenv.Load()
	return &DBConfig{
		os.Getenv("GORM_DRIVER"),
		os.Getenv("GORM_CONN"),
	}
}
