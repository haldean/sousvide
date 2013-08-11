package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"
)

type Timer struct {
	Id        int64
	Name      string
	SetTime   time.Duration
	ExpiresAt time.Time
	TimeRemaining time.Duration
	Expired       bool
	Notified bool
}

type Timers []*Timer

var timers = make(Timers, 0)
var nextId = int64(0)

func StartTimerUpdateLoop() {
	for _ = range time.Tick(time.Second) {
		for _, t := range timers {
			t.TimeRemaining = time.Now().Sub(t.ExpiresAt)
			t.Expired = t.TimeRemaining < 0
		}
	}
}

func (t Timers) Len() int      { return len(t) }
func (t Timers) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

func (t Timers) Less(i, j int) bool {
	return t[i].TimeRemaining < t[j].TimeRemaining
}

func AddTimerHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "missing argument name", http.StatusBadRequest)
		return
	}
	h, err := intData(w, r, "h", 0)
	if err != nil {
		return
	}
	m, err := intData(w, r, "m", 0)
	if err != nil {
		return
	}
	s, err := intData(w, r, "s", 0)
	if err != nil {
		return
	}
	if (h == 0 && m == 0 && s == 0) || (h < 0 || m < 0 || s < 0) {
		http.Error(w, "must set timer for time in the future", http.StatusBadRequest)
		return
	}

	t := &Timer{
		Id:   nextId,
		Name: name,
		SetTime: time.Duration(h)*time.Hour +
			time.Duration(m)*time.Minute +
			time.Duration(s)*time.Second,
	}
	t.ExpiresAt = time.Now().Add(t.SetTime)

	nextId++
	log.Printf("set timer %v\n", t)
	timers = append(timers, t)
	sort.Sort(timers)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func GetTimersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	b, err := json.Marshal(timers)
	if err != nil {
		log.Panicf("could not marshal timer data to json: %v", err)
	}
	w.Write(b)
	for _, t := range timers {
		t.Notified = t.Expired
	}
}

func DeleteTimerHandler(w http.ResponseWriter, r *http.Request) {
	id, err := intData(w, r, "id", -1)
	if err != nil {
		return
	} else if id == -1 {
		http.Error(w, "must specify ID to delete", http.StatusBadRequest)
		return
	}
	idx := -1
	for i, t := range timers {
		if t.Id == id {
			idx = i
			break
		}
	}
	if idx == -1 {
		http.Error(
			w, fmt.Sprintf("could not find ID %d", id), http.StatusBadRequest)
		return
	}
	timers[idx] = timers[len(timers)-1]
	timers = timers[:len(timers)-1]
	sort.Sort(timers)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
