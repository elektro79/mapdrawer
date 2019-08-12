package mapDrawer

import (
	"fmt"
	"image/color"

	"regexp"
	"strconv"

	"github.com/llgcode/draw2d"
)

type Drawer struct {
	width            int
	height           int
	drawerObjectList []DrawerObject
	margin           int
	maxZoom          uint8
	tileDrawer       DrawerTile
	conversor        Converter
}

type DrawerObject interface {
	Draw(*Converter)
	GetLocs() [][2]float64
}

type DrawerTile interface {
	Draw(*Converter)
}

func NewDrawer(gc draw2d.GraphicContext, w, h, margin int) *Drawer {
	d := &Drawer{
		width:            w,
		height:           h,
		maxZoom:          18,
		margin:           margin,
		tileDrawer:       ImageTileStreetMap(),
		drawerObjectList: make([]DrawerObject, 0),
	}
	d.conversor.Gc = gc
	return d
}

func (d *Drawer) GetSize() (int, int) {
	return d.width, d.height
}

func (d *Drawer) GetBounds() (float64, float64, float64, float64, uint8) {
	return d.conversor.MinLoc[0], d.conversor.MinLoc[1], d.conversor.MaxLoc[0], d.conversor.MinLoc[1], d.conversor.Zoom
}

func (d *Drawer) Prepend(do DrawerObject) {
	d.drawerObjectList = append([]DrawerObject{do}, d.drawerObjectList...)
}

func (d *Drawer) Add(do DrawerObject) {
	d.drawerObjectList = append(d.drawerObjectList, do)
}

func (d *Drawer) SetDrawAreaFromDraw() {
	d.conversor.setBounds(d.drawerObjectList)
}

func (d *Drawer) Draw() {
	if d.conversor.Width == 0 {
		d.SetDrawAreaFromDraw()
		d.conversor.set(d.width, d.height, d.maxZoom, d.margin)
	}
	d.tileDrawer.Draw(&d.conversor)
	d.conversor.draw(d.drawerObjectList)
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

func StringToColor(scol string) (color.Color, error) {
	format := "#%02x%02x%02x%02x"
	var r, g, b, a uint8
	n, err := fmt.Sscanf(scol, format, &r, &g, &b, &a)
	if err != nil {
		return nil, err
	}
	if n != 4 {
		return nil, fmt.Errorf("color: %v is not a hex-color alpha", scol)
	}
	return color.NRGBA{r, g, b, a}, nil
}
