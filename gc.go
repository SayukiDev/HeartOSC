package main

import "runtime/debug"

const mb = 1024 * 1024

func init() {
	debug.SetMemoryLimit(25 * mb)
}
