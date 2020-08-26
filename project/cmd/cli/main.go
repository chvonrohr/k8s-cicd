package main

import (
	"github.com/spf13/cobra"
	"gitlab.com/letsboot/core/kubernetes-course/solution/code/core/internal/model"
	"gitlab.com/letsboot/core/kubernetes-course/solution/code/core/internal/sdk"
	"log"
	"strconv"
	"time"
)

func main() {

	var endpoint string

	var cmdSites = &cobra.Command{
		Use:   "sites",
		Short: "list all sites",
		Long:  "Retrieve a list of all sites from configured backend.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			client := sdk.NewClient(endpoint)
			sites, err := client.GetSites()
			if err != nil {
				log.Fatal(err)
			}
			for _, site := range sites {
				log.Printf("%04d - %s", site.ID, site.Url)
			}
		},
	}

	var cmdSiteCreate = &cobra.Command{
		Use:   "create",
		Short: "create a new site",
		Long:  "Create a new site with a given url and optional interval.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := sdk.NewClient(endpoint)
			url := args[0]
			var interval time.Duration
			if len(args) == 1 {
				interval = 7 * 24 * time.Hour
			} else {
				i, err := strconv.Atoi(args[1])
				if err != nil {
					log.Fatal(err)
				}
				interval = time.Duration(i)
			}

			err := client.CreateSite(model.Site{
				Url:      url,
				Interval: interval,
			})
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	var cmdCrawls = &cobra.Command{
		Use:   "crawls [site-id]",
		Short: "list all crawls for site",
		Long:  "Retrieve a list of all crawls for given site id from configured backend.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := sdk.NewClient(endpoint)
			siteId, err := strconv.Atoi(args[0])
			if err != nil {
				log.Fatal(err)
			}
			crawls, err := client.GetCrawlsForSite(siteId)
			if err != nil {
				log.Fatal(err)
			}
			for _, crawl := range crawls {
				log.Printf("%04d - %s - %s", crawl.ID, crawl.CreatedAt, crawl.Site.Url)
			}
		},
	}

	var cmdCrawlCreate = &cobra.Command{
		Use:   "create",
		Short: "create a new crawl",
		Long:  "Create a new crawl for a given site. This immediately queues the crawl.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := sdk.NewClient(endpoint)
			siteId, err := strconv.Atoi(args[0])
			if err != nil {
				log.Fatal(err)
			}
			err = client.CreateCrawl(model.Crawl{SiteID: uint(siteId)})
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	var cmdPages = &cobra.Command{
		Use:   "pages [site-id]",
		Short: "list all pages for crawl",
		Long:  "Retrieve a list of all pages for given crawl id from configured backend.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := sdk.NewClient(endpoint)
			crawlId, err := strconv.Atoi(args[0])
			if err != nil {
				log.Fatal(err)
			}
			pages, err := client.GetPagesForCrawl(crawlId)
			if err != nil {
				log.Fatal(err)
			}
			for _, page := range pages {
				log.Printf("%04d - [%s, %03d] - %s", page.ID, page.State, page.StatusCode, page.Url)
			}
		},
	}

	var rootCmd = &cobra.Command{Use: "cli"}
	rootCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "e", "http://localhost:8080", "endpoint for backend")
	rootCmd.AddCommand(cmdSites, cmdCrawls, cmdPages)
	cmdSites.AddCommand(cmdSiteCreate)
	cmdCrawls.AddCommand(cmdCrawlCreate)
	_ = rootCmd.Execute()

}
