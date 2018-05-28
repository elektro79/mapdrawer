package mapDrawer

import (
	"io"
	"log"
	"net/http"
)

const nworkers = 10

var CacheManager cache

var works chan downloadUrl

type downloadUrl interface {
	SetReader(io.Reader)
	GetUrl() string
	CacheId() string
}

type cache interface {
	Read(string) (io.ReadCloser, error)
	Write(string) io.WriteCloser
}

func DownloadUrl(du downloadUrl) {
	works <- du
}

func init() {
	var dc = &dummyCache{}
	CacheManager = dc
	works = make(chan downloadUrl, nworkers)
	for n := 0; n < nworkers; n++ {
		go worker(works, n)
	}
}

func worker(works <-chan downloadUrl, n int) {
	client := &http.Client{}
	for work := range works {
		cid := work.CacheId()
		if r, err := CacheManager.Read(cid); err == nil {
			work.SetReader(r)
			r.Close()
			continue
		}
		req, err := http.NewRequest("GET", work.GetUrl(), nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/66.0.3359.139 Chrome/66.0.3359.139 Safari/537.36")
		response, err := client.Do(req)
		if err != nil {
			log.Println(err)
		} else if response.StatusCode != http.StatusOK {
			log.Println(response.Status)
			response.Body.Close()
			work.SetReader(response.Body)
		} else {
			cw := CacheManager.Write(cid)
			tee := io.TeeReader(response.Body, cw)
			work.SetReader(tee)
			cw.Close()
			response.Body.Close()
		}
	}
}
