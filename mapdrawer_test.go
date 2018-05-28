package mapDrawer

import (
	"image"
	"image/jpeg"
	"log"
	"os"
	"os/user"
	"testing"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func TestMain(t *testing.T) {
	var err error
	cu, _ := user.Current()
	CacheManager, err = NewCacheHd(cu.HomeDir + "/.cache/testMapDrawer/")
	if err != nil {
		log.Fatal(err)
	}
	draw2d.SetFontFolder(cu.HomeDir + "/ttf")
	img := image.NewRGBA(image.Rect(0, 0, 1024, 1024))
	d := NewDrawer(draw2dimg.NewGraphicContext(img), 1024, 1024, 30)
	marker := &Marker{FillColor: "#03DC0380", Locs: [][2]float64{
		{39.478706359863, -0.3582658469677},
	}}
	d.Add(marker)
	marker = &Marker{FillColor: "#03DC03ff", Size: tiny, Label: "H", Locs: [][2]float64{
		{37.181018829346, -5.7789669036865},
	}}
	d.Add(marker)

	path := &Path{Color: "#03DC03ff", GeoDesic: true, Arrow: true, Weight: 3.0, Locs: [][2]float64{
		{38.036354064941, 140.16046142578},
		{19.675472259521, -155.97915649414},
	}}
	d.Add(path)
	path = &Path{Color: "#03DC03ff", GeoDesic: true, Arrow: true, Weight: 3.0, Locs: [][2]float64{
		{37.181018829346, -5.7789669036865},
		{63.4382591247559, 10.3109569549561},
	}}
	d.Add(path)
	marker = &Marker{FillColor: "#03DC03ff", Label: "8", Locs: [][2]float64{
		{38.036354064941, 140.16046142578},
		{19.675472259521, -155.97915649414},
	}}
	d.Add(marker)
	f, err := os.Create("/tmp/img.jpg")
	check(err)
	defer f.Close()
	d.Draw()
	jpeg.Encode(f, img, &jpeg.Options{Quality: 90})
}
