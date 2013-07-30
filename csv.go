package main

import (
	"fmt"
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
