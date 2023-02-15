package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func getOsEnvironments() map[string]string {
	osEnvs := os.Environ()
	resp := make(map[string]string, len(osEnvs))

	for _, env := range os.Environ() {
		vars := strings.Split(env, "=")
		key := vars[0]
		val := vars[1]
		resp[key] = val
	}

	return resp
}

func joinEnvironments(osEnvs map[string]string, appEnvs Environment) {
	for key, val := range appEnvs {
		if val.NeedRemove {
			delete(osEnvs, key)
		} else {
			osEnvs[key] = val.Value
		}
	}
}

func prepareEnvironments(cmdEnvs map[string]string) []string {
	resp := make([]string, 0, len(cmdEnvs))
	for key, val := range cmdEnvs {
		resp = append(resp, key+"="+val)
	}

	return resp
}

func runCmd(cmdName string, args, env []string) error {
	cmd := exec.Command(cmdName, args...)
	cmd.Env = env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// RunCmd runs a command + arguments (commands) with environment variables from env.
func RunCmd(args []string, env Environment) (returnCode int) {
	cmd := args[0]
	var cmdArgs []string
	if len(args) > 1 {
		cmdArgs = args[1:]
	}

	osEnvs := getOsEnvironments()
	joinEnvironments(osEnvs, env)
	cmdEnvs := prepareEnvironments(osEnvs)

	if err := runCmd(cmd, cmdArgs, cmdEnvs); err != nil {
		fmt.Println(err)
		return 1
	}

	return 0
}

func getEnvDirName() (string, error) {
	if len(os.Args) > 1 {
		return os.Args[1], nil
	}

	return "", errors.New("env directory name not send")
}

func getCmdArguments() []string {
	if len(os.Args) > 2 {
		return os.Args[2:]
	}

	return nil
}

func Execute() error {
	envDir, err := getEnvDirName()
	if err != nil {
		return err
	}

	args := getCmdArguments()

	envs, err := ReadDir(envDir)
	if err != nil {
		panic(err)
	}

	fmt.Println(RunCmd(args, envs))

	return nil
}
