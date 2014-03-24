package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sync"
	"time"
)

const (
	InterruptDelay = 1 * time.Second
	LogFile        = "runlog.txt"
	HistoryLength  = 8192
	LowpassSamples = 2
	AccErrorWindow = 32
)

var StubGpio = flag.Bool("stub_gpio", false, "stub GPIO calls for testing")
var FakeTemp = flag.Bool("fake_temp", false, "use fake temperature values")
var PidFile = flag.String("pid_file", "pid.json", "file to save PID values in")
var StartEnabled = flag.Bool("enabled", false, "start with heater enabled")
var StartTarget = flag.Float64("target", 0, "initial target temperature, in C")

var Stream chan HistorySample

type SousVide struct {
	Heating     bool
	Enabled     bool
	Temp        Celsius
	Target      Celsius
	History     []HistorySample
	Pid         PidParams
	Gpio        GpioParams
	DataLock    sync.Mutex
	AccError    float64
	MaxError    float64
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
	MaxError float64
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
		MaxError: s.MaxError,
		Pid:      s.Pid,
		POutput:  s.lastPOutput,
		IOutput:  s.lastIOutput,
		DOutput:  s.lastDOutput,
		COutput:  s.lastControl,
	}
}

func (s *SousVide) checkpoint() {
	snapshot := s.Snapshot()
	if len(s.History) == HistoryLength {
		for i := 0; i < HistoryLength-1; i++ {
			s.History[i] = s.History[i+1]
		}
		s.History[len(s.History)-1] = snapshot
	} else {
		s.History = append(s.History, snapshot)
	}
	Stream <- snapshot

	s.AccError = 0
	s.MaxError = 0
	N := len(s.History)
	l := float64(0)
	for i := N - 1; i >= N-AccErrorWindow-1 && i >= 0; i-- {
		ae := s.History[i].AbsError
		s.AccError += math.Abs(ae)
		if ae < s.MaxError {
			// find the most negative error
			s.MaxError = ae
		}
		l++
	}
	s.MaxError *= -1
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

func (s *SousVide) SavePid() {
	fd, err := os.Create(*PidFile)
	if err != nil {
		fmt.Printf("could not save PID values: %v\n", err)
		return
	}
	defer fd.Close()

	b, err := json.Marshal(s.Pid)
	if err != nil {
		fmt.Printf("could not save PID values: %v\n", err)
		return
	}
	fd.Write(b)
}

func (s *SousVide) LoadPid() error {
	bytes, err := ioutil.ReadFile(*PidFile)
	if err != nil {
		fmt.Printf("could not load PID values: %v\n", err)
		return err
	}
	err = json.Unmarshal(bytes, &s.Pid)
	if err != nil {
		fmt.Printf("could not load PID values: %v\n", err)
		return err
	}
	return nil
}

func main() {
	flag.Parse()

	s := New()
	err := s.LoadPid()
	if err != nil {
		s.Pid.P = 10
		s.Pid.D = 20
		s.SavePid()
	}
	s.Gpio.Stub = *StubGpio
	s.Target = Celsius(*StartTarget)
	s.Enabled = *StartEnabled

	err = s.InitGpio()
	if err != nil {
		fmt.Printf("could not initialize gpio: %v\n", err)
		return
	}

	err = s.InitTherm()
	if err != nil {
		fmt.Printf("could not initialize thermocouple: %v\n", err)
		return
	}

	Stream = StartSockServer()
	go StartTimerUpdateLoop()
	go s.StartControlLoop()
	go StartBroadcast()
	s.StartServer()
}
