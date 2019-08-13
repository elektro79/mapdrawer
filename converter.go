package mapDrawer

import (
	"math"

	"github.com/llgcode/draw2d"
)

type Converter struct {
	Width      int
	Height     int
	TilePixelX int
	TilePixelY int
	MinLoc     [2]float64
	MaxLoc     [2]float64
	Zoom       uint8
	ntile      int
	Gc         draw2d.GraphicContext
}

func (c *Converter) GetXY(ll [2]float64) (x, y float64) {
	n := float64(uint(1)<<uint(c.Zoom)) * 256.0
	latr := ll[0] * math.Pi / 180
	x = mod((ll[1]+180.0)/360.0*n-float64(c.TilePixelX), n)
	y = (1.0-math.Log(math.Tan(latr)+(1.0/math.Cos(latr)))/math.Pi)/2.0*n - float64(c.TilePixelY)
	return
}

func (c *Converter) setBounds(do []DrawerObject) {
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
	c.MinLoc = minLoc
	c.MaxLoc = maxLoc
}

func (c *Converter) getBounds() (float64, float64, float64, float64, uint8) {
	npx := float64(c.ntile) * 256.0
	lng1 := float64(c.TilePixelX)*360/npx - 180
	lng2 := float64(c.TilePixelX+c.Width)*360/npx - 180
	y1 := 180 - float64(c.TilePixelY)*360/npx
	y2 := 180 - float64(c.TilePixelY+c.Height)*360/npx
	lat1 := 360/math.Pi*math.Atan(math.Exp(y1*math.Pi/180)) - 90
	lat2 := 360/math.Pi*math.Atan(math.Exp(y2*math.Pi/180)) - 90
	return lat1, lng1, lat2, lng2, c.Zoom
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
func (c *Converter) set(w int, h int, maxZoom uint8, margin int) {
	c.Width = w
	c.Height = h
	minX := (c.MinLoc[1] + 180.0) / 360.0
	latMinr := c.MinLoc[0] * math.Pi / 180
	minY := (1.0 - math.Log(math.Tan(latMinr)+(1.0/math.Cos(latMinr)))/math.Pi) / 2.0
	if c.MinLoc[0] == c.MaxLoc[0] && c.MinLoc[1] == c.MaxLoc[1] {
		c.Zoom = maxZoom
		npx := float64(uint(1)<<uint(c.Zoom)) * 256
		c.TilePixelX = int((minX * npx) - (float64(w) / 2))
		c.TilePixelY = int((minY * npx) - (float64(h) / 2))
		c.ntile = int(uint(1) << uint(c.Zoom))
	} else {
		latMaxr := c.MaxLoc[0] * math.Pi / 180
		dy := ((1.0 - math.Log(math.Tan(latMaxr)+(1.0/math.Cos(latMaxr)))/math.Pi) / 2.0) - minY
		dx := ((c.MaxLoc[1] + 180.0) / 360.0) - minX
		for zoom := maxZoom; zoom >= 0; zoom-- {
			npx := float64(uint(1)<<uint(zoom)) * 256.0
			tw := int(npx*dx) + margin
			th := int(npx*dy) + margin
			if tw < w && th < h {
				c.TilePixelX = int(minX*npx) - ((w - tw) / 2)
				c.TilePixelY = int(minY*npx) - ((h - th) / 2)
				if c.TilePixelY < 0 {
					c.TilePixelY = 0
				}
				c.Zoom = zoom
				c.ntile = int(uint(1) << uint(c.Zoom))
				return
			} else if int(npx) <= w {
				c.TilePixelX = 0
				c.TilePixelY = 0
				c.Zoom = zoom
				c.ntile = int(uint(1) << uint(c.Zoom))
				return
			}

		}
	}
}

func (c *Converter) draw(do []DrawerObject) {
	for _, d := range do {
		d.Draw(c)
	}
}

func divmod(numerator, denominator int) (quotient, remainder int) {
	quotient = numerator / denominator
	remainder = numerator % denominator
	return
}

//no generic necesary? :-D
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

//% return negative? WTF
func modInt(n, m int) int {
	s := n % m
	if s < 0 {
		s += m
	}
	return s
}
