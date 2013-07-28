package main

import (
	"github.com/ajstarks/svgo"
	"github.com/haldean/chart"
	"github.com/haldean/chart/svgg"
	"image/color"
	"math"
	"net/http"
)

const (
	imgWidth  = 600
	imgHeight = 400
	useSvg = true
)

func (s *SousVide) GenerateChart(w http.ResponseWriter, req *http.Request) {
	c := chart.ScatterChart{}

	c.Key.Hide = true
	c.XRange.TicSetting.Hide = true
	c.YRange.MinMode.Fixed = true
	c.YRange.MinMode.Value = 0
	c.YRange.TicSetting.Grid = 1
	c.YRange.TicSetting.HideLabels = true

	s.HistoryLock.Lock()
	h := &s.History

	c.XRange.Fixed(0, float64(h.End) + 1, float64(h.End / 10))

	temps := make([]chart.EPoint, 0, h.End)
	targets := make([]chart.EPoint, 0, h.End)
	errs := make([]chart.EPoint, 0, h.End)
	var ep chart.EPoint
	for i, _ := range h.Times[:h.End] {
		ep = chart.EPoint{
			X:      float64(i),
			Y:      h.Temps[i],
			DeltaX: math.NaN(),
			DeltaY: math.NaN(),
		}
		temps = append(temps, ep)

		ep = chart.EPoint{
			X:      float64(i),
			Y:      h.Targets[i],
			DeltaX: math.NaN(),
			DeltaY: math.NaN(),
		}
		targets = append(targets, ep)

		ep = chart.EPoint{
			X:      float64(i),
			Y:      math.Abs(h.Temps[i] - h.Targets[i]),
			DeltaX: math.NaN(),
			DeltaY: math.NaN(),
		}
		errs = append(errs, ep)
	}
	s.HistoryLock.Unlock()

	c.AddData("Temperature", temps, chart.PlotStyleLines, chart.Style{
		LineColor: color.NRGBA{0xFF, 0x00, 0x00, 0xFF}, LineWidth: 2,
	})
	c.AddData("Target", targets, chart.PlotStyleLines, chart.Style{
		LineColor: color.NRGBA{0x00, 0x00, 0xFF, 0xFF}, LineWidth: 2,
	})
	c.AddData("Error", errs, chart.PlotStyleLines, chart.Style{
		LineColor: color.NRGBA{0x00, 0x00, 0x00, 0x66}, LineWidth: 2,
	})

	w.Header().Set("Content-type", "image/svg+xml")
	svgs := svg.New(w)
	svgs.Start(imgWidth, imgHeight)
	svgs.Title("Temperature history (\u00B0C)")
	canvas := svgg.New(
		svgs, imgWidth, imgHeight, "monospace", 12,
		color.RGBA{0xFF, 0xFF, 0xFF, 0xFF})

	canvas.Begin()
	c.Plot(canvas)
	canvas.End()
	svgs.End()
}
