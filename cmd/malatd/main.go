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
		Name:   "run",
		Usage:  "Compile and run go project",
		Before: initProject,
		Action: func(c *cli.Context) error {
			run.RunProject()
			return nil
		},
	}

	// new a README.md
	newDocCom := cli.Command{
		Name:  "doc",
		Usage: "Generate a project README.md(malatd doc || malatd doc -r ${root_group})",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "root_group, r",
				Usage: "The project name uri or input",
			},
		},
		Before: initProject,
		Action: func(c *cli.Context) error {
			create.CreateDoc(c.String("r"))
			return nil
		},
	}

	app.Commands = []cli.Command{newCom, runCom, newDocCom}
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
