package main

import (
	"fmt"
	"github.com/aperezg/gorm-cmd"
	"github.com/jinzhu/gorm"
)

func main() {
	gorm_cmd.AddMigrationToExec(Up_20180830235959, Down_2018030235959)
	gorm_cmd.Run()
}

func Up_20180830235959(tx *gorm.DB) error {
	fmt.Println("Executing Up_20180830235959")
	return nil
}

func Down_2018030235959(tx *gorm.DB) error {
	fmt.Println("Executing Down_2018030235959")
	return nil
}
