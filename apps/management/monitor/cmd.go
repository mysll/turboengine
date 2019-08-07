package monitor

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"turboengine/core/api"
	"turboengine/core/service"

	"github.com/urfave/cli"
)

func run(srv api.Service) {

	app := cli.NewApp()
	app.Name = "management"
	app.Usage = "management tools"
	app.Version = "1.0.0"

	app.Commands = []cli.Command{
		{
			Name:  "shut",
			Usage: "shut id/all, shut one or all service",
			Action: func(c *cli.Context) (err error) {
				id := c.Args().Get(0)
				if id == "" {
					return fmt.Errorf("usage: %s", c.Command.Usage)
				}

				if id == "all" {
					err = srv.Pub(service.SERVICE_SHUT_ALL, []byte(""))
				} else {
					err = srv.Pub(service.SERVICE_SHUT, []byte(id))
				}

				if err != nil {
					return
				}
				fmt.Println("ok")
				return nil
			},
		},
		{
			Name:  "quit",
			Usage: "quit",
			Action: func(c *cli.Context) error {
				srv.Close()
				return errors.New("quit")
			},
		},
	}
	app.Action = func(c *cli.Context) error {
		for {
			fmt.Print("cmd> ")
			var input string
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			input = scanner.Text()

			cmdArgs := strings.Split(input, " ")
			if len(cmdArgs) == 0 {
				continue
			}
			s := []string{app.Name}
			s = append(s, cmdArgs...)
			err := c.App.Run(s)
			if err != nil {
				if err.Error() == "quit" {
					fmt.Println("bye")
					break
				}
				fmt.Println("err:", err)
			}
		}
		return nil
	}
	app.Run(os.Args)
}
