package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// No mutex needed for only just ONE user!
type server struct {
	timings  []*Timing
	duration time.Duration
}

type scoreParams struct {
	WithOob  bool
	HasMinus bool
	Polling  bool
	Value    string
}

func (s *server) InitRoutes(e *echo.Echo) {
	e.GET("/", s.handleIndex)
	e.GET("/score", s.handleGetScore)
	e.POST("/timing/consume", s.handleAddTime(CONSUME))
	e.POST("/timing/topup", s.handleAddTime(TOPUP))
	e.POST("/timing/stop", s.handleStopTime)
}

func (s *server) handleIndex(c echo.Context) error {
	score, hasMinus := s.score()
	polling := false
	last := s.lastTiming()
	templ := "controller-start"
	if last != nil && !last.IsStopped() {
		polling = true
		if last.TimeType == CONSUME {
			templ = "controller-consume"
		} else {
			templ = "controller-topup"
		}
	}
	return c.Render(http.StatusOK, "index", map[string]any{
		"Templ":   templ,
		"Timings": s.timings,
		"Score": scoreParams{
			WithOob:  false,
			HasMinus: hasMinus,
			Polling:  polling,
			Value:    score,
		},
	})
}

func (s *server) handleGetScore(c echo.Context) error {
	hasPolling := c.QueryParams().Has("polling")
	withOob := c.QueryParams().Has("withOob")
	score, hasMinus := s.score()
	return c.Render(http.StatusOK, "score", scoreParams{
		WithOob:  withOob,
		HasMinus: hasMinus,
		Polling:  hasPolling,
		Value:    score,
	})
}

func (s *server) handleAddTime(timingType TimingType) echo.HandlerFunc {
	return func(c echo.Context) error {
		timing := s.addTiming(timingType)
		score, hasMinus := s.score()
		return c.Render(http.StatusOK, fmt.Sprintf("controller-%s", timingType), map[string]any{
			"Timing": timing,
			"Score": scoreParams{
				WithOob:  true,
				HasMinus: hasMinus,
				Polling:  true,
				Value:    score,
			},
		})
	}
}

func (s *server) handleStopTime(c echo.Context) error {
	last := s.stopLastTiming()
	score, hasMinus := s.score()
	return c.Render(http.StatusOK, "controller-start", map[string]any{
		"Timing": last,
		"Score": scoreParams{
			WithOob:  true,
			HasMinus: hasMinus,
			Polling:  false,
			Value:    score,
		},
	})
}

func (s *server) score() (string, bool) {
	d := s.duration
	last := s.lastTiming()
	if last != nil && !last.IsStopped() {
		delta := time.Now().Sub(last.Start)
		if last.TimeType == CONSUME {
			d -= delta
		} else {
			d += delta
		}
	}
	d = d.Truncate(time.Second)
	if d == 0 {
		return "-", true
	}
	ss := d.String()
	return ss, strings.HasPrefix(ss, "-")
}

func (s *server) stopLastTiming() *Timing {
	last := s.lastTiming()
	if last != nil && last.Stop.IsZero() {
		last.Stop = time.Now()
		delta := last.Stop.Sub(last.Start)
		if last.TimeType == CONSUME {
			s.duration -= delta
		} else {
			s.duration += delta
		}
	}
	return last
}

func (s *server) addTiming(timingType TimingType) *Timing {
	s.stopLastTiming()
	timing := &Timing{
		Start:    time.Now(),
		TimeType: TimingType(timingType),
	}
	s.timings = append(s.timings, timing)
	return timing
}

func (s *server) lastTiming() *Timing {
	if len(s.timings) <= 0 {
		return nil
	}
	return s.timings[len(s.timings)-1]
}
