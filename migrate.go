package gorm_cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sort"
	"time"
)

//Up execute all next migrations
func Up(path, currentVersion string) error {
	err := getMigrations(path, currentVersion)
	if err != nil {
		return err
	}

	for _, m := range registeredMigrations {
		o, err := exec.Command("go", "run", m.Source, upMigrationType).Output()
		if err != nil {
			return err
		}

		fmt.Println(string(o))
	}

	if len(registeredMigrations) <= 0 {
		log.Println("Nothing to execute")
	}

	return nil
}

func getMigrations(path, currentVersion string) error {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("%s not exists", path)
	}

	migrationFiles, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	migrations := Migrations{}
	vmigrations, err := getValidMigrations(migrationFiles, currentVersion)

	if err != nil {
		return err
	}

	for _, v := range vmigrations {
		for _, f := range migrationFiles {
			version, _ := extractVersionOfFile(f.Name())
			if version == v {
				migrations[v] = &Migration{Version: v, Source: path + f.Name()}
			}
		}
	}
	registeredMigrations = migrations
	return nil
}

func getValidMigrations(files []os.FileInfo, currentVersion string) ([]string, error) {
	var migrations sort.StringSlice
	currentVersionTime, _ := time.Parse(migrationFileFormatDate, currentVersion)

	for _, f := range files {
		version, err := extractVersionOfFile(f.Name())
		if err != nil {
			return []string{}, err
		}

		vt, err := time.Parse(migrationFileFormatDate, version)
		if err != nil {
			return []string{}, err
		}

		if vt.Before(currentVersionTime) || vt.Equal(currentVersionTime) {
			continue
		}

		migrations = append(migrations, version)
	}

	sort.Sort(sort.Reverse(migrations[:]))
	return migrations, nil
}