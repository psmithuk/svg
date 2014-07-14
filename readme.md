# SVG

This package is a fork of [github.com/ajstarks/svgo](github.com/ajstarks/svgo) which has been made slightly more suitable for print graphics.

The key difference is that all integer positioning and dimensions have been changed to floating point. This improves accuracy when the SVG file isn't output to screen but manages to break almost every single API method in the process.

My personal use case was in the creation of visualisations which need to be imported to Adobe Illustrator for pre-press.


## API and Documentation

All comments and documentation from `ajstarks/svgo` have been retained. Most of the filters have been removed, whilst several new helper methods for typography have been introduced.

See [https://godoc.org/github.com/](https://godoc.org/github.com/psmithuk/svg) for more details.

To keep things working across web and print, use a mapping of 1.2 pixels to a point. You can pass a pixel width `px` to the canvas start method, see example.


## Example

```go
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
```