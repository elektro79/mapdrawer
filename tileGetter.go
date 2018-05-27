package mapDrawer

import (
	"io"
	"log"
	"net/http"
)

const nworkers = 10

var works = make(chan downloadUrl, 10)

type downloadUrl interface {
	SetReaderClose(io.ReadCloser)
	GetUrl() string
}

func DownloadUrl(du downloadUrl) {
	works <- du
}

func init() {
	for n := 0; n < nworkers; n++ {
		go worker(works)
	}
}

func worker(works <-chan downloadUrl) {
	client := &http.Client{}
	for work := range works {
		req, err := http.NewRequest("GET", work.GetUrl(), nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/66.0.3359.139 Chrome/66.0.3359.139 Safari/537.36")
		response, err := client.Do(req)
		if err != nil {
			log.Println(err)
		} else if response.StatusCode != http.StatusOK {
			log.Println(response.Status)
			response.Body.Close()
		}
		work.SetReaderClose(response.Body)
	}
}
