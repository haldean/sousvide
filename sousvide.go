package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"sync"
	"time"
)

const (
	InterruptDelay = 3 * time.Second
	LogFile        = "runlog.txt"
	HistoryLength  = 2048
	LowpassSamples = 2
	AccErrorWindow = 32
)

var StubGpio = flag.Bool("stub_gpio", false, "stub GPIO calls for testing")
var FakeTemp = flag.Bool("fake_temp", false, "use fake temperature values")

type SousVide struct {
	Heating     bool
	Enabled 	bool
	Temp        Celsius
	Target      Celsius
	History     []HistorySample
	Pid         PidParams
	Gpio        GpioParams
	DataLock    sync.Mutex
	AccError    float64
	lastPOutput float64
	lastIOutput float64
	lastDOutput float64
	lastControl float64
}

type HistorySample struct {
	Time     time.Time
	Enabled  bool
	Heating  bool
	Temp     float64
	Target   float64
	AbsError float64
	AccError float64
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

type GpioParams struct {
	ThermFd     *os.File
	ThermReader *bufio.Reader
	HeaterFd    *os.File
	Stub        bool
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
		Enabled:  s.Enabled,
		Heating:  s.Heating,
		Temp:     float64(s.Temp),
		Target:   float64(s.Target),
		AbsError: float64(s.Error()),
		AccError: s.AccError,
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

	s.AccError = 0
	N := len(s.History)
	l := float64(0)
	for i := N - 1; i >= N-AccErrorWindow-1 && i >= 0; i-- {
		s.AccError += math.Abs(s.History[i].AbsError)
		l++
	}
	s.AccError /= l
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
	flag.Parse()

	s := New()
	s.Pid.P = 10
	s.Pid.I = 0.1
	s.Pid.D = 10
	s.Gpio.Stub = *StubGpio

	err := s.InitGpio()
	if err != nil {
		fmt.Printf("could not initialize gpio: %v\n", err)
		return
	}

	err = s.InitTherm()
	if err != nil {
		fmt.Printf("could not initialize thermocouple: %v\n", err)
		return
	}

	go s.StartControlLoop()
	s.StartServer()
}
