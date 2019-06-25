package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Donner"
	app.Usage = `Donner is a generic command wrapper. It let's you define strategies to wrap commands in things like 'docker-compose exec' or 'docker container run'. 
	This is can come in very handy when developing applications in containers. Donner allows defining a wrapping strategy on a per command basis. 
	So you don't have to worry which service to use or whether you should use 'docker-compose exec' or 'docker-compose run' when executing a command.`
	app.Commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "run a command",
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "strict,s", Usage: "enable strict mode"},
				cli.BoolFlag{Name: "fallback,f", Usage: "fallback to local commands"},
			},
			Action: func(c *cli.Context) error {
				command := []string{"docker-compose", "exec", "app"}
				command = append(command, c.Args()...)

				// TODO handle 'yaml' case
				dat, err := ioutil.ReadFile(".donner.yml")
				if err != nil {
					return err
				}

				res, err := ParseFile(dat)
				if err != nil {
					return err
				}

				// TODO dispatch os call
				fmt.Printf("%+v", res)

				return nil
			},
		},
		{
			Name:    "aliases",
			Aliases: []string{"a"},
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "strict,s", Usage: "enable strict mode"},
				cli.BoolFlag{Name: "fallback,f", Usage: "fallback to local commands"},
			},
			Usage: "generate aliases",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
