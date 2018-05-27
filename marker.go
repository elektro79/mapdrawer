package mapDrawer

import (
	"fmt"
	"image/color"
	"math"

	"github.com/llgcode/draw2d"
)

type Size uint8

const (
	normal Size = iota
	mid
	small
	tiny
)

var scaleName = []float64{0.05625, 0.046875, 0.034375, 0.0203125}

type Marker struct {
	Size  Size
	Color color.Color
	Label string
	Icon  string
	Locs  [][2]float64
}

func (m *Marker) Draw(c *Converter) {
	gc := c.Gc
	for _, loc := range m.Locs {
		gc.Save()
		gc.SetFontData(draw2d.FontData{Name: "luxi", Family: draw2d.FontFamilyMono, Style: draw2d.FontStyleBold | draw2d.FontStyleItalic})
		gc.SetLineJoin(draw2d.RoundJoin)
		x, y := c.GetXY(loc)
		gc.Translate(x, y)
		scale := float64(c.Height) * scaleName[m.Size]
		gc.SetLineWidth(1.0)
		gc.SetFillColor(m.Color)
		gc.SetStrokeColor(color.Black)
		gc.BeginPath()
		DrawPath(gc, "M 0,0 c -0.02803654,-0.1546784 -0.1030343,-0.2943913 -0.18714392,-0.4232202 -0.0380496,-0.0659319 -0.0841096,-0.13092197 -0.0991292,-0.2077379 -0.0240313,-0.12600324 0.0430562,-0.26152998 0.15360021,-0.31668256 0.12806693,-0.0680249 0.30209377,-0.0219773 0.37729178,0.10779347 0.0570744,0.0910489 0.0580757,0.21516832 0,0.30768233 C 0.19154977,-0.42301088 0.1147496,-0.32756657 0.06568565,-0.2156916 0.03364388,-0.1476666 0.01061386,-0.07472286 -0.00340441,0 Z", scale)
		gc.FillStroke()
		if m.Label != "" {
			gc.BeginPath()
			gc.SetFillColor(color.Black)
			gc.SetFontSize(393.846153846 * scaleName[m.Size])
			//gc.SetFontSize(8)
			left, top, right, botton := gc.GetStringBounds(m.Label[0:1])
			fmt.Println(left, right, top, botton)
			gc.FillStringAt(m.Label[0:1], -(right / 2), (-0.71672547*scale)-(top/2))

		} else {
			gc.BeginPath()
			gc.SetFillColor(color.Black)
			gc.ArcTo(0.0, -0.71672547*scale, 0.10*scale, 0.10*scale, 0.0, 2*math.Pi)
			gc.FillStroke()
		}
		gc.Restore()
	}

	/*for _, loc := range m.Locs {
		gc.Save()
		x, y := c.GetXY(loc)
		fmt.Println(x, y, loc[0], loc[1])
		gc.Translate(x, y)
		scale := float64(c.height) * scaleName[m.Size]
		gc.Scale(scale, scale)
		drawPath(gc, "M 0,0 c -0.02803654,-0.1546784 -0.1030343,-0.2943913 -0.18714392,-0.4232202 -0.0380496,-0.0659319 -0.0841096,-0.13092197 -0.0991292,-0.2077379 -0.0240313,-0.12600324 0.0430562,-0.26152998 0.15360021,-0.31668256 0.12806693,-0.0680249 0.30209377,-0.0219773 0.37729178,0.10779347 0.0570744,0.0910489 0.0580757,0.21516832 0,0.30768233 C 0.19154977,-0.42301088 0.1147496,-0.32756657 0.06568565,-0.2156916 0.03364388,-0.1476666 0.01061386,-0.07472286 -0.00340441,0 Z")
		//gc.SetColor(m.Color)
		gc.SetFillColor(m.Color)
		gc.FillPreserve()
		gc.Fill()
		if m.Label != "" {
			gc.Restore()
			gc.Save()
			gc.Translate(x, y)
			gc.SetRGB(0, 0, 0)
			gc.Stroke()
			gc.FillPreserve()
			gc.DrawStringAnchored(m.Label[0:1], 0, -0.71672547*scale, 0.5, 0.5)
		} else {
			gc.FillPreserve()
			gc.Stroke()
			gc.DrawCircle(0, -0.71672547, 0.10)
		}
		gc.Restore()
	}*/
}

func (m *Marker) GetLocs() [][2]float64 {
	return m.Locs
}
