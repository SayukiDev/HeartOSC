package options

import (
	"encoding/json"
	"os"
)

type Options struct {
	savePath           string
	Device             string `json:"Device"`
	OSCHost            string `json:"OSCHost"`
	OscPort            int    `json:"OscPort"`
	Parameter          string `json:"Parameter"`
	EnableRandomOffset bool   `json:"EnableRandomOffset"`
	EnableSmoothing    bool   `json:"EnableSmoothing"`
}

func New(path string) *Options {
	return &Options{
		savePath:           path,
		OSCHost:            "127.0.0.1",
		OscPort:            9000,
		Parameter:          "VRCOSC/Heartrate/Value",
		EnableRandomOffset: false,
		EnableSmoothing:    false,
	}
}

func (c *Options) Load() (isDefault bool, err error) {
	f, err := os.Open(c.savePath)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}
		return false, err
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(c)
	if err != nil {
		return false, err
	}
	return false, nil
}

func (c *Options) Save() error {
	f, err := os.OpenFile(c.savePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(c)
	if err != nil {
		return err
	}
	return nil
}
