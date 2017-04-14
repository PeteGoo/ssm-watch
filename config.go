package main

import (
	"flag"
	"fmt"
	"strings"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join(*i, ",")
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func parseArgs(args []string) (*Config, []string, error) {
	config := &Config{
		Variables: make(map[string]string),
	}

	var myFlags arrayFlags
	execCommand := flag.NewFlagSet("exec", flag.ContinueOnError)
	execCommand.Var(&myFlags, "var", "var FOO=my_ssm_parameter_key")

	switch args[0] {
	case "exec":
		execCommand.Parse(args[1:])
		command := execCommand.Args()

		if len(myFlags) < 1 {
			return nil, nil, fmt.Errorf("No variables specified.")
		}

		for _, pair := range myFlags {

			kvp := strings.SplitN(pair, "=", 2)
			if len(kvp) < 2 {
				return nil, nil, fmt.Errorf("Could not parse variable %s", pair)
			}
			config.Variables[kvp[0]] = kvp[1]
		}

		if len(command) < 1 {
			return nil, nil, fmt.Errorf("No command detected.")
		}

		return config, command, nil
	}

	return nil, nil, fmt.Errorf("%q is not a valid command.\n", args[0])
}

type Config struct {
	Command string
	Variables map[string]string
}
