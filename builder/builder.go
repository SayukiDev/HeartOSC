package main

import (
	"os"
	"os/exec"
)

func main() {
	err := ExecStd(nil, "go", "generate")
	if err != nil {
		panic(err)
	}
	err = ExecStd(
		nil,
		"go",
		"build",
		"-o",
		"HeartOSC.exe",
		"-trimpath",
		"-ldflags=-s -w",
	)
	if err != nil {
		panic(err)
	}
	clean()
}

func ExecStd(env []string, name string, params ...string) error {
	cmd := exec.Command(name, params...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env
	return cmd.Run()
}
