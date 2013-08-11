package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var port = flag.Int("port", 80, "port for web interface")

func intData(w http.ResponseWriter, req *http.Request, arg string, def int64) (int64, error) {
	valStr := req.FormValue(arg)
	if valStr == "" {
		return def, nil
	}
	val, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("could not parse %s: %v", arg, err),
			http.StatusBadRequest)
		return 0, err
	}
	return val, nil
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

		if len(s.History) == 0 {
			resp.WriteHeader(http.StatusNoContent)
			return
		}

		b, err := json.Marshal(s.History[len(s.History)-1])
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
		s.DataLock.Lock()
		defer s.DataLock.Unlock()

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
		s.SavePid()
		http.Redirect(resp, req, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/enable", func(w http.ResponseWriter, r *http.Request) {
		s.Enabled = !s.Enabled
		log.Printf("set enabled to %v", s.Enabled)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/csv", s.DumpCsv)
	http.HandleFunc("/plot", s.GenerateChart2)
	http.HandleFunc("/timer", AddTimerHandler)
	http.HandleFunc("/timers", GetTimersHandler)
	http.HandleFunc("/delete_timer", DeleteTimerHandler)
	http.Handle("/", http.FileServer(http.Dir("static/")))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
