package main

import (
	"github.com/ajstarks/svgo"
	"math"
	"net/http"
)

const (
	ImgWidth  = 600
	ImgHeight = 400
)

func (s *SousVide) GenerateChart2(w http.ResponseWriter, r *http.Request) {
	if len(s.History) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	N := len(s.History)

	maxVal := float64(0)
	for _, h := range s.History {
		if h.Temp > maxVal {
			maxVal = h.Temp
		}
		if h.Target > maxVal {
			maxVal = h.Target
		}
	}
	maxY := 10 * math.Ceil(maxVal / 10)
	pxPerUnitY := float64(ImgHeight) / maxY

	maxX := float64(N - 1)
	pxPerUnitX := float64(ImgWidth) / maxX

	w.Header().Set("Content-type", "image/svg+xml")
	svgs := svg.New(w)
	svgs.Start(ImgWidth, ImgHeight)
	svgs.Title("Temperature history (\u00B0C)")

	// draw "metadata": heating, etc.
	if N > 1 {
		for i, h := range s.History {
			if h.Heating {
				x0 := int(float64(i) * pxPerUnitX)
				svgs.Rect(x0, 0, int(math.Ceil(pxPerUnitX)), ImgHeight,
				"fill:#F7F7F7")
			}
		}
	}

	// draw grid before data so it's under everything
	even := true
	for i := 0; i <= ImgHeight; i += int(5 * pxPerUnitY) {
		y := ImgHeight - i
		if even {
			svgs.Line(0, y, ImgWidth, y, "stroke:#DDD; stroke-width:1")
		} else {
			svgs.Line(0, y, ImgWidth, y, "stroke:#EEE; stroke-width:1")
		}
		even = !even
	}

	// draw data
	if N > 1 {
		xs := make([]int, N)
		temps := make([]int, N)
		targets := make([]int, N)
		for i, h := range s.History {
			xs[i] = int(float64(i) * pxPerUnitX)
			temps[i] = ImgHeight - int(h.Temp * pxPerUnitY)
			targets[i] = ImgHeight - int(h.Target * pxPerUnitY)
		}
		svgs.Polyline(xs, temps, "stroke:#FF0000; stroke-width:1; fill:none")
		svgs.Polyline(xs, targets, "stroke:#0000FF; stroke-width:1; fill:none")
	}

	// draw axes last so they're on top of everything else
	svgs.Line(0, ImgHeight, ImgWidth, ImgHeight, "stroke:#000000; stroke-width:3")
	svgs.Line(0, 0, 0, ImgHeight, "stroke:#000000; stroke-width:3")
	svgs.Line(ImgWidth, 0, ImgWidth, ImgHeight, "stroke:#000000; stroke-width:3")

	svgs.End()
}
