//go:build windows

package main

import (
	"os"
	"path/filepath"
)

func clean() {
	d, err := os.ReadDir("./")
	if err != nil {
		panic(err)
	}
	for _, f := range d {
		if filepath.Ext(f.Name()) == ".syso" {
			os.Remove(f.Name())
		}
	}
}
