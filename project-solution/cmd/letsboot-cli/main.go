package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"gitlab.com/letsboot/core/kubernetes-course/project-solution/internal/model"
	"gitlab.com/letsboot/core/kubernetes-course/project-solution/internal/sdk"
)

func main() {

	var endpoint string
	var rootCmd = &cobra.Command{Use: "letsboot-cli"}

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
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := sdk.NewClient(endpoint)
			url := args[0]

			interval, err := cmd.Flags().GetDuration("interval")
			if err != nil {
				log.Fatal(err)
			}

			site, err := client.CreateSite(model.Site{
				Url:      url,
				Interval: interval,
			})
			if err != nil {
				log.Fatal(err)
			}

			c, err := cmd.Flags().GetBool("crawl")
			if err != nil {
				log.Fatal(err)
			}

			if c {
				err = client.CreateCrawl(model.Crawl{SiteID: site.ID})
				if err != nil {
					log.Fatal(err)
				}
			}
		},
	}

	cmdSiteCreate.Flags().DurationP("interval", "i", 7*24*time.Hour, "crawl interval")
	cmdSiteCreate.Flags().Bool("crawl", false, "immediately start crawl")

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

	var cmdSchedule = &cobra.Command{
		Use:   "schedule",
		Short: "schedule all pending crawls",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			client := sdk.NewClient(endpoint)
			r, err := client.Schedule()
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Scheduled %d crawls: %v", r.ScheduledCount, r.ScheduledCrawls)
		},
	}

	var cmdCompletions = &cobra.Command{
		Use:   "completions",
		Short: "generate completions for various shells",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	var cmdCompletionsZsh = &cobra.Command{
		Use:   "zsh",
		Short: "generates zsh completions for cli, outputs to stdout",
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Root().GenZshCompletion(os.Stdout)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	cmdCompletions.AddCommand(cmdCompletionsZsh)

	var cmdCompletionsBash = &cobra.Command{
		Use:   "bash",
		Short: "generates bash completions for cli, outputs to stdout",
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Root().GenBashCompletion(os.Stdout)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	cmdCompletions.AddCommand(cmdCompletionsBash)

	var cmdCompletionsPowerShell = &cobra.Command{
		Use:   "powershell",
		Short: "generates powershell completions for cli, outputs to stdout",
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Root().GenPowerShellCompletion(os.Stdout)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	cmdCompletions.AddCommand(cmdCompletionsPowerShell)

	var cmdCompletionsFish = &cobra.Command{
		Use:   "fish",
		Short: "generates fish completions for cli, outputs to stdout",
		Run: func(cmd *cobra.Command, args []string) {
			desc, err := cmd.Flags().GetBool("desc")
			if err != nil {
				log.Fatal(err)
			}
			err = cmd.Root().GenFishCompletion(os.Stdout, desc)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	cmdCompletionsFish.Flags().Bool("desc", false, "include desc")
	cmdCompletions.AddCommand(cmdCompletionsFish)

	rootCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "e", "http://localhost:8080", "endpoint for backend")
	rootCmd.AddCommand(cmdSites, cmdCrawls, cmdPages, cmdSchedule, cmdCompletions)
	cmdSites.AddCommand(cmdSiteCreate)
	cmdCrawls.AddCommand(cmdCrawlCreate)
	_ = rootCmd.Execute()

}
