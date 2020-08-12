package crawler

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/spf13/viper"
	"gitlab.com/letsboot/core/kubernetes-course/solution/code/core/internal/sdk"
	"golang.org/x/net/html"
)

var (
	mkdirOnce = sync.Once{}
)

// Crawl processes a uri by making a GET request to it and extracting anchor tags from it.
// It then returns a page response containing information about the request.
//
// Optionally, it can also dump the raw http response body to a file (it will do so if the flag `crawler.dump` is set.
// In this case, it will write the file to the given data directory `crawler.data`.
func Crawl(uri string, crawlId int) (response sdk.PageResponse, err error) {
	request, _ := http.NewRequest("GET", uri, nil)
	// set a custom user agent - some websites block default library user agents like the go useragent
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
	var reader io.Reader = r.Body

	if viper.GetBool("crawler.dump") {
		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return response, err
		}
		dataDir := viper.GetString("crawler.data")
		mkdirOnce.Do(func() {
			err = os.MkdirAll(dataDir, os.ModeDir|os.ModePerm)
		})
		if err != nil {
			return response, err
		}
		sanitizedUri := strings.TrimPrefix(strings.TrimPrefix(uri, "http://"), "https://")
		p := fmt.Sprintf("%s/%04d/%s", dataDir, crawlId, sanitizedUri)
		err = os.MkdirAll(path.Dir(p), os.ModeDir|os.ModePerm)
		if err != nil {
			return response, err
		}
		err = ioutil.WriteFile(p, bs, os.ModePerm)
		if err != nil {
			return response, err
		}
		reader = bytes.NewReader(bs)
	}

	node, err := html.Parse(reader)
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
