package main

import (
	"testing"

	"gitlab.com/letsboot/core/kubernetes-course/project-vision/internal/crawler"
	"gitlab.com/letsboot/core/kubernetes-course/project-vision/internal/util"
)

func TestCrawl(t *testing.T) {
	// set up flags (viper.Get to retrieve)
	crawler.InitialiseFlags()
	// set up configuration files and parse flags
	util.InitialiseConfig("crawler")

	var (
	// query backend endpoint url from config library
	)

	// create a new crawler sdk client (view internal/sdk)
	_, _ = crawler.Crawl("https://www.sidekicks.ch", 0)

}
