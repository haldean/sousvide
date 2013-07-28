package main

import (
	"fmt"
	"github.com/ajstarks/svgo"
	"math"
	"net/http"
)

const (
	bottomMargin = 30
	height       = 400
	leftMargin   = 30
	rightMargin  = 10
	textSize     = 10
	topMargin    = 10
	width        = 600
)

func text(c *svg.SVG, x, y int, t string) {
	c.Gtransform(fmt.Sprintf("translate(0,%d)", height-textSize))
	c.Gtransform("scale(1,-1)")
	c.Text(x, height-(y+textSize/2+2), t,
		"font-family=\"monospace\"", "font-size=\"10px\"")
	c.Gend()
	c.Gend()
}

func (s *SousVide) GenerateGraph(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "image/svg+xml")
	c := svg.New(w)
	c.Start(width, height)
	c.Gtransform(fmt.Sprintf("translate(0,%d)", height))
	c.Gtransform("scale(1,-1)")

	c.Line(leftMargin, bottomMargin, leftMargin, height-topMargin,
		"stroke=\"black\"", "stroke-width=\"1\"")
	c.Line(leftMargin, bottomMargin, width-rightMargin, bottomMargin,
		"stroke=\"black\"", "stroke-width=\"1\"")

	s.HistoryLock.Lock()
	defer s.HistoryLock.Unlock()

	t := &s.History
	maxf := float64(0)
	for i := range t.Temps {
		if t.Temps[i] > maxf {
			maxf = t.Temps[i]
		}
		if t.Targets[i] > maxf {
			maxf = t.Targets[i]
		}
	}
	maxf += 5
	// round up to the nearest ten
	max := int(10 * math.Ceil(maxf/10))
	// keep maxf around for convenience
	maxf = float64(max)

	pxPerTemp := float64(height-topMargin-bottomMargin) / maxf

	// draw ticks on y axis
	for i := 0; i <= max; i += 10 {
		y := bottomMargin + int(pxPerTemp*float64(i))
		c.Line(leftMargin-4, y, leftMargin, y,
			"stroke=\"black\"", "stroke-width=\"1\"")
		text(c, 5, y, fmt.Sprintf("%d", i))
	}

	c.Gend()
	c.Gend()
	c.End()
}
