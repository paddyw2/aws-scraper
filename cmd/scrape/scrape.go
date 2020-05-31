package scrape

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/go-scraper/pkg/logging"
	"github.com/go-scraper/pkg/scraper"
)

var VerboseLevel int = 0
var TargetSite string
var LocalFile string
var TargetListFile string

var rootCmd = &cobra.Command{
	Use:   "scraper",
	Short: "go-scraper scrapes websites for AWS-releated IPs",
}

var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "Scrapes websites for AWS-releated IPs",
	Long: `The HTML and JS from a site are downloaded and parsed for hostnames and IPs
Any IPs matching AWS-releated services are logged along with the related service`,
	Args: ScrapeArgsValidator,
	RunE: ScrapeCommand,
}

func ScrapeArgsValidator(cmd *cobra.Command, args []string) error {
    if TargetSite == "" && LocalFile == "" && TargetListFile == "" {
        return errors.New("At least one flag must be specified - see help for details")
    } else if TargetListFile == "" && TargetSite == "" {
        return errors.New("The name of the site must be specified")
    }
    return nil
}

func ScrapeCommand(cmd *cobra.Command, args []string) error {
    logger := logging.NewLogger(VerboseLevel)
    logger.Debug("Validating flags...")

    scrapeController := scraper.NewScrapeController(logger, VerboseLevel)
    var err error
    if LocalFile != "" {
        err = scrapeController.ScrapeLocalFile(TargetSite, LocalFile)
    } else if TargetListFile != "" {
        err = scrapeController.ScrapeSiteList(TargetListFile)
    } else {
        err = scrapeController.ScrapeSite(TargetSite)
    }

    return err
}

func Execute() {
	rootCmd.AddCommand(scrapeCmd)
	rootCmd.PersistentFlags().CountVarP(&VerboseLevel, "verbose", "v", "Verbose logging level 1")
	scrapeCmd.Flags().StringVarP(&TargetSite, "target", "t", "", "Target site to scrape")
	scrapeCmd.Flags().StringVarP(&LocalFile, "local", "l", "", "Local file to scrape")
	scrapeCmd.Flags().StringVarP(&TargetListFile, "target-list", "f", "", "Local file with a list of targets to scrape")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
