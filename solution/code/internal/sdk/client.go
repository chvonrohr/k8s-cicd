package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"gitlab.com/letsboot/core/kubernetes-course/solution/code/core/internal/model"
)

type Client struct {
	HttpClient *http.Client
	Endpoint   string
}
type PageResponse struct {
	StatusCode  int
	Urls        []string
	ContentType string
	Ok          bool
}

func NewClient(endpoint string) *Client {
	return &Client{
		HttpClient: &http.Client{
			Transport: &http.Transport{
				TLSHandshakeTimeout:   30 * time.Second,
				IdleConnTimeout:       30 * time.Second,
				ResponseHeaderTimeout: 30 * time.Second,
				ExpectContinueTimeout: 30 * time.Second,
			},
			Timeout: 30 * time.Second,
		},
		Endpoint: endpoint,
	}
}

func (c *Client) PageCallback(page model.Page, response PageResponse) error {
	bs, err := json.Marshal(&response)
	if err != nil {
		return err
	}
	r := bytes.NewReader(bs)
	_, err = c.HttpClient.Post(fmt.Sprintf("%s/callback/%d", c.Endpoint, page.ID), "application/json", r)
	return err
}

func (c *Client) CreateSite(site model.Site) error {
	bs, err := json.Marshal(&site)
	if err != nil {
		return err
	}
	r := bytes.NewReader(bs)
	_, err = c.HttpClient.Post(fmt.Sprintf("%s/sites", c.Endpoint), "application/json", r)
	return err
}
func (c *Client) CreateCrawl(crawl model.Crawl) error {
	bs, err := json.Marshal(&crawl)
	if err != nil {
		return err
	}
	r := bytes.NewReader(bs)
	_, err = c.HttpClient.Post(fmt.Sprintf("%s/crawls", c.Endpoint), "application/json", r)
	return err
}

func (c *Client) GetSites() ([]model.Site, error) {
	r, err := c.HttpClient.Get(fmt.Sprintf("%s/sites", c.Endpoint))
	if err != nil {
		return nil, err
	}
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var sites []model.Site
	err = json.Unmarshal(bs, &sites)
	return sites, err
}
func (c *Client) GetCrawlsForSite(site int) ([]model.Crawl, error) {
	r, err := c.HttpClient.Get(fmt.Sprintf("%s/crawls?site=%d", c.Endpoint, site))
	if err != nil {
		return nil, err
	}
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var crawls []model.Crawl
	err = json.Unmarshal(bs, &crawls)
	return crawls, err
}
func (c *Client) GetPagesForCrawl(crawl int) ([]model.Page, error) {
	r, err := c.HttpClient.Get(fmt.Sprintf("%s/pages?crawl=%d", c.Endpoint, crawl))
	if err != nil {
		return nil, err
	}
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var pages []model.Page
	err = json.Unmarshal(bs, &pages)
	return pages, err
}
