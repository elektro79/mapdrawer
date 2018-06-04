package mapDrawer

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"io"
	"log"
	"math/rand"
	"sync"
)

type tileImage struct {
	imageTile *ImageTile
	tilex     int
	tiley     int
	zoom      uint8
	Dx        int
	Dy        int
	X         int
	Y         int
	W         int
	H         int
	complete  chan bool
	mutex     sync.Locker
	Image     draw.Image
}

func (t *tileImage) GetUrl() string {
	return fmt.Sprintf(t.imageTile.Url, t.imageTile.Mirror[rand.Intn(len(t.imageTile.Mirror))], t.zoom, t.tilex, t.tiley)
}

func (t *tileImage) CacheId() string {
	return fmt.Sprintf("%s-%d-%d-%d", t.imageTile.Name, t.zoom, t.tilex, t.tiley)
}

func (t *tileImage) SetReader(rc io.Reader) {
	img, _, err := image.Decode(rc)
	if err != nil {
		log.Println(err)
		t.complete <- true
		return
	}
	t.mutex.Lock()
	draw.Draw(t.Image, image.Rect(t.Dx, t.Dy, t.W+t.Dx, t.H+t.Dy), img, image.Point{t.X, t.Y}, draw.Src)
	t.mutex.Unlock()
	t.complete <- true

}

type ImageTile struct {
	Url    string
	Name   string
	Mirror []string
}

func ImageTileStreetMap() *ImageTile {
	return &ImageTile{Name: "OSM", Url: "http://%[1]s.tile.openstreetmap.org/%[2]d/%[3]d/%[4]d.png", Mirror: []string{"a", "b"}}
}

func (i *ImageTile) Draw(c *Converter) {
	tx, px := divmod(modInt(c.TilePixelX, 256*c.ntile), 256)
	ty, py := divmod(c.TilePixelY, 256)
	tty := ty
	hTile := 256 - py
	iy := py
	yTile := py
	my := 0
	var mutex sync.Mutex
	complete := make(chan bool)
	img := image.NewRGBA(image.Rect(0, 0, c.Width, c.Height))
	cnt := 0
	go func() {
		for {
			ttx := tx
			wTile := 256 - px
			ix := px
			xTile := px
			mx := 0
			for {
				if tty >= 0 && tty < c.ntile && mx < (256*c.ntile) {
					cnt++
					DownloadUrl(&tileImage{
						imageTile: i,
						tilex:     modInt(ttx, c.ntile),
						tiley:     tty,
						zoom:      c.Zoom,
						Dx:        mx,
						Dy:        my,
						X:         xTile,
						Y:         yTile,
						W:         wTile,
						H:         hTile,
						complete:  complete,
						mutex:     &mutex,
						Image:     img})
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
	}()
	for _ = range complete {
		cnt--
		if cnt == 0 {
			close(complete)
		}
	}
	c.Gc.DrawImage(img)
}
