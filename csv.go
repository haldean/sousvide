package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (h HistorySample) ToCsv() string {
	return fmt.Sprintf("%d,%v,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f\n",
		h.Time.Unix(), h.Heating, h.Temp, h.Target, h.AbsError, h.Pid.P,
		h.Pid.I, h.Pid.D, h.POutput, h.IOutput, h.DOutput, h.COutput)

}

func (s *SousVide) DumpCsv(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Time,Heating,Temperature,Target,Error,\"P Coeff\",\"I " +
		"Coeff\",D Coeff\",\"P Term\",\"I Term\",\"D Term\",Controller\n"))
	for _, h := range s.History {
		w.Write([]byte(h.ToCsv()))
	}
}

func (s *SousVide) DumpJson(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-type", "application/json")

	if len(s.History) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	b, err := json.Marshal(s.History)
	if err != nil {
		log.Panicf("could not marshal historical data to json: %v", err)
	}
	w.Write(b)
}
