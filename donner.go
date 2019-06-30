package main

import (
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

// ErrUndefinedCommand is thrown if a command specified can't be found in the yaml definition
var ErrUndefinedCommand = errors.New("the command you're trying to run doesn't exist in the yaml definition")

// ErrMissingCommand is thrown if no handler for execution is provided
var ErrMissingCommand = errors.New("no command for execution specified")

func main() {
	app := cli.NewApp()
	app.Name = "Donner"
	app.Usage = `Donner is a generic command wrapper. It let's you define strategies to wrap commands in things like 'docker-compose exec' or 'docker container run'.
	 This is can come in very handy when developing applications in containers. Donner allows defining a wrapping strategy on a per command basis. 
	 So you don't have to worry which service to use or whether you should use 'docker-compose exec' or 'docker-compose run' when executing a command.`
	// TODO implement flags
	app.Flags = []cli.Flag{
		cli.BoolFlag{Name: "strict,s", Usage: "enable strict mode"},
		cli.BoolFlag{Name: "fallback,f", Usage: "fallback to local commands"},
	}
	app.Commands = []cli.Command{
		{
			Name:            "run",
			Aliases:         []string{"r"},
			Usage:           "run a command",
			SkipFlagParsing: true,
			Action: func(c *cli.Context) error {
				// TODO handle 'yaml' case
				dat, err := ioutil.ReadFile(".donner.yml")
				if err != nil {
					return err
				}

				cfg, err := parseFile(dat)
				if err != nil {
					return err
				}

				return execCommand(cfg, c.Args())
			},
		},
		{
			Name:    "aliases",
			Aliases: []string{"a"},
			Usage:   "generate aliases",
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

// execCommand dispatches the call to the OS
func execCommand(cfg *Cfg, params []string) error {
	if len(params) < 1 {
		// TODO show usage?
		return ErrMissingCommand
	}

	// check if command specified exists in parsed file
	command, ok := cfg.Commands[params[0]]
	if !ok {
		return ErrUndefinedCommand
	}

	// extract handler from corresponding strategy
	designatedHandler := cfg.Strategies[string(command)]
	execHandler := availableHandlers[designatedHandler.Handler]

	// construct os call
	cliArgs := []string{execHandler.Argument}

	if designatedHandler.Remove {
		cliArgs = append(cliArgs, "--rm")
	}

	if designatedHandler.Service != "" {
		cliArgs = append(cliArgs, designatedHandler.Service)
	}

	if designatedHandler.Image != "" {
		cliArgs = append(cliArgs, designatedHandler.Image)
	}

	for _, p := range params {
		cliArgs = append(cliArgs, p)
	}

	cmd := exec.Command(execHandler.BaseCommand, cliArgs...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	// TODO structured output?
	fmt.Println(string(out))

	return nil
}
