package app

import (
	"HeartOSC/options"
	"HeartOSC/service"
	"flag"

	tea "charm.land/bubbletea/v2"
)

var path = flag.String("path", "./config.json", "config path")

func start() (s *service.Service, isDefault bool, err error) {
	flag.Parse()
	c := options.New(*path)
	isDefault, err = c.Load()
	if err != nil {
		return
	}
	s, err = service.New(c)
	if err != nil {
		return
	}
	err = s.Start()
	if err != nil {
		return
	}
	return s, isDefault, nil
}

type startedMsg struct {
	Srv *service.Service
}

type startModel struct {
}

func newStartModel() startModel {
	return startModel{}
}

func (s *startModel) Init() tea.Cmd {
	return tea.Batch(func() tea.Msg {
		srv, _, err := start()
		if err != nil {
			return err
		}
		return startedMsg{Srv: srv}
	})
}

func (s *startModel) View() string {
	return "起動中..."
}
