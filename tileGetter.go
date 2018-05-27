package mapDrawer

import (
	"errors"
	"image"
	_ "image/png"
	"log"
	"net/http"
	"sync/atomic"
)

type tileImage struct {
	Url   string
	Dx    int
	Dy    int
	X     int
	Y     int
	W     int
	H     int
	Image image.Image
}

func downloadTiles(tiles []*tileImage, r chan *tileImage) {
	total := uint32(len(tiles))
	download := uint32(0)
	for _, tile := range tiles {
		go func(tile *tileImage) {
			b, err := downloadf(tile.Url)
			if err != nil {
				tile.Image = nil
			} else {
				tile.Image = b
			}
			r <- tile
			if total == atomic.AddUint32(&download, 1) {
				close(r)
			}
		}(tile)
	}
}

func downloadf(url string) (image.Image, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/66.0.3359.139 Chrome/66.0.3359.139 Safari/537.36")
	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		log.Println(url)
		log.Fatalln(response.Status)
		return nil, errors.New(response.Status)
	}
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	return img, err
}
