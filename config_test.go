package main

import (
	"testing"
	"fmt"
)

func TestParseArgs_ExecWithVariable_ReturnsConfig(t *testing.T) {
	args := []string{
		"exec",
		"--var",
		"FOO=foo",
		"ping",
		"127.0.0.1",
	}

	config, command, err := parseArgs(args)

	if err != nil {
		t.Fatal(err.Error())
	}

	assertEqual(t, len(config.Variables), 1, "")
	assertEqual(t, config.Variables["FOO"], "foo", "")

	assertEqual(t, len(command), 2, "")
	assertEqual(t, command[0], "ping", "")
	assertEqual(t, command[1], "127.0.0.1", "")

}

func TestParseArgs_WithNoVariables_ReturnsError(t *testing.T) {
	args := []string{
		"exec",
		"ping",
		"127.0.0.1",
	}

	_, _, err := parseArgs(args)

	if err == nil {
		t.Error("Expected an error ")
	}
}

func TestParseArgs_WithNoCommand_ReturnsError(t *testing.T) {
	args := []string{
		"exec",
		"-var",
		"FOO=foo",
	}

	_, _, err := parseArgs(args)

	if err == nil {
		t.Error("Expected an error ")
	}
}

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}
