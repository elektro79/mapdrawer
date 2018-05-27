package mapDrawer

import (
	"io"

	"regexp"
	"strconv"

	"github.com/llgcode/draw2d"
)

type Drawer struct {
	width            int
	height           int
	drawerObjectList []drawerObject
	margin           int
	maxZoom          uint8
	tileDrawer       DrawerTile
	conversor        Converter
}

type drawerObject interface {
	Draw(*Converter)
	GetLocs() [][2]float64
}

type DrawerTile interface {
	Draw(*Converter)
}

func NewDrawer() *Drawer {
	return &Drawer{
		width:            640,
		height:           640,
		maxZoom:          18,
		tileDrawer:       ImageTileStreetMap(),
		drawerObjectList: make([]drawerObject, 0),
	}
}

func (d *Drawer) Add(do drawerObject) {
	d.drawerObjectList = append(d.drawerObjectList, do)
}

func (d *Drawer) SetDrawAreaFromDraw() {
	d.conversor.setBounds(d.drawerObjectList)
}

func (d *Drawer) Draw(w io.Writer) {
	if d.conversor.Width == 0 {
		d.SetDrawAreaFromDraw()
		d.conversor.set(d.width, d.height, d.maxZoom, d.margin)
	}
	d.tileDrawer.Draw(&d.conversor)
	d.conversor.draw(d.drawerObjectList, w)
}

func DrawPath(gc draw2d.GraphicContext, d string, scale float64) {
	var lx, ly float64
	r := regexp.MustCompile(",| ")
	ld := r.Split(d, -1)
	pos := 0
	for pos < len(ld) {
		switch ld[pos] {
		case "M":
			lx = 0
			ly = 0
			gc.MoveTo(parseFloat(ld[pos+1], scale), parseFloat(ld[pos+2], scale))
			pos += 3
		case "m":
			lx = parseFloat(ld[pos+1], scale)
			ly = parseFloat(ld[pos+2], scale)
			gc.MoveTo(lx, ly)
			pos += 3
		case "C":
			for {
				lx = parseFloat(ld[pos+5], scale)
				ly = parseFloat(ld[pos+6], scale)
				gc.CubicCurveTo(
					//gc.CubicTo(
					parseFloat(ld[pos+1], scale),
					parseFloat(ld[pos+2], scale),
					parseFloat(ld[pos+3], scale),
					parseFloat(ld[pos+4], scale),
					lx, ly)
				pos += 7
				if pos >= len(ld) {
					break
				}
				if _, err := strconv.ParseFloat(ld[pos], 64); err != nil {
					break
				}
				pos--
			}
		case "c":
			for {
				tx := parseFloat(ld[pos+5], scale) + lx
				ty := parseFloat(ld[pos+6], scale) + ly
				gc.CubicCurveTo(
					parseFloat(ld[pos+1], scale)+lx,
					parseFloat(ld[pos+2], scale)+ly,
					parseFloat(ld[pos+3], scale)+lx,
					parseFloat(ld[pos+4], scale)+ly,
					tx, ty)
				lx = tx
				ly = ty
				pos += 7
				if pos >= len(ld) {
					break
				}
				if _, err := strconv.ParseFloat(ld[pos], 64); err != nil {
					break
				}
				pos--
			}
		case "Z":
			gc.Close()
			pos++
		case "z":
			gc.Close()
			pos++
		}
	}
}

func parseFloat(s string, scale float64) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return f * scale
}
