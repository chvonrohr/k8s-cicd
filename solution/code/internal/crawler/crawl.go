package crawler

import (
	"errors"
	"golang.org/x/net/html"
	"net/http"
)

func Crawl(url string) ([]string, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.New("invalid status code")
	}
	if r.Header.Get("content-type") != "text/html" {
		return nil, errors.New("invalid content type")
	}
	z := html.NewTokenizer(r.Body)
}
