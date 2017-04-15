package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	log "github.com/Sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	// log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	// log.SetLevel(log.DebugLevel)
}


func main() {
	config, command, err := parseArgs(os.Args[1:])

	if err != nil {
		log.WithError(err)
		os.Exit(1)
	}

	if config.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	env := environ(os.Environ())

	updateParameters(config, &env)

	cmd := exec.Command(command[0], command[1:]...)

	ticker := time.NewTicker(time.Second * time.Duration(config.Interval))

	go func() {
		for range ticker.C {
			if changed := updateParameters(config, &env); changed {
				// TODO: Signal child process
			}
		}
	}()

	cmd.Env = env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	var waitStatus syscall.WaitStatus

	if err := cmd.Run(); err != nil {

		ticker.Stop()

		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			os.Exit(waitStatus.ExitStatus())
		}
		if err != nil {
			log.WithError(err)
			os.Exit(1)
		} /**/

	}
}

func updateParameters(config *Config, env *environ) bool {
	for k, v := range config.Variables {

		log.Debugf("Checking SSM parameter %s\n", k)

		value, err := getSsmParameter(v)
		if err != nil {
			log.WithError(err)
			os.Exit(1)
		}

		if !env.Exists(k) {
			log.Debugf("Setting %s for the first time\n", k)
			env.Set(k, value)
			return true
		} else if !env.IsSame(k, value) {
			log.Debugf("Updating to new value for %s\n", k)
			env.Set(k, value)
			return true
		}
	}
	return false
}

// environ is a slice of strings representing the environment, in the form "key=value".
type environ []string

func getSsmParameter(key string) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")})

	if err != nil {
		return "", err
	}

	client := ssm.New(sess)
	//result, err := client.DescribeParameters(&ssm.DescribeParametersInput{})
	result, err := client.GetParameters(&ssm.GetParametersInput{
		Names:          []*string{&key},
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", err
	}
	if len(result.InvalidParameters) > 0 {

		return "", fmt.Errorf("Invalid parameter: %s", key)
	}
	return *result.Parameters[0].Value, nil
}

// Unset an environment variable by key
func (e *environ) Unset(key string) {
	for i := range *e {
		if strings.HasPrefix((*e)[i], key+"=") {
			(*e)[i] = (*e)[len(*e)-1]
			*e = (*e)[:len(*e)-1]
			break
		}
	}
}

// Set adds an environment variable, replacing any existing ones of the same key
func (e *environ) Set(key, val string) {
	e.Unset(key)
	*e = append(*e, key+"="+val)
}

func (e *environ) Exists(key string) bool{
	for _, val := range *e {
		if strings.HasPrefix(val, key + "=") {
			return true
		}
	}
	return false
}

func (e *environ) IsSame(key string, value string) bool{
	for _, val := range *e {
		if val == key+"="+value {
			return true
		}
	}
	return false
}

