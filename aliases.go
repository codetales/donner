package main

import (
	"fmt"
	"strings"
)

func printAliases(cfg *Cfg, strictMode, fallbackMode bool) {
	commands := cfg.ListCommands()
	outputs := make([]string, len(commands))

	flags := make([]string, 0, 2)
	if strictMode {
		flags = append(flags, "--strict")
	}

	if fallbackMode {
		flags = append(flags, "--fallback")
	}

	for i, c := range commands {
		output := append(flags, c)
		outputs[i] = strings.Join(output, " ")
	}

	fmt.Println()
	for i, c := range commands {
		fmt.Printf("alias %s='donner run %s';\n", c, outputs[i])
	}

	aliasCommand := strings.Join(append([]string{"donner", "aliases"}, flags...), " ")

	fmt.Printf("\n# copy and paste the output into your terminal or run\n")
	fmt.Printf("#  eval $(%s)\n", aliasCommand)
}
