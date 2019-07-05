package util

import (
	"github.com/gizak/termui/v3/widgets"
)

// generate a display of the drone's bearing relative to start position
func NewCompass(x1 int, y1 int, x2 int, y2 int) * widgets.PieChart {
	pc := widgets.NewPieChart()
	pc.Title = "Bearing"
	pc.SetRect(x1, y1, x2, y2)
	pc.Data = []float64{.1, .9 }

	return pc
}