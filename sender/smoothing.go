package sender

import (
	"math"
	"math/rand"
	"time"
)

const smoothingAlpha = 0.3

func (s *Sender) applySmoothing(target int32) int32 {
	if !s.smoothInit {
		s.smoothed = float64(target)
		s.smoothInit = true
		return target
	}
	s.smoothed += (float64(target) - s.smoothed) * smoothingAlpha
	return int32(math.Round(s.smoothed))
}

const min = -2
const max = 2

func (s *Sender) randomOffset(rate int32) (filtered int32) {
	if s.lastUpdate.Before(time.Now().Add(-time.Second * 2)) {
		r := rand.Int31n(max-min+1) + min
		filtered = s.rate + r
		return
	}
	return rate
}
