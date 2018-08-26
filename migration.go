package gorm_cmd

import (
	"github.com/jinzhu/gorm"
	"log"
	"path/filepath"
)

type Migration struct {
	Version  string
	Next     string
	Previous string
	Source   string
	UpFn     func(db *gorm.DB) error
	DownFn   func(db *gorm.DB) error
}

const (
	upMigrationType   = "Up"
	downMigrationType = "Down"
)

func (m *Migration) Up(db *gorm.DB) error {
	if err := m.exec(db, upMigrationType); err != nil {
		return err
	}

	log.Println(upMigrationType, ":", filepath.Base(m.Source))
	return nil
}

func (m *Migration) Down(db *gorm.DB) error {
	if err := m.exec(db, downMigrationType); err != nil {
		return err
	}

	log.Println(downMigrationType, ":", filepath.Base(m.Source))
	return nil
}

func (m *Migration) exec(db *gorm.DB, typeMigration string) error {

	//Create table version if not exists
	err := InitVersion(db)
	if err != nil {
		log.Fatal("Can't create version table:", err)
	}

	tx := db.Begin()
	if tx.Error != nil {
		log.Fatal("Error starting transaction: ", tx.Error)
	}

	fn := m.UpFn
	if typeMigration == downMigrationType {
		fn = m.DownFn
	}

	if fn == nil {
		log.Fatal("The function type", typeMigration, "is not defined")
	}

	if err := m.UpFn(tx); err != nil {
		tx.Rollback()
		log.Fatalf("Error executing migration %s :%v", filepath.Base(m.Source), err)
		return err
	}

	if err := UpdateVersion(tx, m.Version, typeMigration); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
