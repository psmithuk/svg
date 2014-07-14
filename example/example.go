package main

import (
	"github.com/psmithuk/svg"
	"os"
)

var (
	canvas = svg.New(os.Stdout)

	width  = 500
	height = 500

	centreX = float64(width) / 2.0
	centreY = float64(height) / 2.0
)

func main() {
	canvas.Start(width, height, "px")
	canvas.Circle(centreX, centreY, float64(width)/3.0)
	canvas.Text(centreX, centreY, "Hello, SVG", "text-anchor:middle;font-size:10px;fill:white")
	canvas.End()
}
