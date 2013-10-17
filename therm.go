package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func findSerial() (string, error) {
	dir, err := os.Open("/sys/bus/w1/devices/")
	if err != nil {
		return "", err
	}
	defer dir.Close()

	devices, err := dir.Readdirnames(0)
	if err != nil {
		return "", err
	}
	for _, d := range devices {
		if !strings.Contains(d, "w1 bus master") {
			return d, nil
		}
	}
	return "", errors.New("no 1-wire devices. did you load the kernel modules?")
}

func (s *SousVide) InitTherm() error {
	var err error
	if *FakeTemp {
		s.Temp = s.Target
		return nil
	}

	if s.Gpio.Stub {
		s.Gpio.ThermFd, err = os.OpenFile(
			"test_temp.txt", os.O_RDONLY|os.O_SYNC, 0666)
		if err != nil {
			return err
		}
	} else {
		serial, err := findSerial()
		if err != nil {
			return err
		}

		s.Gpio.ThermFd, err = os.OpenFile(
			fmt.Sprintf("/sys/bus/w1/devices/%s/w1_slave", serial),
			os.O_RDONLY|os.O_SYNC, 0666)
		if err != nil {
			return err
		}
	}

	s.Gpio.ThermReader = bufio.NewReader(s.Gpio.ThermFd)
	return nil
}

func (s *SousVide) MeasureTemp() error {
	if *FakeTemp {
		if s.Heating {
			s.Temp += Celsius(rand.Float64())
		} else {
			s.Temp -= Celsius(rand.Float64())
		}
		if s.Temp < 0 {
			s.Temp = 0
		}
		return nil
	}

	s.Gpio.ThermFd.Seek(0, 0)
	line, err := s.Gpio.ThermReader.ReadString('\n')
	if err != nil {
		return err
	}
	line = strings.TrimSpace(line)
	data := strings.Split(strings.Split(line, "=")[1], " ")
	if len(data) != 2 {
		return errors.New(
			fmt.Sprintf("malformed line from 1-wire interface: %s", line))
	}
	if data[1] != "YES" {
		// read the next line to flush the buffer. when we get a NO status
		// there's still a line with the last temperature it read
		s.Gpio.ThermReader.ReadString('\n')
		return errors.New(fmt.Sprintf(
			"thermocouple did not return 'YES' status, got %s", line))
	}

	line, err = s.Gpio.ThermReader.ReadString('\n')
	if err != nil {
		return err
	}
	line = strings.TrimSpace(line)
	val := strings.Split(line, "=")[1]
	floatVal, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return err
	}

	s.Temp = Celsius(floatVal / 1000)
	log.Printf("read temperature %f", s.Temp)
	return nil
}
