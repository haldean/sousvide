package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type apiData struct {
	Temps     []float64
	Targets   []float64
	LogErrors []float64
	PidParams PidParams
}

func floatData(w http.ResponseWriter, req *http.Request, arg string) (float64, error) {
	valStr := req.FormValue(arg)
	if valStr == "" {
		http.Error(
			w, fmt.Sprintf("missing argument %s", arg), http.StatusBadRequest)
		return 0, errors.New("argument not supplied in request")
	}
	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("could not parse %s: %v", arg, err),
			http.StatusBadRequest)
		return 0, err
	}
	return val, nil
}

func (s *SousVide) StartServer() {
	http.HandleFunc("/api_data", func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("Content-type", "application/json")

		s.DataLock.Lock()
		defer s.DataLock.Unlock()

		h := &s.History
		a := apiData{
			h.Temps[:h.End],
			h.Targets[:h.End],
			h.LogErrors[:h.End],
			s.Pid,
		}

		b, err := json.Marshal(a)
		if err != nil {
			log.Panicf("could not marshal temp data to json: %v", err)
		}
		resp.Write(b)
	})

	http.HandleFunc("/target", func(resp http.ResponseWriter, req *http.Request) {
		s.DataLock.Lock()
		defer s.DataLock.Unlock()

		t, err := floatData(resp, req, "target")
		if err != nil {
			return
		}
		s.Target = Celsius(t)
		s.checkpoint()
		http.Redirect(resp, req, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/pid", func(resp http.ResponseWriter, req *http.Request) {
		log.Printf("acquire lock")
		s.DataLock.Lock()
		defer s.DataLock.Unlock()
		log.Printf("acquired lock")

		p, err := floatData(resp, req, "p")
		if err != nil {
			return
		}
		i, err := floatData(resp, req, "i")
		if err != nil {
			return
		}
		d, err := floatData(resp, req, "d")
		if err != nil {
			return
		}
		s.Pid.P = p
		s.Pid.I = i
		s.Pid.D = d
		s.checkpoint()
		http.Redirect(resp, req, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/plot", s.GenerateChart)

	http.Handle("/", http.FileServer(http.Dir("static/")))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
