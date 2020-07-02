package crawler

import (
	"fmt"
	"gitlab.com/letsboot/core/kubernetes-course/solution/code/core/internal/sdk"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func Crawl(uri string) (response sdk.PageResponse, err error) {
	request, _ := http.NewRequest("GET", uri, nil)
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36")
	r, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}
	response.StatusCode = r.StatusCode
	if r.StatusCode < 200 || r.StatusCode >= 300 {
		return response, fmt.Errorf("invalid status code: %s", r.Status)
	}
	response.ContentType = r.Header.Get("content-type")
	if !strings.Contains(response.ContentType, "text/html") {
		return response, fmt.Errorf("invalid content type: %s", response.ContentType)
	}
	response.Ok = true
	node, err := html.Parse(r.Body)
	if err != nil {
		return response, err
	}
	parsedUri, _ := url.Parse(uri)
	anchors := crawlNode(node)
	response.Urls = make([]string, 0)
	for _, anchor := range anchors {
		anchorUri := getHref(anchor)
		if strings.HasPrefix(anchorUri, "/") {
			anchorUri = fmt.Sprintf("%s://%s%s", parsedUri.Scheme, parsedUri.Host, anchorUri)
		}
		if isValidUri(parsedUri, anchorUri) {
			response.Urls = append(response.Urls, anchorUri)
		}
	}
	return
}

func isValidUri(parent *url.URL, uri string) bool {
	if uri == "" {
		return false
	}
	parsedUri, err := url.Parse(uri)
	if err != nil {
		log.Println(err)
		return false
	}
	return parent.Host == parsedUri.Host
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
