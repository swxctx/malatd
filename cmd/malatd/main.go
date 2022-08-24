package main

import (
	"os"

	"github.com/swxctx/malatd/cmd/run"

	"github.com/swxctx/malatd/cmd/create"
	"github.com/swxctx/malatd/cmd/info"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Malatd project command"
	app.Version = "1.0.0"
	app.Author = "swxctx"
	app.Usage = "A deployment tools of malatd frameware"

	// new a project
	newCom := cli.Command{
		Name:   "gen",
		Usage:  "Generate a malatd project",
		Before: initProject,
		Action: func(c *cli.Context) error {
			create.CreateProject()
			return nil
		},
	}

	// run a project
	runCom := cli.Command{
		Name:      "run",
		Usage:     "Compile and run go project",
		UsageText: `malatd run`,
		Flags: []cli.Flag{
			cli.StringSliceFlag{
				Name:  "watch_exts, x",
				Value: (*cli.StringSlice)(&[]string{".go", ".ini", ".yaml", ".toml", ".xml"}),
				Usage: "Specified to increase the listening file suffix",
			},
			cli.StringSliceFlag{
				Name:  "notwatch, n",
				Value: (*cli.StringSlice)(&[]string{}),
				Usage: "Not watch files or directories",
			},
			cli.StringFlag{
				Name:  "app_path, p",
				Usage: "The path(relative/absolute) of the project",
			},
		},
		Before: initProject,
		Action: func(c *cli.Context) error {
			run.RunProject()
			return nil
		},
	}

	app.Commands = []cli.Command{newCom, runCom}
	app.Run(os.Args)
}

// initProject 初始化项目
func initProject(c *cli.Context) error {
	appPath := c.String("app_path")
	if len(appPath) == 0 {
		appPath = c.Args().First()
	}
	if len(appPath) == 0 {
		appPath = "./"
	}
	return info.Init(appPath)
}
