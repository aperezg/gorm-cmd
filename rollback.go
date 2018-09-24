package gorm_cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func Down(path, currentVersion string) error {
	err := getSpecificMigration(path, currentVersion)
	if err != nil {
		return err
	}

	if  len(registeredMigrations) <= 0 {
		log.Println("nothing to execute")
		return nil
	}

	m := registeredMigrations[0]
	o, err := exec.Command("go", "run", m.Source, downMigrationType).Output()
	if err != nil {
		return err
	}

	fmt.Println(string(o))
	return nil
}

func getSpecificMigration(path, version string) error {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("%s not exists", path)
	}

	migrationFiles, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	migrations := Migrations{}
	for _, f := range migrationFiles {
		v, err := extractVersionOfFile(f.Name())
		if err != nil {
			return err
		}

		if v != version {
			continue
		}

		migrations = append(migrations, &Migration{Version: v, Source: path + f.Name()})
	}

	registeredMigrations = migrations
	return nil
}
