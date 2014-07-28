// Package svg provides an API for generating Scalable Vector Graphics (SVG)
package svg

import (
	"fmt"
	"io"

	"encoding/xml"
	"strings"
)

// SVG defines the location of the generated SVG
type SVG struct {
	Writer io.Writer
}

const (
	svginit = `<?xml version="1.0" encoding="utf-8"?>
<!-- Generator:  github.com/psmithuk/svg  -->
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg
	version="1.1"
	id="Layer_1"
	xmlns="http://www.w3.org/2000/svg"
	xmlns:xlink="http://www.w3.org/1999/xlink"
	x="0px" y="0px"
	width="%dpx" height="%dpx"
 	viewBox="0 0 %d %d"
 	enable-background="new 0 0 %d %d"
 	xml:space="preserve">`

	vbfmt      = `viewBox="%f %f %f %f"`
	emptyclose = "/>\n"
)

// New is the SVG constructor, specifying the io.Writer where the generated SVG is written.
func New(w io.Writer) *SVG { return &SVG{w} }

func (svg *SVG) print(a ...interface{}) (n int, errno error) {
	return fmt.Fprint(svg.Writer, a...)
}

func (svg *SVG) println(a ...interface{}) (n int, error error) {
	return fmt.Fprintln(svg.Writer, a...)
}

func (svg *SVG) printf(format string, a ...interface{}) (n int, errno error) {
	return fmt.Fprintf(svg.Writer, format, a...)
}

// Structure, Metadata, Scripting, Transformation, and Links

// Start begins the SVG document with the width w and height h
// Other attributes may be optionally added, for example viewbox or additional namespaces
// Standard Reference: http://www.w3.org/TR/SVG11/struct.html#SVGElement
func (svg *SVG) Start(w int, h int) {
	svg.printf(svginit, w, h, w, h, w, h)
}

// End the SVG document
func (svg *SVG) End() { svg.println("</svg>") }

// Script defines a script with a specified type, (for example "application/javascript").
// if the first variadic argument is a link, use only the link reference.
// Otherwise, treat those arguments as the text of the script (marked up as CDATA).
// if no data is specified, just close the script element
func (svg *SVG) Script(scriptype string, data ...string) {
	svg.printf(`<script type="%s"`, scriptype)
	switch {
	case len(data) == 1 && islink(data[0]):
		svg.printf(" %s/>\n", href(data[0]))

	case len(data) > 0:
		svg.printf(">\n<![CDATA[\n")
		for _, v := range data {
			svg.println(v)
		}
		svg.printf("]]>\n</script>\n")

	default:
		svg.println(`/>`)
	}
}

// Gstyle begins a group, with the specified style.
// Standard Reference: http://www.w3.org/TR/SVG11/struct.html#GElement
func (svg *SVG) Gstyle(s string) { svg.println(group("style", s)) }

// Gtransform begins a group, with the specified transform
// Standard Reference: http://www.w3.org/TR/SVG11/coords.html#TransformAttribute
func (svg *SVG) Gtransform(s string) { svg.println(group("transform", s)) }

// Translate begins coordinate translation, end with Gend()
// Standard Reference: http://www.w3.org/TR/SVG11/coords.html#TransformAttribute
func (svg *SVG) Translate(x, y float64) { svg.Gtransform(translate(x, y)) }

// Scale scales the coordinate system by n, end with Gend()
// Standard Reference: http://www.w3.org/TR/SVG11/coords.html#TransformAttribute
func (svg *SVG) Scale(n float64) { svg.Gtransform(scale(n)) }

// ScaleXY scales the coordinate system by dx and dy, end with Gend()
// Standard Reference: http://www.w3.org/TR/SVG11/coords.html#TransformAttribute
func (svg *SVG) ScaleXY(dx, dy float64) { svg.Gtransform(scaleXY(dx, dy)) }

// SkewX skews the x coordinate system by angle a, end with Gend()
// Standard Reference: http://www.w3.org/TR/SVG11/coords.html#TransformAttribute
func (svg *SVG) SkewX(a float64) { svg.Gtransform(skewX(a)) }

// SkewY skews the y coordinate system by angle a, end with Gend()
// Standard Reference: http://www.w3.org/TR/SVG11/coords.html#TransformAttribute
func (svg *SVG) SkewY(a float64) { svg.Gtransform(skewY(a)) }

// SkewXY skews x and y coordinates by ax, ay respectively, end with Gend()
// Standard Reference: http://www.w3.org/TR/SVG11/coords.html#TransformAttribute
func (svg *SVG) SkewXY(ax, ay float64) { svg.Gtransform(skewX(ax) + " " + skewY(ay)) }

// Rotate rotates the coordinate system by r degrees, end with Gend()
// Standard Reference: http://www.w3.org/TR/SVG11/coords.html#TransformAttribute
func (svg *SVG) Rotate(r float64) { svg.Gtransform(rotate(r)) }

// TranslateRotate translates the coordinate system to (x,y), then rotates to r degrees, end with Gend()
func (svg *SVG) TranslateRotate(x, y float64, r float64) {
	svg.Gtransform(translate(x, y) + " " + rotate(r))
}

// RotateTranslate rotates the coordinate system r degrees, then translates to (x,y), end with Gend()
func (svg *SVG) RotateTranslate(x, y float64, r float64) {
	svg.Gtransform(rotate(r) + " " + translate(x, y))
}

// Group begins a group with arbitrary attributes
func (svg *SVG) Group(s ...string) { svg.printf("<g %s\n", endstyle(s, `>`)) }

// Gid begins a group, with the specified id
func (svg *SVG) Gid(s string) {
	svg.print(`<g id="`)
	xml.Escape(svg.Writer, []byte(s))
	svg.println(`">`)
}

// Gend ends a group (must be paired with Gsttyle, Gtransform, Gid).
func (svg *SVG) Gend() { svg.println(`</g>`) }

// ClipPath defines a clip path
func (svg *SVG) ClipPath(s ...string) { svg.printf(`<clipPath %s`, endstyle(s, `>`)) }

// ClipEnd ends a ClipPath
func (svg *SVG) ClipEnd() {
	svg.println(`</clipPath>`)
}

// Def begins a defintion block.
// Standard Reference: http://www.w3.org/TR/SVG11/struct.html#DefsElement
func (svg *SVG) Def() { svg.println(`<defs>`) }

// DefEnd ends a defintion block.
func (svg *SVG) DefEnd() { svg.println(`</defs>`) }

// Marker defines a marker
// Standard reference: http://www.w3.org/TR/SVG11/painting.html#MarkerElement
func (svg *SVG) Marker(id string, x, y, width, height float64, s ...string) {
	svg.printf(`<marker id="%s" refX="%f" refY="%f" markerWidth="%f" markerHeight="%f" %s`,
		id, x, y, width, height, endstyle(s, ">\n"))
}

// MarkEnd ends a marker
func (svg *SVG) MarkerEnd() { svg.println(`</marker>`) }

// Pattern defines a pattern with the specified dimensions.
// The putype can be either "user" or "obj", which sets the patternUnits
// attribute to be either userSpaceOnUse or objectBoundingBox
// Standard reference: http://www.w3.org/TR/SVG11/pservers.html#Patterns
func (svg *SVG) Pattern(id string, x, y, width, height float64, putype string, s ...string) {
	puattr := "userSpaceOnUse"
	if putype != "user" {
		puattr = "objectBoundingBox"
	}
	svg.printf(`<pattern id="%s" x="%f" y="%f" width="%f" height="%f" patternUnits="%s" %s`,
		id, x, y, width, height, puattr, endstyle(s, ">\n"))
}

// PatternEnd ends a marker
func (svg *SVG) PatternEnd() { svg.println(`</pattern>`) }

// Desc specified the text of the description tag.
// Standard Reference: http://www.w3.org/TR/SVG11/struct.html#DescElement
func (svg *SVG) Desc(s string) { svg.tt("desc", s) }

// Title specified the text of the title tag.
// Standard Reference: http://www.w3.org/TR/SVG11/struct.html#TitleElement
func (svg *SVG) Title(s string) { svg.tt("title", s) }

// Link begins a link named "name", with the specified title.
// Standard Reference: http://www.w3.org/TR/SVG11/linking.html#Links
func (svg *SVG) Link(href string, title string) {
	svg.printf("<a xlink:href=\"%s\" xlink:title=\"", href)
	xml.Escape(svg.Writer, []byte(title))
	svg.println("\">")
}

// LinkEnd ends a link.
func (svg *SVG) LinkEnd() { svg.println(`</a>`) }

// Use places the object referenced at link at the location x, y, with optional style.
// Standard Reference: http://www.w3.org/TR/SVG11/struct.html#UseElement
func (svg *SVG) Use(x float64, y float64, link string, s ...string) {
	svg.printf(`<use %s %s %s`, loc(x, y), href(link), endstyle(s, emptyclose))
}

// Mask creates a mask with a specified id, dimension, and optional style.
func (svg *SVG) Mask(id string, x float64, y float64, w float64, h float64, s ...string) {
	svg.printf(`<mask id="%s" x="%f" y="%f" width="%f" height="%f" %s`, id, x, y, w, h, endstyle(s, `>`))
}

// MaskEnd ends a Mask.
func (svg *SVG) MaskEnd() { svg.println(`</mask>`) }

// Shapes

// Circle centered at x,y, with radius r, with optional style.
// Standard Reference: http://www.w3.org/TR/SVG11/shapes.html#CircleElement
func (svg *SVG) Circle(x float64, y float64, r float64, s ...string) {
	svg.printf(`<circle cx="%f" cy="%f" r="%f" %s`, x, y, r, endstyle(s, emptyclose))
}

// Ellipse centered at x,y, centered at x,y with radii w, and h, with optional style.
// Standard Reference: http://www.w3.org/TR/SVG11/shapes.html#EllipseElement
func (svg *SVG) Ellipse(x float64, y float64, w float64, h float64, s ...string) {
	svg.printf(`<ellipse cx="%f" cy="%f" rx="%f" ry="%f" %s`,
		x, y, w, h, endstyle(s, emptyclose))
}

// Polygon draws a series of line segments using an array of x, y coordinates, with optional style.
// Standard Reference: http://www.w3.org/TR/SVG11/shapes.html#PolygonElement
func (svg *SVG) Polygon(x []float64, y []float64, s ...string) {
	svg.poly(x, y, "polygon", s...)
}

// Rect draws a rectangle with upper left-hand corner at x,y, with width w, and height h, with optional style
// Standard Reference: http://www.w3.org/TR/SVG11/shapes.html#RectElement
func (svg *SVG) Rect(x float64, y float64, w float64, h float64, s ...string) {
	svg.printf(`<rect %s %s`, dim(x, y, w, h), endstyle(s, emptyclose))
}

// CenterRect draws a rectangle with its center at x,y, with width w, and height h, with optional style
func (svg *SVG) CenterRect(x float64, y float64, w float64, h float64, s ...string) {
	svg.Rect(x-(w/2), y-(h/2), w, h, s...)
}

// Roundrect draws a rounded rectangle with upper the left-hand corner at x,y,
// with width w, and height h. The radii for the rounded portion
// are specified by rx (width), and ry (height).
// Style is optional.
// Standard Reference: http://www.w3.org/TR/SVG11/shapes.html#RectElement
func (svg *SVG) Roundrect(x float64, y float64, w float64, h float64, rx float64, ry float64, s ...string) {
	svg.printf(`<rect %s rx="%f" ry="%f" %s`, dim(x, y, w, h), rx, ry, endstyle(s, emptyclose))
}

// Square draws a square with upper left corner at x,y with sides of length l, with optional style.
func (svg *SVG) Square(x float64, y float64, l float64, s ...string) {
	svg.Rect(x, y, l, l, s...)
}

// Paths

// Path draws an arbitrary path, the caller is responsible for structuring the path data
func (svg *SVG) Path(d string, s ...string) {
	svg.printf(`<path d="%s" %s`, d, endstyle(s, emptyclose))
}

// Arc draws an elliptical arc, with optional style, beginning coordinate at sx,sy, ending coordinate at ex, ey
// width and height of the arc are specified by ax, ay, the x axis rotation is r
// if sweep is true, then the arc will be drawn in a "positive-angle" direction (clockwise), if false,
// the arc is drawn counterclockwise.
// if large is true, the arc sweep angle is greater than or equal to 180 degrees,
// otherwise the arc sweep is less than 180 degrees
// http://www.w3.org/TR/SVG11/paths.html#PathDataEllipticalArcCommands
func (svg *SVG) Arc(sx float64, sy float64, ax float64, ay float64, r float64, large bool, sweep bool, ex float64, ey float64, s ...string) {
	svg.printf(`%s A%s %f %s %s %s" %s`,
		ptag(sx, sy), coord(ax, ay), r, onezero(large), onezero(sweep), coord(ex, ey), endstyle(s, emptyclose))
}

// Bezier draws a cubic bezier curve, with optional style, beginning at sx,sy, ending at ex,ey
// with control points at cx,cy and px,py.
// Standard Reference: http://www.w3.org/TR/SVG11/paths.html#PathDataCubicBezierCommands
func (svg *SVG) Bezier(sx float64, sy float64, cx float64, cy float64, px float64, py float64, ex float64, ey float64, s ...string) {
	svg.printf(`%s C%s %s %s" %s`,
		ptag(sx, sy), coord(cx, cy), coord(px, py), coord(ex, ey), endstyle(s, emptyclose))
}

// Qbez draws a quadratic bezier curver, with optional style
// beginning at sx,sy, ending at ex, sy with control points at cx, cy
// Standard Reference: http://www.w3.org/TR/SVG11/paths.html#PathDataQuadraticBezierCommands
func (svg *SVG) Qbez(sx float64, sy float64, cx float64, cy float64, ex float64, ey float64, s ...string) {
	svg.printf(`%s Q%s %s" %s`,
		ptag(sx, sy), coord(cx, cy), coord(ex, ey), endstyle(s, emptyclose))
}

// Qbezier draws a Quadratic Bezier curve, with optional style, beginning at sx, sy, ending at tx,ty
// with control points are at cx,cy, ex,ey.
// Standard Reference: http://www.w3.org/TR/SVG11/paths.html#PathDataQuadraticBezierCommands
func (svg *SVG) Qbezier(sx float64, sy float64, cx float64, cy float64, ex float64, ey float64, tx float64, ty float64, s ...string) {
	svg.printf(`%s Q%s %s T%s" %s`,
		ptag(sx, sy), coord(cx, cy), coord(ex, ey), coord(tx, ty), endstyle(s, emptyclose))
}

// Lines

// Line draws a straight line between two points, with optional style.
// Standard Reference: http://www.w3.org/TR/SVG11/shapes.html#LineElement
func (svg *SVG) Line(x1 float64, y1 float64, x2 float64, y2 float64, s ...string) {
	svg.printf(`<line x1="%f" y1="%f" x2="%f" y2="%f" %s`, x1, y1, x2, y2, endstyle(s, emptyclose))
}

// Polyline draws connected lines between coordinates, with optional style.
// Standard Reference: http://www.w3.org/TR/SVG11/shapes.html#PolylineElement
func (svg *SVG) Polyline(x []float64, y []float64, s ...string) {
	svg.poly(x, y, "polyline", s...)
}

// Image places at x,y (upper left hand corner), the image with
// width w, and height h, referenced at link, with optional style.
// Standard Reference: http://www.w3.org/TR/SVG11/struct.html#ImageElement
func (svg *SVG) Image(x float64, y float64, w float64, h float64, link string, s ...string) {
	svg.printf(`<image %s %s %s`, dim(x, y, w, h), href(link), endstyle(s, emptyclose))
}

// Text places the specified text, t at x,y according to the style specified in s
// Standard Reference: http://www.w3.org/TR/SVG11/text.html#TextElement
func (svg *SVG) Text(x float64, y float64, t string, s ...string) {
	svg.printf(`<text %s %s`, loc(x, y), endstyle(s, ">"))
	xml.Escape(svg.Writer, []byte(t))
	svg.println(`</text>`)
}

// Textpath places text optionally styled text along a previously defined path
// Standard Reference: http://www.w3.org/TR/SVG11/text.html#TextPathElement
func (svg *SVG) Textpath(t string, pathid string, s ...string) {
	svg.printf("<text %s<textPath xlink:href=\"%s\">", endstyle(s, ">"), pathid)
	xml.Escape(svg.Writer, []byte(t))
	svg.println(`</textPath></text>`)
}

// Textlines places a series of lines of text starting at x,y, at the specified size, fill, and alignment.
// Each line is spaced according to the spacing argument
func (svg *SVG) Textlines(x, y float64, s []string, size, spacing float64, fill, align string, units string) {
	svg.Gstyle(fmt.Sprintf("font-size:%f%s;fill:%s;text-anchor:%s", size, units, fill, align))
	for _, t := range s {
		svg.Text(x, y, t)
		y += spacing
	}
	svg.Gend()
}

// Colors

// RGB specifies a fill color in terms of a (r)ed, (g)reen, (b)lue triple.
// Standard reference: http://www.w3.org/TR/css3-color/
func (svg *SVG) RGB(r int, g int, b int) string {
	return fmt.Sprintf(`fill:rgb(%d,%d,%d)`, r, g, b)
}

// RGBA specifies a fill color in terms of a (r)ed, (g)reen, (b)lue triple and opacity.
func (svg *SVG) RGBA(r int, g int, b int, a float64) string {
	return fmt.Sprintf(`fill-opacity:%.2f; %s`, a, svg.RGB(r, g, b))
}

// Utility

// Grid draws a grid at the specified coordinate, dimensions, and spacing, with optional style.
func (svg *SVG) Grid(x float64, y float64, w float64, h float64, n float64, s ...string) {

	if len(s) > 0 {
		svg.Gstyle(s[0])
	}
	for ix := x; ix <= x+w; ix += n {
		svg.Line(ix, y, ix, y+h)
	}

	for iy := y; iy <= y+h; iy += n {
		svg.Line(x, iy, x+w, iy)
	}
	if len(s) > 0 {
		svg.Gend()
	}

}

// Support functions

// style returns a style name,attribute string
func style(s string) string {
	if len(s) > 0 {
		return fmt.Sprintf(`style="%s"`, s)
	}
	return s
}

// pp returns a series of polygon points
func (svg *SVG) pp(x []float64, y []float64, tag string) {
	svg.print(tag)
	if len(x) != len(y) {
		svg.print(" ")
		return
	}
	lx := len(x) - 1
	for i := 0; i < lx; i++ {
		svg.print(coord(x[i], y[i]) + " ")
	}
	svg.print(coord(x[lx], y[lx]))
}

// endstyle modifies an SVG object, with either a series of name="value" pairs,
// or a single string containing a style
func endstyle(s []string, endtag string) string {
	if len(s) > 0 {
		nv := ""
		for i := 0; i < len(s); i++ {
			if strings.Index(s[i], "=") > 0 {
				nv += (s[i]) + " "
			} else {
				nv += style(s[i])
			}
		}
		return nv + endtag
	}
	return endtag

}

// tt creates a xml element, tag containing s
func (svg *SVG) tt(tag string, s string) {
	svg.print("<" + tag + ">")
	xml.Escape(svg.Writer, []byte(s))
	svg.println("</" + tag + ">")
}

// poly compiles the polygon element
func (svg *SVG) poly(x []float64, y []float64, tag string, s ...string) {
	svg.pp(x, y, "<"+tag+" points=\"")
	svg.print(`" ` + endstyle(s, "/>\n"))
}

// onezero returns "0" or "1"
func onezero(flag bool) string {
	if flag {
		return "1"
	}
	return "0"
}

// pct returns a percetage, capped at 100
func pct(n uint8) uint8 {
	if n > 100 {
		return 100
	}
	return n
}

// islink determines if a string is a script reference
func islink(link string) bool {
	return strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "#") ||
		strings.HasPrefix(link, "../") || strings.HasPrefix(link, "./")
}

// group returns a group element
func group(tag string, value string) string { return fmt.Sprintf(`<g %s="%s">`, tag, value) }

// scale return the scale string for the transform
func scale(n float64) string { return fmt.Sprintf(`scale(%g)`, n) }

// scaleXY return the scale string for the transform
func scaleXY(dx, dy float64) string { return fmt.Sprintf(`scale(%g,%g)`, dx, dy) }

// skewx returns the skewX string for the transform
func skewX(angle float64) string { return fmt.Sprintf(`skewX(%g)`, angle) }

// skewx returns the skewX string for the transform
func skewY(angle float64) string { return fmt.Sprintf(`skewY(%g)`, angle) }

// rotate returns the rotate string for the transform
func rotate(r float64) string { return fmt.Sprintf(`rotate(%g)`, r) }

// translate returns the translate string for the transform
func translate(x, y float64) string { return fmt.Sprintf(`translate(%f,%f)`, x, y) }

// coord returns a coordinate string
func coord(x float64, y float64) string { return fmt.Sprintf(`%f,%f`, x, y) }

// ptag returns the beginning of the path element
func ptag(x float64, y float64) string { return fmt.Sprintf(`<path d="M%s`, coord(x, y)) }

// loc returns the x and y coordinate attributes
func loc(x float64, y float64) string { return fmt.Sprintf(`x="%f" y="%f"`, x, y) }

// href returns the href name and attribute
func href(s string) string { return fmt.Sprintf(`xlink:href="%s"`, s) }

// dim returns the dimension string (x, y coordinates and width, height)
func dim(x float64, y float64, w float64, h float64) string {
	return fmt.Sprintf(`x="%f" y="%f" width="%f" height="%f"`, x, y, w, h)
}

// // tablevalues outputs a series of values as a XML attribute
// func (svg *SVG) tablevalues(s string, t []float64) {
// 	svg.printf(` %s="`, s)
// 	for i := 0; i < len(t)-1; i++ {
// 		svg.printf("%g ", t[i])
// 	}
// 	svg.printf(`%g"%s`, t[len(t)-1], emptyclose)
// }

// imgchannel validates the image channel indicator
func imgchannel(c string) string {
	switch c {
	case "R", "G", "B", "A":
		return c
	case "r", "g", "b", "a":
		return strings.ToUpper(c)
	case "red", "green", "blue", "alpha":
		return strings.ToUpper(c[0:1])
	case "Red", "Green", "Blue", "Alpha":
		return c[0:1]
	}
	return "R"
}
