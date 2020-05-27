package scraper

import (
    "fmt"
    "errors"
    "os"
    "github.com/spf13/cobra"
    "github.com/go-scraper/pkg/logger"
)

var TargetSite string
var LocalFile string
var TargetListFile string

var rootCmd = &cobra.Command{
    Use: "scraper",
    Short: "go-scraper scrapes websites for AWS-releated IPs",
}

var scrapeCmd = &cobra.Command{
  Use:   "scrape",
  Short: "Scrapes websites for AWS-releated IPs",
  Long: `The HTML and JS from a site are downloaded and parsed for hostnames and IPs
Any IPs matching AWS-releated services are logged along with the related service`,
  Args: func(cmd *cobra.Command, args []string) error {
      if TargetSite == "" && LocalFile == "" && TargetListFile == "" {
          return errors.New("At least one flag must be specified - see help for details")
      }
      return nil
  },
  RunE: func(cmd *cobra.Command, args []string) error {
      logger.Debug("Validating flags...")

      var err error
      if TargetSite != "" {
          err = ScrapeSite(TargetSite)
      } else if TargetListFile != "" {
          err = ScrapeSiteList(TargetListFile)
      } else {
          err = ScrapeLocalFile(LocalFile)
      }

      return err
  },
}

func Execute() {
//  if err != nil {
//    log.Fatalf("can't initialize zap logger: %v", err)
//  }
//
  rootCmd.AddCommand(scrapeCmd)
  scrapeCmd.Flags().StringVarP(&TargetSite, "target", "t", "", "Target site to scrape")
  scrapeCmd.Flags().StringVarP(&LocalFile, "local", "l", "", "Local file to scrape")
  scrapeCmd.Flags().StringVarP(&TargetListFile, "target-list", "f", "", "Local file with a list of targets to scrape")

  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}
