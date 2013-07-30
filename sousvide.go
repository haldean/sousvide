package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

const (
	InterruptDelay = 1 * time.Second
	LogFile        = "runlog.txt"
	HistoryLength  = 2048
	LowpassSamples = 3
)

type SousVide struct {
	Heating     bool
	Temp        Celsius
	Target      Celsius
	History     []HistorySample
	Pid         PidParams
	DataLock    sync.Mutex
	lastPOutput float64
	lastIOutput float64
	lastDOutput float64
	lastControl float64
}

type HistorySample struct {
	Time     time.Time
	Heating  bool
	Temp     float64
	Target   float64
	AbsError float64
	Pid      PidParams
	POutput  float64
	IOutput  float64
	DOutput  float64
	COutput  float64
}

type PidParams struct {
	P float64
	I float64
	D float64
}

type Celsius float64

func New() *SousVide {
	s := new(SousVide)
	s.History = make([]HistorySample, 0, HistoryLength)
	return s
}

func (s *SousVide) Snapshot() HistorySample {
	return HistorySample{
		Time:     time.Now(),
		Heating:  s.Heating,
		Temp:     float64(s.Temp),
		Target:   float64(s.Target),
		AbsError: float64(s.Error()),
		Pid:      s.Pid,
		POutput:  s.lastPOutput,
		IOutput:  s.lastIOutput,
		DOutput:  s.lastDOutput,
		COutput:  s.lastControl,
	}
}

func (s *SousVide) checkpoint() {
	if len(s.History) == HistoryLength {
		for i := 0; i < HistoryLength-1; i++ {
			s.History[i] = s.History[i+1]
		}
		s.History[len(s.History)-1] = s.Snapshot()
	} else {
		s.History = append(s.History, s.Snapshot())
	}
}

func (s *SousVide) StartControlLoop() {
	tick := time.Tick(InterruptDelay)
	for _ = range tick {
		s.DataLock.Lock()
		if s.Heating {
			s.Temp += Celsius(10 * rand.Float64())
		} else {
			s.Temp -= Celsius(10 * rand.Float64())
		}
		log.Printf("read temperature %f deg C", s.Temp)

		co := s.ControllerResult()
		s.Heating = co > 0
		log.Printf("controller returned %v", co)

		s.checkpoint()
		s.DataLock.Unlock()
	}
}

func (s *SousVide) ControllerResult() Celsius {
	s.lastPOutput = s.Pid.P * float64(s.Error())

	if len(s.History) > 0 {
		integral := float64(0)
		for _, h := range s.History {
			integral += h.AbsError
		}
		integral /= float64(len(s.History))
		s.lastIOutput = s.Pid.I * integral
	}

	// ignore derivative term if we have no history to use
	if len(s.History) > LowpassSamples {
		// use weighted window over three samples instead of two to act as a
		// low-pass filter
		N := len(s.History)
		d := (s.History[N-LowpassSamples-1].Temp - s.History[N-1].Temp) / 2
		s.lastDOutput = s.Pid.D * d
	}

	s.lastControl = s.lastPOutput + s.lastIOutput + s.lastDOutput
	return Celsius(s.lastControl)
}

func (s *SousVide) SetTarget(target Celsius) {
	s.DataLock.Lock()
	defer s.DataLock.Unlock()

	s.Target = target
	s.checkpoint()
}

func (s *SousVide) Error() Celsius {
	return s.Target - s.Temp
}

func main() {
	s := New()
	s.Target = 200
	s.Pid.P = 10
	s.Pid.I = 0.1
	s.Pid.D = 10

	go s.StartControlLoop()
	s.StartServer()
}
