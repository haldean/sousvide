package main

import (
	"log"
	"math"
	"math/rand"
	"sync"
	"time"
)

const (
	InterruptDelay = 1 * time.Second
	LogFile        = "runlog.txt"
	HistoryLength  = 2048
)

type SousVide struct {
	Temp        Celsius
	Target      Celsius
	History     TempHistory
	HistoryLock sync.Mutex
}

// A ring buffer to store historical data for plotting
type TempHistory struct {
	Times     [HistoryLength]time.Time
	Temps     [HistoryLength]float64
	Targets   [HistoryLength]float64
	LogErrors [HistoryLength]float64
	End       int
}

type Celsius float64

func (s *SousVide) checkpoint() {
	s.HistoryLock.Lock()
	defer s.HistoryLock.Unlock()

	// this would be better implemented by a ring buffer, but it doesn't
	// actually buy me anything because on every change I have to write it to a
	// flat array to plot it anyway.

	h := &s.History
	if h.End == HistoryLength {
		for i := 0; i < HistoryLength-1; i++ {
			h.Times[i] = h.Times[i+1]
			h.Temps[i] = h.Temps[i+1]
			h.Targets[i] = h.Targets[i+1]
			h.LogErrors[i] = h.LogErrors[i+1]
		}
		h.End -= 1
	}

	h.Times[h.End] = time.Now()
	h.Temps[h.End] = float64(s.Temp)
	h.Targets[h.End] = float64(s.Target)
	h.LogErrors[h.End] = math.Abs(math.Log10(math.Abs(float64(s.Error()))))
	h.End += 1
}

func (s *SousVide) StartControlLoop() {
	tick := time.Tick(InterruptDelay)
	for _ = range tick {
		s.Temp -= 0.1*s.Error() + Celsius(rand.Float64()-0.5)
		log.Printf("read temperature %f deg C", s.Temp)
		s.checkpoint()
	}
}

func (s *SousVide) SetTarget(target Celsius) {
	s.Target = target
	s.checkpoint()
}

func (s *SousVide) Error() Celsius {
	return s.Temp - s.Target
}

func main() {
	s := new(SousVide)
	s.Target = 200
	go s.StartControlLoop()
	s.StartServer()
}
