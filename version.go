package gorm_cmd

import (
	"errors"
	"github.com/jinzhu/gorm"
	"path/filepath"
	"time"
)

type Version struct {
	Version string `gorm:"primary_key"`
}

const migrationFileFormatDate = "20060102150405"

//currentVersion get the last version
func CurrentVersion(db *gorm.DB) Version {
	var last Version
	db.Last(&last)
	return last
}

func initVersion(db *gorm.DB) error {
	if !db.HasTable(&Version{}) {
		return db.AutoMigrate(&Version{}).Error
	}
	return nil
}

func updateVersion(tx *gorm.DB, version string, typeMigration string) error {

	if typeMigration == upMigrationType {
		return newVersion(tx, version)
	} else if typeMigration == downMigrationType {
		return rollbackVersion(tx, version)
	}

	return nil
}

func extractVersionOfFile(name string) (string, error) {
	file := filepath.Base(name)

	if chkExt := filepath.Ext(file); chkExt != ".go" {
		return "", errors.New("migration file type not recognized")
	}

	v := file[0:14]
	_, err := time.Parse(migrationFileFormatDate, v)
	if err != nil {
		return "", errors.New("migration file must start with a valid format date (" + migrationFileFormatDate + ")")
	}

	return v, nil
}

func newVersion(tx *gorm.DB, version string) error {
	return tx.Create(&Version{version}).Error
}

func rollbackVersion(tx *gorm.DB, version string) error {
	if db := tx.Delete(&Version{version}); db != nil {
		return db.Error
	}

	return nil
}
