// hat tip to github.com/aqua/raspberrypi for inspiration for the GPIO code

package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

const HeaterGpioPin = 18

func checkHeaterExported() error {
	_, err := os.Stat(fmt.Sprintf("/sys/class/gpio/gpio%d", HeaterGpioPin))
	if err == nil {
		return nil
	}
	if !os.IsNotExist(err) {
		return err
	}

	fd, err := os.OpenFile(
		"/sys/class/gpio/export", os.O_WRONLY|os.O_SYNC, 0666)
	if err != nil {
		return err
	}
	defer fd.Close()

	_, err = fmt.Fprintf(fd, "%d\n", HeaterGpioPin)
	return err
}

func setHeaterOutputMode() error {
	fd, err := os.OpenFile(
		fmt.Sprintf("/sys/class/gpio/gpio%d/direction", HeaterGpioPin),
		os.O_WRONLY|os.O_SYNC, 0666)
	if err != nil {
		return err
	}
	defer fd.Close()

	fmt.Fprintf(fd, "out")
	return nil
}

func (s *SousVide) InitGpio() error {
	if s.Gpio.Stub {
		s.Gpio.HeaterFd = os.Stdout
		return nil
	}

	err := checkHeaterExported()
	if err != nil {
		return err
	}
	err = setHeaterOutputMode()
	if err != nil {
		return err
	}

	s.Gpio.HeaterFd, err = os.OpenFile(
		fmt.Sprintf("/sys/class/gpio/gpio%d/value", HeaterGpioPin),
		os.O_WRONLY|os.O_SYNC, 0666)
	if err != nil {
		return err
	}

	return nil
}

func (s *SousVide) StartControlLoop() {
	tick := time.Tick(InterruptDelay)
	for _ = range tick {
		s.DataLock.Lock()
		err := s.MeasureTemp()
		if err != nil {
			log.Printf("could not read temperature: %v", err)
		} else {
			co := s.ControllerResult()
			s.Heating = co > 0
			s.UpdateHardware()
			s.checkpoint()
		}
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

func (s *SousVide) UpdateHardware() {
	var heatVal string
	if s.Heating {
		heatVal = "1\n"
	} else {
		heatVal = "0\n"
	}
	fmt.Fprintf(s.Gpio.HeaterFd, heatVal)
}
