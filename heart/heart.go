package heart

import (
	"tinygo.org/x/bluetooth"
)

var (
	adapter = bluetooth.DefaultAdapter
)

func Start() error {
	return adapter.Enable()
}

func Close() error {
	return nil
}
