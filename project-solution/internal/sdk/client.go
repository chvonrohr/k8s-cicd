package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"gitlab.com/letsboot/core/kubernetes-course/project-solution/internal/model"
)

// Client is the main resource used to interact with the crawler application.
type Client struct {
	HttpClient *http.Client
	Endpoint   string
}

// PageResponse is a response as given by the crawler. This is mainly used internally.
type PageResponse struct {
	StatusCode  int
	Urls        []string
	ContentType string
	Ok          bool
	Dumped      bool
}

type SchedulerResponse struct {
	ScheduledCount  int
	ScheduledCrawls []model.Crawl
}

// NewClient returns a new client to the given backend endpoint.
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

// PageCallback is used by the crawler to call back a page response to the backend.
func (c *Client) PageCallback(page model.Page, response PageResponse) error {
	bs, err := json.Marshal(&response)
	if err != nil {
		return err
	}
	r := bytes.NewReader(bs)
	_, err = c.HttpClient.Post(fmt.Sprintf("%s/callback/%d", c.Endpoint, page.ID), "application/json", r)
	return err
}

// CreateSite can be used to create a new site with a given struct.
func (c *Client) CreateSite(site model.Site) (model.Site, error) {
	bs, err := json.Marshal(&site)
	if err != nil {
		return site, err
	}
	r := bytes.NewReader(bs)
	response, err := c.HttpClient.Post(fmt.Sprintf("%s/sites", c.Endpoint), "application/json", r)
	if err != nil {
		return site, err
	}
	bs, err = ioutil.ReadAll(response.Body)
	err = json.Unmarshal(bs, &site)
	if err != nil {
		return model.Site{}, err
	}
	return site, err
}

// CreateCrawl can be used to create a new crawl with a given struct.
func (c *Client) CreateCrawl(crawl model.Crawl) error {
	bs, err := json.Marshal(&crawl)
	if err != nil {
		return err
	}
	r := bytes.NewReader(bs)
	_, err = c.HttpClient.Post(fmt.Sprintf("%s/crawls", c.Endpoint), "application/json", r)
	return err
}

// GetSites can be used to get a list of sites from the api.
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

// GetCrawlsForSite returns a list of all crawls fora  given site.
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

// GetPagesForCrawl returns a list of pages for a given crawl.
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

func (c *Client) Schedule() (SchedulerResponse, error) {
	r, err := c.HttpClient.Post(fmt.Sprintf("%s/schedule", c.Endpoint), "", nil)
	if err != nil {
		return SchedulerResponse{}, err
	}
	var response *SchedulerResponse
	bs, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(bs, &response)
	if err != nil {
		return SchedulerResponse{}, err
	}
	return *response, err
}
