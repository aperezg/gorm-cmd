package gorm_cmd

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"path/filepath"
	"runtime"
)

type Migration struct {
	Version  string
	Next     string
	Previous string
	Source   string
	UpFn     func(db *gorm.DB) error
	DownFn   func(db *gorm.DB) error
}

type Migrations map[string]*Migration

var registeredMigrations Migrations
var execMigration *Migration

const (
	upMigrationType   = "Up"
	downMigrationType = "Down"
)

//Up execute external code of migration
func (m *Migration) Up(db *gorm.DB) error {
	if err := m.exec(db, upMigrationType); err != nil {
		return err
	}

	log.Println("OK:", upMigrationType, "-", m.Version)
	return nil
}

//Down rollback external code of migration
func (m *Migration) Down(db *gorm.DB) error {
	if err := m.exec(db, downMigrationType); err != nil {
		return err
	}

	log.Println(downMigrationType, ":", filepath.Base(m.Source))
	return nil
}

func AddMigrationToExec(up func(tx *gorm.DB) error, down func(tx *gorm.DB) error) {
	_, file, _, _ := runtime.Caller(1)
	version, _ := extractVersionOfFile(file)

	execMigration = &Migration{UpFn: up, DownFn: down, Version: version}
}

func (m *Migration) exec(db *gorm.DB, typeMigration string) error {
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("error starting transaction: %v", tx.Error)
	}

	fn := m.UpFn
	if typeMigration == downMigrationType {
		fn = m.DownFn
	}

	if fn == nil {
		return fmt.Errorf("the function type %s is not defined", typeMigration)
	}

	if err := m.UpFn(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("error executing migration %s :%v", filepath.Base(m.Source), err)
	}

	if err := updateVersion(tx, m.Version, typeMigration); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
