package scraper

import (
    "os"
    "errors"
    "bufio"
    "go.uber.org/zap"
    "github.com/go-scraper/pkg/logging"
)

func ScrapeSite(targetSite string, verbose bool) error {
    logger := logging.NewLogger(verbose)
    logger.Debug("Scraping site...")
    return nil
}

func ScrapeSiteList(targetSiteListFilename string, verbose bool) error {
    return nil
}

func ScrapeLocalFile(localFilename string, verbose bool) error {
    logger := logging.NewLogger(verbose)
    logger.Debug("Scraping local file...")
    s := Scraper{logger: logger}
    s.localFilename = localFilename
    s.scrapeLocalFile()
    return nil
}

type Scraper struct {
    logger *zap.SugaredLogger
    localFilename string
    targetSiteUrl string
    targetSiteListFilename string
    discoveredHostnames []string
    discoveredIps []string
    listOfAwsServices []string
}

func (scraper *Scraper) scrapeLocalFile() error {
    if scraper.localFilename == "" {
        return errors.New("localFilename cannot be empty")
    }

    file, err := os.Open(scraper.localFilename)

    scraper.check(err, "File could not be opened")
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        scraper.logger.Debug(scanner.Text())
    }

    scraper.check(scanner.Err(), "Scanner failed")
    return nil
}

func (scraper *Scraper) check(e error, msg string) {
    if e != nil {
        scraper.logger.Fatal(msg)
        panic(e)
    }
}


