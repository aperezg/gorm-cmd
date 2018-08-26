package main

import (
	"github.com/urfave/cli"
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
		Name:    createCommandName,
		Aliases: []string{"c"},
	}

}

func migrate() cli.Command {
	return cli.Command{
		Name:    migrateCommandName,
		Aliases: []string{"m"},
	}
}

func rollback() cli.Command {
	return cli.Command{
		Name:    rollbackCommandName,
		Aliases: []string{"b"},
	}
}

func redo() cli.Command {
	return cli.Command{
		Name:    redoCommandName,
		Aliases: []string{"r"},
	}
}

func version() cli.Command {
	return cli.Command{
		Name:    versionCommandName,
		Aliases: []string{"v"},
	}
}
