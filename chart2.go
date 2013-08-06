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
	svgs.Start(imgWidth, imgHeight)
	svgs.Title("Temperature history (\u00B0C)")

	// draw "metadata": heating, etc.
	if N > 1 {
		for i, h := range s.History {
			if h.Heating {
				x0 := int(float64(i) * pxPerUnitX)
				svgs.Rect(x0, 0, int(math.Ceil(pxPerUnitX)), ImgHeight,
				"fill:#F0F0F0")
			}
		}
	}

	// draw grid before data so it's under everything
	for i := 0; i <= ImgHeight; i += int(10 * pxPerUnitY) {
		y := ImgHeight - i
		svgs.Line(0, y, ImgWidth, y, "stroke:#DDD; stroke-width:1")
	}

	// draw data
	if N > 1 {
		lastX := int(0)
		lastTempY := int(s.History[0].Temp * pxPerUnitY)
		lastTargetY := int(s.History[0].Target * pxPerUnitY)
		for i, h := range s.History[1:] {
			x := int(float64(i + 1) * pxPerUnitX)
			tempY := int(h.Temp * pxPerUnitY)
			targetY := int(h.Target * pxPerUnitY)
			svgs.Line(lastX, ImgHeight - lastTempY, x, ImgHeight - tempY,
				"stroke:#FF0000; stroke-width:2")
			svgs.Line(lastX, ImgHeight - lastTargetY, x, ImgHeight - targetY,
				"stroke:#0000FF; stroke-width:2")
			lastX = x
			lastTempY = tempY
			lastTargetY = targetY
		}
	}

	// draw axes last so they're on top of everything else
	svgs.Line(0, ImgHeight, ImgWidth, ImgHeight, "stroke:#000000; stroke-width:3")
	svgs.Line(0, 0, 0, ImgHeight, "stroke:#000000; stroke-width:3")
	svgs.Line(ImgWidth, 0, ImgWidth, ImgHeight, "stroke:#000000; stroke-width:3")

	svgs.End()
}
