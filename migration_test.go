package gorm_cmd

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	db := OpenDB()
	cleanVersioning(db)

	os.Exit(m.Run())
}

func TestAddMigrationToExec(t *testing.T) {
	AddMigrationToExec(up_test_20180830235959, down_test_2018030235959)
	assert.NotNil(t, execMigration)
}

func TestMigration_Up(t *testing.T) {
	version := "20180830235959"
	db := OpenDB()
	defer cleanVersioning(db)

	migration := &Migration{UpFn: up_test_20180830235959, Version: version}
	err := migration.Up(db)
	v := CurrentVersion(db)

	assert.NoError(t, err)
	assert.Equal(t, version, v.Version)
}

func TestMigration_Down(t *testing.T) {
	version := "20180830235959"

	db := OpenDB()
	defer cleanVersioning(db)
	newVersion(db, version)

	migration := &Migration{DownFn:down_test_2018030235959, Version: version}
	err := migration.Down(db)
	v := CurrentVersion(db)

	assert.NoError(t, err)
	assert.NotEqual(t, version, v.Version)
}

func TestMigration_Up_Error(t *testing.T) {
	var data = []struct {
		m *Migration
		db *gorm.DB
	}{
		{new(Migration), dbWithError()}, //Wrong DB
		{&Migration{UpFn: wrong_fn}, OpenDB()}, // Wrong Run Upfn
		{&Migration{UpFn: up_test_20180830235959, Version: "20180830235959"}, OpenDB()}, //Duplicate Version
	}

	for _, d := range data {
		if d.db != nil {
			cleanVersioning(d.db)
			newVersion(d.db, "20180830235959")
		}

		err := d.m.Up(d.db)
		d.db.Close()
		assert.Error(t, err)
	}
}

func TestMigration_Down_Error(t *testing.T) {
	var data = []struct {
		m *Migration
		db *gorm.DB
	} {
		{new(Migration), dbWithError()}, //Wrong DB
		{&Migration{DownFn: wrong_fn}, OpenDB()}, // Wrong Run Downfn
	}

	for _, d := range data {
		if d.db != nil {
			cleanVersioning(d.db)
		}

		err := d.m.Down(d.db)
		d.db.Close()
		assert.Error(t, err)
	}
}

func cleanVersioning(db *gorm.DB) {
	db.Delete(&Version{})
}

func wrong_fn(tx *gorm.DB) error {
	return errors.New("FAIL")
}

func up_test_20180830235959(tx *gorm.DB) error {
	fmt.Println("Executing Up_20180830235959")
	return nil
}

func down_test_2018030235959(tx *gorm.DB) error {
	fmt.Println("Executing Down_2018030235959")
	return nil
}

func dbWithError() *gorm.DB {
	db := OpenDB()
	db.AddError(errors.New("FAIL"))

	return db
}
