package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	oS := []string{runtime.GOOS}
	arch := []string{runtime.GOARCH}
	if len(os.Args) >= 2 {
		oS = []string{os.Args[1]}
	}
	if len(os.Args) >= 3 {
		arch = []string{os.Args[2]}
	}

	err := os.MkdirAll("output", 0755)
	if err != nil {
		panic(err)
	}
	for _, a := range arch {
		for _, o := range oS {
			err := ExecStd([]string{
				"GOOS=" + o,
				"GOARCH=" + a,
			}, "go", "generate")
			if err != nil {
				panic(fmt.Errorf("generate %s %s error: %s", o, a, err))
			}
			name := fmt.Sprintf("output/HeartOSC_%s_%s/HeartOSC", o, a)
			if o == "windows" {
				name += ".exe"
			}
			err = ExecStd(
				[]string{
					"CGO_ENABLED=1",
					"GOOS=" + o,
					"GOARCH=" + a,
				},
				"go",
				"build",
				"-o",
				name,
				"-trimpath",
				"-ldflags=-s -w",
			)
			if err != nil {
				panic(fmt.Errorf("build %s %s error: %s", o, a, err))
			}
		}
	}
}

func ExecStd(env []string, name string, params ...string) error {
	cmd := exec.Command(name, params...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = append(os.Environ(), env...)
	return cmd.Run()
}
