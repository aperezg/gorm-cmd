package gorm_cmd

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddMigrationToExec(t *testing.T) {
	AddMigrationToExec(up_test_20180830235959, down_test_2018030235959)
	assert.NotNil(t, execMigration)
}

func TestMigration_Up_Wrong_Run_Upfn(t *testing.T) {
	db := OpenDB()
	migration := &Migration{UpFn: wrong_fn}

	err := migration.Up(db)
	assert.Error(t, err)
}

func TestMigration_Up_Wrong_Version(t *testing.T) {
	db := OpenDB()
	cleanVersioning(db)
	newVersion(db, "20180830235959")

	migration := &Migration{UpFn: up_test_20180830235959, Version: "20180830235959"}
	err := migration.Up(db)
	assert.Error(t, err)
}

func TestMigration_Up(t *testing.T) {
	db := OpenDB()
	cleanVersioning(db)

	migration := &Migration{UpFn: up_test_20180830235959, Version: "20180830235959"}
	err := migration.Up(db)
	assert.NoError(t, err)
}

func TestMigration_Up_Error_DB(t *testing.T) {
	db := gorm.DB{}
	migration := &Migration{}

	err := migration.Up(&db)
	assert.Error(t, err)
}

func cleanVersioning(db *gorm.DB) {
	db.Delete(&Version{})
}

func wrong_fn(tx *gorm.DB) error {
	return errors.New("Fail")
}

func up_test_20180830235959(tx *gorm.DB) error {
	fmt.Println("Executing Up_20180830235959")
	return nil
}

func down_test_2018030235959(tx *gorm.DB) error {
	fmt.Println("Executing Down_2018030235959")
	return nil
}
