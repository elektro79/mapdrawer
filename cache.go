package mapDrawer

import (
	"errors"
	"io"
	"log"
	"os"
)

type ignore int

func (i ignore) Write(b []byte) (int, error) { return len(b), nil }
func (i ignore) Close() error                { return nil }

type dummyCache struct{}

func (d *dummyCache) Read(url string) (io.ReadCloser, error) {
	return nil, errors.New("No cache")
}

func (d *dummyCache) Write(url string) io.WriteCloser {
	var dc ignore
	return &dc
}

type cacheHd struct {
	dir string
}

func NewCacheHd(dir string) (*cacheHd, error) {
	if err := os.MkdirAll(dir, 0770); err != nil {
		return nil, err
	}
	return &cacheHd{dir: dir}, nil
}

func (c *cacheHd) Read(uid string) (io.ReadCloser, error) {
	f, err := os.Open(c.dir + uid)
	if err != nil {
		return nil, err
	}
	return f, nil
}
func (c *cacheHd) Write(uid string) io.WriteCloser {
	f, err := os.Create(c.dir + uid)
	if err != nil {
		var dc ignore
		log.Println(err)
		return &dc
	}
	return f
}
