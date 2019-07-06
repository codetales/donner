package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/urfave/cli"
)

// ErrMissingCommand is thrown if no handler for execution is provided
var ErrMissingCommand = errors.New("no command for execution specified")

func main() {
	app := cli.NewApp()
	app.Name = "Donner"
	app.Usage = `Donner is a generic command wrapper. It let's you define strategies to wrap commands in things like 'docker-compose exec' or 'docker container run'.
	 This is can come in very handy when developing applications in containers. Donner allows defining a wrapping strategy on a per command basis. 
	 So you don't have to worry which service to use or whether you should use 'docker-compose exec' or 'docker-compose run' when executing a command.`
	app.Commands = []cli.Command{
		{
			Name:           "run",
			Aliases:        []string{"r"},
			Usage:          "run a command",
			SkipArgReorder: true,
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "strict,s", Usage: "enable strict mode"},
			},
			Action: func(c *cli.Context) error {
				cfg, err := readConfig()
				if err != nil {
					return err
				}
				return execCommand(cfg, c.Args(), c.Bool("strict"))
			},
		},
		{
			Name:    "aliases",
			Aliases: []string{"a"},
			Usage:   "generate aliases",
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "strict,s", Usage: "enable strict mode"},
				cli.BoolFlag{Name: "fallback,f", Usage: "fallback to local commands"},
			},
			Action: func(c *cli.Context) error {
				cfg, err := readConfig()
				if err != nil {
					return err
				}
				printAliases(cfg, c.Bool("strict"), c.Bool("fallback"))

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
func execCommand(cfg *Cfg, command []string, strict bool) error {
	if len(command) < 1 {
		// TODO show usage?
		return ErrMissingCommand
	}

	handler, err := cfg.GetHandlerFor(command[0], strict)
	if err != nil {
		return err
	}
	wrappedCommand := handler.WrapCommand(command)

	cmd := exec.Command(wrappedCommand[0], wrappedCommand[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// This block ensures that we return the same exit code in case the command failed
			waitStatus := exitError.Sys().(syscall.WaitStatus)
			os.Exit(waitStatus.ExitStatus())
		} else {
			// This block handles the case where Donner could not start the command at all (missing, permission, ...)
			_, _ = os.Stderr.WriteString(fmt.Sprintf("%s\n", err.Error()))
			os.Exit(1) // TODO: What would be a good exit code to return here?
		}
	}

	return nil
}

func readConfig() (*Cfg, error) {
	// TODO handle 'yaml' case
	dat, err := ioutil.ReadFile(".donner.yml")
	if err != nil {
		return nil, err
	}

	cfg, err := GenerateConfig(dat)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
