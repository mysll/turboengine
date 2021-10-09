package main

import (
	"os"
	"turboengine/apps/tools/turbogen"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "turbogen"
	app.Usage = "turbo engine toolkit"

	app.Commands = []cli.Command{
		{
			Name:  "create",
			Usage: "create (service | module | proto | entity) --path output",

			Subcommands: []cli.Command{
				{
					Name:   "service",
					Usage:  "create service",
					Action: turbogen.CreateService,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "path,p",
							Value: "./",
							Usage: "output path",
						},
					},
				},
				{
					Name:   "module",
					Usage:  "create module",
					Action: turbogen.CreateModule,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "path,p",
							Value: "./",
							Usage: "output path",
						},
					},
				},
				{
					Name:   "proto",
					Usage:  "create proto",
					Action: turbogen.CreateProto,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "path,p",
							Value: "./",
							Usage: "output path",
						},
					},
				},
				{
					Name:   "entity",
					Usage:  "create entity",
					Action: turbogen.CreateEntity,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "path,p",
							Value: "./",
							Usage: "output path",
						},
					},
				},
			},
		},
	}
	app.Run(os.Args)
}
