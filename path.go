package mapDrawer

import (
	"image/color"
	"math"

	"github.com/llgcode/draw2d"
)

type Path struct {
	Color    string       `json:"color"`
	Locs     [][2]float64 `json:"locs"`
	Weight   float64      `json:"weight"`
	GeoDesic bool         `json:"geodesic"`
	Arrow    bool         `json:"arrow"`
	locs     [][2]float64
}

func (p *Path) Draw(c *Converter) {
	var col color.Color
	var err error
	if col, err = StringToColor(p.Color); err != nil {
		col = color.Black
	}
	gc := c.Gc
	gc.SetLineWidth(p.Weight)
	gc.SetLineCap(draw2d.RoundCap)
	gc.SetLineJoin(draw2d.RoundJoin)
	locs := p.Locs
	if p.GeoDesic {
		p.setGeodesic()
		locs = p.locs
	}
	x, y := c.GetXY(locs[0])
	gc.MoveTo(x, y)
	px := x
	dir := 0.0
	for _, loc := range locs {
		x, y = c.GetXY(loc)
		px -= x
		if c.Zoom > 2 || dir == 0 || (dir > 0 && px > 0) || (dir < 0 && px < 0) {
			gc.LineTo(x, y)
		} else {
			gc.MoveTo(x, y)
		}
		dir = px
		px = x
	}
	gc.SetFillColor(col)
	gc.SetStrokeColor(col)
	gc.Stroke()
	if p.Arrow {
		px, py := c.GetXY(locs[len(locs)-2])
		ang := math.Atan2(px-x, py-y)
		gc.BeginPath()
		gc.SetLineWidth(1.0)
		gc.MoveTo(x, y)
		gc.LineTo(math.Sin(ang-0.39)*10+x, math.Cos(ang-0.39)*10+y)
		gc.LineTo(math.Sin(ang+0.39)*10+x, math.Cos(ang+0.39)*10+y)
		gc.LineTo(x, y)
		gc.FillStroke()
	}
}

func (p *Path) genGeodesic(pend [2]float64) {
	pinit := p.locs[len(p.locs)-1]
	R := 6367000.0
	lat1 := pinit[0] * math.Pi / 180
	lat2 := pend[0] * math.Pi / 180
	lng1 := pinit[1] * math.Pi / 180
	lng2 := pend[1] * math.Pi / 180
	dLng := lng2 - lng1
	segments := (math.Floor(math.Abs(dLng * R / 50000)))

	if segments > 1 {
		d := 2 * math.Asin(math.Sqrt(math.Pow((math.Sin((lat1-lat2)/2)), 2)+math.Cos(lat1)*math.Cos(lat2)*math.Pow((math.Sin((lng1-lng2)/2)), 2)))
		sinLat1 := math.Sin(lat1)
		sinLat2 := math.Sin(lat2)
		cosLat1 := math.Cos(lat1)
		cosLat2 := math.Cos(lat2)
		cosLat1cosLng1 := cosLat1 * math.Cos(lng1)
		cosLat1sinLng1 := cosLat1 * math.Sin(lng1)
		cosLat2cosLng2 := cosLat2 * math.Cos(lng2)
		cosLat2SinLng2 := cosLat2 * math.Sin(lng2)

		for i := 1.0; i < segments; i++ {
			f := i / segments
			A := math.Sin((1-f)*d) / math.Sin(d)
			B := math.Sin(f*d) / math.Sin(d)
			x := A*cosLat1cosLng1 + B*cosLat2cosLng2
			y := A*cosLat1sinLng1 + B*cosLat2SinLng2
			z := A*sinLat1 + B*sinLat2
			latr := math.Atan2(z, math.Sqrt(math.Pow(x, 2.0)+math.Pow(y, 2.0)))
			lngr := math.Atan2(y, x)
			p.locs = append(p.locs, [2]float64{latr * 180 / math.Pi, lngr * 180 / math.Pi})
		}
	}
	p.locs = append(p.locs, pend)
}

func (p *Path) genGeodesic2(pend [2]float64) {
	pinit := p.locs[len(p.locs)-1]
	R := 6367000.0
	d2r := math.Pi / 180.0
	r2d := 180.0 / math.Pi
	lat1 := pinit[0] * d2r
	lat2 := pend[0] * d2r
	lng1 := pinit[1] * d2r
	lng2 := pend[1] * d2r
	dLng := lng2 - lng1
	segments := (math.Floor(math.Abs(dLng * R / 50000)))
	if segments > 1 {
		sinLat1 := math.Sin(lat1)
		sinLat2 := math.Sin(lat2)
		cosLat1 := math.Cos(lat1)
		cosLat2 := math.Cos(lat2)
		sinLat1CosLat2 := sinLat1 * cosLat2
		sinLat2CosLat1 := sinLat2 * cosLat1
		cosLat1CosLat2SinDLng := cosLat1 * cosLat2 * math.Sin(dLng)
		i := 1.0
		for i < segments {
			iLng := lng1 + dLng*(i/segments)
			iLat := math.Atan((sinLat1CosLat2*math.Sin(lng2-iLng) + sinLat2CosLat1*math.Sin(iLng-lng1)) / cosLat1CosLat2SinDLng)
			nloc := [2]float64{iLat * r2d, iLng * r2d}
			p.locs = append(p.locs, nloc)
			i += 1.0
		}
	}
	p.locs = append(p.locs, pend)
}

func (p *Path) setGeodesic() {
	if len(p.locs) == 0 {
		p.locs = append(p.locs, p.Locs[0])
		for _, loc := range p.Locs[1:] {
			p.genGeodesic(loc)
		}
	}
}

func (p *Path) GetLocs() [][2]float64 {
	if p.GeoDesic {
		p.setGeodesic()
		return p.locs
	} else {
		return p.Locs
	}
}
