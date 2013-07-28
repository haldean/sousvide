package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type apiData struct {
	Temps     []float64
	Targets   []float64
	LogErrors []float64
}

func (s *SousVide) StartServer() {
	http.HandleFunc("/api_data", func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("Content-type", "application/json")

		s.HistoryLock.Lock()
		defer s.HistoryLock.Unlock()

		h := &s.History
		a := apiData{
			h.Temps[:h.End],
			h.Targets[:h.End],
			h.LogErrors[:h.End],
		}

		b, err := json.Marshal(a)
		if err != nil {
			log.Panicf("could not marshal temp data to json: %v", err)
		}
		resp.Write(b)
	})

	http.HandleFunc("/target", func(resp http.ResponseWriter, req *http.Request) {
		t_str := req.FormValue("target")
		if t_str == "" {
			http.Error(resp, "no target specified", http.StatusBadRequest)
			return
		}
		target, err := strconv.ParseFloat(t_str, 64)
		if err != nil {
			http.Error(
				resp, fmt.Sprintf("could not parse target temp: %v", err),
				http.StatusBadRequest)
			return
		}
		s.Target = Celsius(target)
		s.checkpoint()
		http.Redirect(resp, req, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/plot", s.GenerateChart)

	http.Handle("/", http.FileServer(http.Dir("static/")))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
