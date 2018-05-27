package mapDrawer

import (
	"image/color"
	"os"
	"os/user"
	"testing"

	"github.com/llgcode/draw2d"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func TestMain(t *testing.T) {
	cu, _ := user.Current()
	draw2d.SetFontFolder(cu.HomeDir + "/ttf")
	d := NewDrawer()
	d.width = 1024
	d.height = 1024
	d.margin = 30
	marker := &Marker{Color: color.RGBA{0x03, 0xDC, 0x03, 0xff}, Label: "8", Locs: [][2]float64{
		{39.478706359863, -0.3582658469677},
	}}
	d.Add(marker)
	marker = &Marker{Color: color.RGBA{0x03, 0xDC, 0x03, 0xff}, Size: tiny, Label: "H", Locs: [][2]float64{
		{37.181018829346, -5.7789669036865},
	}}
	d.Add(marker)

	path := &Path{Color: color.RGBA{0x03, 0xDC, 0x03, 0xff}, GeoDesic: true, Arrow: true, Weight: 3.0, Locs: [][2]float64{
		{38.036354064941, 140.16046142578},
		{19.675472259521, -155.97915649414},
	}}
	d.Add(path)
	path = &Path{Color: color.RGBA{0x03, 0xDC, 0x03, 0xff}, GeoDesic: true, Arrow: true, Weight: 3.0, Locs: [][2]float64{
		{37.181018829346, -5.7789669036865},
		{63.4382591247559, 10.3109569549561},
	}}
	d.Add(path)
	marker = &Marker{Color: color.RGBA{0x03, 0xDC, 0x03, 0xff}, Label: "8", Locs: [][2]float64{
		{38.036354064941, 140.16046142578},
		{19.675472259521, -155.97915649414},
	}}
	d.Add(marker)
	f, err := os.Create("/tmp/img.png")
	check(err)
	defer f.Close()
	d.Draw(f)
}
