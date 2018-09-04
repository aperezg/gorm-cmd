package main

import (
	"github.com/aperezg/gorm-cmd"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {

	app := cli.NewApp()
	app.Name = "gorm-cmd"
	app.Usage = "WIP"
	app.Author = "Adrián Pérez <me@adrianpg.com>"

	app.Commands = []cli.Command{
		create(),
		migrate(),
		rollback(),
		redo(),
		version(),
	}

	app.Run(os.Args)
}

const (
	createCommandName   = "create"
	migrateCommandName  = "migrate"
	rollbackCommandName = "rollback"
	redoCommandName     = "redo"
	versionCommandName  = "version"
)

func create() cli.Command {
	return cli.Command{
		Name: createCommandName,
	}
}

func migrate() cli.Command {
	return cli.Command{
		Name: migrateCommandName,
		Action: func(c *cli.Context) error {
			db := gorm_cmd.OpenDB()
			defer db.Close()

			currentVersion := gorm_cmd.CurrentVersion(db)
			if err := gorm_cmd.Up("examples/", currentVersion.Version); err != nil {
				log.Fatal(err)
			}

			return nil
		},
	}
}

func rollback() cli.Command {
	return cli.Command{
		Name: rollbackCommandName,
	}
}

func redo() cli.Command {
	return cli.Command{
		Name: redoCommandName,
	}
}

func version() cli.Command {
	return cli.Command{
		Name: versionCommandName,
	}
}
