package sender

import (
	"HeartOSC/heart"
	"path"
	"time"

	"github.com/hypebeast/go-osc/osc"
)

const ParamPrefix = "/avatar/parameters/"

type Sender struct {
	c                  *osc.Client
	closeC             chan struct{}
	rate               int32
	smoothed           float64
	smoothInit         bool
	lastUpdate         time.Time
	ParamName          string
	EnableRandomOffset bool
	EnableSmoothing    bool
}

func New(addr string, port int, paramName string, enableFilter bool, enableSmoothing bool) *Sender {
	return &Sender{
		c:                  osc.NewClient(addr, port),
		ParamName:          path.Join(ParamPrefix, paramName),
		EnableRandomOffset: enableFilter,
		EnableSmoothing:    enableSmoothing,
	}
}

func (s *Sender) Send(value int32) error {
	return s.c.Send(osc.NewMessage(s.ParamName, value))
}

func (s *Sender) Start() error {
	s.closeC = make(chan struct{})
	go func() {
		for ; ; <-time.Tick(time.Second * 1) {
			select {
			case <-s.closeC:
				return
			default:
			}
			target := heart.GetHeartRate()
			isNewTarget := target != s.rate
			if isNewTarget {
				s.rate = target
				s.lastUpdate = time.Now()
			}

			var out int32
			switch {
			case s.EnableSmoothing:
				out = s.applySmoothing(target)
			case isNewTarget:
				out = target
			case s.EnableRandomOffset:
				out = s.randomOffset(target)
			default:
				continue
			}

			err := s.Send(out)
			if err != nil {
				panic(err)
			}
		}
	}()
	return nil
}

func (s *Sender) Close() error {
	close(s.closeC)
	return nil
}
