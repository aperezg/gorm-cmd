package gorm_cmd

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

func Run() {
	log.SetOutput(os.Stdout)

	if execMigration == nil {
		log.Fatal("Migration to execute not registered")
	}

	db := OpenDB()
	defer db.Close()

	migrationType := os.Args[1]
	if migrationType == upMigrationType {
		err := execMigration.Up(db)
		if err != nil {
			log.Fatal(err)
		}
	} else if migrationType == downMigrationType {
		err := execMigration.Down(db)
		if err != nil {
			log.Fatal(err)
		}
	}
}


