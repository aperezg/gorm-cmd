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

//initVersion Initialize table version on db
func initVersion(db *gorm.DB) error {
	if !db.HasTable(&Version{}) {
		return db.AutoMigrate(&Version{}).Error
	}
	return nil
}

//updateVersion Add or remove version by typeMigration
func updateVersion(tx *gorm.DB, version string, typeMigration string) error {

	if typeMigration == upMigrationType {
		return newVersion(tx, version)
	} else if typeMigration == downMigrationType {
		return rollbackVersion(tx)
	}

	return nil
}

//currentVersion get the last version
func CurrentVersion(db *gorm.DB) Version {
	var last Version
	db.Last(&last)
	return last
}

//ExctractVersionOfFile read the filename and extract the version number on format ('20060102150405')
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

func rollbackVersion(tx *gorm.DB) error {
	var last Version
	tx.Last(&last)
	return tx.Delete(last).Error
}
