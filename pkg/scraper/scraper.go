package scraper

import (
    "os"
    "errors"
    "bufio"
    "regexp"
    "go.uber.org/zap"
    "github.com/go-scraper/pkg/logging"
    "github.com/go-scraper/pkg/tlds"
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
        line := scanner.Text()
        scraper.logger.Debug("Reading: ", line)
        scraper.extractHostnamesIps(line)
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

func (scraper *Scraper) extractHostnamesIps(line string) {
	ipPattern := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
    submatchall := ipPattern.FindAllString(line, -1)
	for _, element := range submatchall {
        scraper.logger.Debug("--> ip: ", element)
        scraper.discoveredIps = append(scraper.discoveredIps, element)
	}

	hostnamePattern := regexp.MustCompile(`["'/]+([a-zA-Z0-9]+(-[a-zA-Z0-9]+)*)+(\.([a-zA-Z]+(-[a-zA-Z]+)*))+`)
    submatchall = hostnamePattern.FindAllString(line, -1)
	for _, element := range submatchall {
        tldPattern := regexp.MustCompile(`\.([a-zA-Z]+(-[a-zA-Z]+)*)$`)
        tldMatch := tldPattern.FindString(element)
        if _, ok := tlds.TLDS[tldMatch[1:]]; ok {
            scraper.logger.Debug("--> hostname: ", element)
            scraper.discoveredHostnames = append(scraper.discoveredHostnames, element)
        }
	}

}
