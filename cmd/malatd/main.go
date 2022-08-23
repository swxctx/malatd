package main

import (
	"os"

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

	app.Commands = []cli.Command{newCom}
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
