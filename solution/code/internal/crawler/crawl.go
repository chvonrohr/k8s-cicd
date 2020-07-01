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
	node, err := html.Parse(r.Body)
	if err != nil {
		return nil, err
	}
	anchors := crawlNode(node)
	urls := make([]string, 0)
	for _, anchor := range anchors {
		url := getHref(anchor)
		if url != "" {
			urls = append(urls, url)
		}
	}
	return urls, nil
}

func getHref(anchor *html.Node) string {
	for _, attribute := range anchor.Attr {
		if attribute.Key == "href" {
			return attribute.Val
		}
	}
	return ""
}

func crawlNode(node *html.Node) []*html.Node {
	anchors := make([]*html.Node, 0)
	if node.Type == html.ElementNode && node.Data == "a" {
		anchors = append(anchors, node)
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		as := crawlNode(child)
		anchors = append(anchors, as...)
	}
	return anchors
}
