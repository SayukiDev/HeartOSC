package service

import (
	"HeartOSC/heart"
	"HeartOSC/options"
	"HeartOSC/sender"
)

type Service struct {
	Sender *sender.Sender
	Conf   *options.Options
}

func New(c *options.Options) (*Service, error) {
	s := sender.New(c.OSCHost, c.OscPort, c.Parameter, c.EnableRandomOffset, c.EnableSmoothing)
	return &Service{
		Sender: s,
		Conf:   c,
	}, nil
}

func (s *Service) Start() error {
	return heart.Start()
}

func (s *Service) Close() error {
	return heart.Close()
}

func (s *Service) StartSending() error {
	return s.Sender.Start()
}

func (s *Service) StopSending() error {
	return s.Sender.Close()
}
