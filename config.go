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
	intervalPtr := execCommand.Int("interval", 60, "interval 60")
	verbosePtr := execCommand.Bool("verbose", false, "verbose")
	execCommand.Var(&myFlags, "var", "var FOO=my_ssm_parameter_key")

	switch args[0] {
	case "exec":
		execCommand.Parse(args[1:])
		config.Interval = *intervalPtr
		config.Verbose = *verbosePtr
		command := execCommand.Args()

		if len(myFlags) < 1 {
			return nil, nil, fmt.Errorf("No variables specified.")
		}

		if err := extractVariablesFromFlags(myFlags, config); err != nil {
			return nil ,nil, err
		}

		if len(command) < 1 {
			return nil, nil, fmt.Errorf("No command detected.")
		}

		return config, command, nil
	}

	return nil, nil, fmt.Errorf("%q is not a valid command.\n", args[0])
}

func extractVariablesFromFlags(myFlags arrayFlags, config *Config) error {
	for _, pair := range myFlags {

		kvp := strings.SplitN(pair, "=", 2)
		if len(kvp) < 2 {
			return fmt.Errorf("Could not parse variable %s", pair)
		}
		config.Variables[kvp[0]] = kvp[1]
	}
	return nil
}

type Config struct {
	Command string
	Variables map[string]string
	Interval int
	Verbose bool
}
