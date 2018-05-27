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
	conversor        Converter
}

type drawerObject interface {
	Draw(*Converter)
	GetLocs() [][2]float64
}

func NewDrawer() *Drawer {
	return &Drawer{
		width:            640,
		height:           640,
		maxZoom:          18,
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
	d.conversor.draw(d.drawerObjectList, w)
}

/*type Conversor struct {
	Width      int
	Height     int
	TilePixelX int
	TilePixelY int
	minLoc     [2]float64
	maxLoc     [2]float64
	zoom       uint8
	ntile      int
	Img        draw.Image
	Gc         draw2d.GraphicContext
}

func (c *Conversor) GetXY(ll [2]float64) (x, y float64) {
	n := float64(uint(1)<<uint(c.zoom)) * 256.0
	latr := ll[0] * math.Pi / 180
	x = mod((ll[1]+180.0)/360.0*n-float64(c.TilePixelX), n)
	y = (1.0-math.Log(math.Tan(latr)+(1.0/math.Cos(latr)))/math.Pi)/2.0*n - float64(c.TilePixelY)
	return
}

func (c *Conversor) setBounds(do []drawerObject) {
	minLoc := [2]float64{0.0, 0.0}
	maxLoc := [2]float64{0.0, 0.0}
	first := true
	for _, d := range do {
		for _, loc := range d.GetLocs() {
			if first {
				minLoc = loc
				maxLoc = loc
				first = false
			} else {
				minLoc[0] = math.Max(minLoc[0], loc[0])
				maxLoc[0] = math.Min(maxLoc[0], loc[0])
				if !dj(loc[1], minLoc[1], maxLoc[1]) {
					if bj(minLoc[1], loc[1]) < cj(maxLoc[1], loc[1]) {
						minLoc[1] = loc[1]
					} else {
						maxLoc[1] = loc[1]
					}
				}
			}
		}
	}
	if minLoc[1] > maxLoc[1] {
		maxLoc[1] += 360.0
	}
	c.minLoc = minLoc
	c.maxLoc = maxLoc
}
func bj(pd1, pd2 float64) float64 {
	return mod(360.0+(pd1-pd2), 360.0)
}

func cj(pd1, pd2 float64) float64 {
	return mod(360.0+(pd2-pd1), 360.0)
}

func dj(lng, mi, ma float64) bool {
	if mi <= ma {
		return (mi <= lng) && (lng <= ma)
	}
	return mi <= lng || lng <= ma
}
func (c *Conversor) set(w int, h int, maxZoom uint8, margin int) {
	c.Width = w
	c.Height = h
	c.Img = image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
	c.Gc = draw2dimg.NewGraphicContext(c.Img)
	minX := (c.minLoc[1] + 180.0) / 360.0
	latMinr := c.minLoc[0] * math.Pi / 180
	minY := (1.0 - math.Log(math.Tan(latMinr)+(1.0/math.Cos(latMinr)))/math.Pi) / 2.0
	if c.minLoc[0] == c.maxLoc[0] && c.minLoc[1] == c.maxLoc[1] {
		c.zoom = 18
		npx := float64(uint(1)<<uint(c.zoom)) * 256
		c.TilePixelX = int((minX * npx) - (float64(w) / 2))
		c.TilePixelY = int((minY * npx) - (float64(h) / 2))
		c.ntile = int(uint(1) << uint(c.zoom))
	} else {
		latMaxr := c.maxLoc[0] * math.Pi / 180
		dy := ((1.0 - math.Log(math.Tan(latMaxr)+(1.0/math.Cos(latMaxr)))/math.Pi) / 2.0) - minY
		dx := ((c.maxLoc[1] + 180.0) / 360.0) - minX
		for zoom := maxZoom; zoom >= 1; zoom-- {
			npx := float64(uint(1)<<uint(zoom)) * 256.0
			tw := int(npx*dx) + margin
			th := int(npx*dy) + margin
			if tw < w && th < h {
				c.TilePixelX = int(minX*npx) - ((w - tw) / 2)
				c.TilePixelY = int(minY*npx) - ((h - th) / 2)
				c.zoom = zoom
				c.ntile = int(uint(1) << uint(c.zoom))
				return
			}
		}
	}
}

func (c *Conversor) draw(do []drawerObject, w io.Writer) {
	tx, px := divmod(c.TilePixelX, 256)
	ty, py := divmod(c.TilePixelY, 256)
	chanTile := make(chan *tileImage)
	tl := make([]*tileImage, 0)
	tty := ty
	hTile := 256 - py
	iy := py
	yTile := py
	my := 0
	for {
		ttx := tx
		wTile := 256 - px
		ix := px
		xTile := px
		mx := 0
		for {
			if tty >= 0 && tty < c.ntile {
				tl = append(tl, &tileImage{Url: fmt.Sprintf("http://bs.tile.openstreetmap.org/%[1]d/%[2]d/%[3]d.png", c.zoom, ttx%c.ntile, tty),
					Dx: mx,
					Dy: my,
					X:  xTile,
					Y:  yTile,
					W:  wTile,
					H:  hTile})
			}
			ttx++
			ix += 256
			mx += wTile
			if mx >= c.Width {
				break
			}
			wTile = min(256, c.Width-mx)
			xTile = 0
		}
		tty++
		iy += 256
		my += hTile
		if my >= c.Height {
			break
		}
		hTile = min(256, c.Height-my)
		yTile = 0
	}
	downloadTiles(tl, chanTile)
	for ct := range chanTile {
		draw.Draw(c.Img, image.Rect(int(ct.Dx), int(ct.Dy), int(ct.W)+int(ct.Dx), int(ct.H)+int(ct.Dy)), ct.Image, image.Point{int(ct.X), int(ct.Y)}, draw.Src)
	}
	for _, d := range do {
		d.Draw(c)
	}
	png.Encode(w, c.Img)
}

func divmod(numerator, denominator int) (quotient, remainder int) {
	quotient = numerator / denominator
	remainder = numerator % denominator
	return
}*/

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
			//gc.ClosePath()
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

/*//no necesary generic? :-D
func min(x, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}

//math.Mod return negative? WTF
func mod(n, m float64) float64 {
	s := math.Mod(n, m)
	if s < 0 {
		s += m
	}
	return s
}
*/
