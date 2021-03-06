package scraper

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
    "strings"
    "bufio"

	"github.com/paddyw2/aws-scraper/pkg/logging"
)

var OutputHeaderPrinted = false

type ScrapeController interface {
	ScrapeSite(targetSite string) error
	ScrapeLocalFile(hostname string, localFilename string) error
	ScrapeSiteList(targetSiteListFilename string) error
}

type scrapeController struct {
	logger       *logging.Logger
	verboseLevel int
    displayIps   bool
	maxLevel     int
	currentLevel int
}

func NewScrapeController(logger *logging.Logger, verboseLevel int, displayIps bool) ScrapeController {
	var maxLevel int = 1
    sc := scrapeController{logger: logger, verboseLevel: verboseLevel, maxLevel: maxLevel, currentLevel: 0, displayIps: displayIps}
	return &sc
}

func (sc *scrapeController) check(e error, msg string) {
	if e != nil {
		sc.logger.Fatal(msg, e)
		panic(e)
	}
}

func (sc *scrapeController) ScrapeSite(targetSite string) error {
    cleanSiteName := strings.Replace(targetSite, "http://", "", -1)
    cleanSiteName = strings.Replace(targetSite, "https://", "", -1)
    cleanSiteName = strings.Replace(cleanSiteName, "/", "", -1)
	fileName := "/tmp/" + cleanSiteName + "-source.txt"
	sc.logger.Info("Downloading " + targetSite + " to " + fileName)
	err := downloadFile(fileName, targetSite)
	if err != nil {
		sc.logger.Fatal("Download did not work", err)
		return errors.New("Download did not work")
	}
	sc.ScrapeLocalFile(targetSite, fileName)
	return nil
}

func (sc *scrapeController) ScrapeSiteList(targetSiteListFilename string) error {

    file, err := os.Open(targetSiteListFilename)

    sc.check(err, "File could not be opened")
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        sc.logger.Debug("Reading: ", line)
        targetSite := strings.TrimSpace(line)
        sc.ScrapeSite(targetSite)
    }

    sc.check(scanner.Err(), "Scanner failed")
	return nil
}

func (sc *scrapeController) ScrapeLocalFile(hostname string, localFilename string) error {
	sc.logger.Info("Scraping local file...")
	logger := logging.NewLogger(sc.verboseLevel)
	s := Scraper{logger: logger, rootHostname: hostname}
	s.localFilename = localFilename
	s.scrapeLocalFile()
	s.markUrlsToCheck()
	s.markUrlsAsAwsService()
	for _, url := range s.discoveredUrls {

		if url.follow {
			logger.Debug("Checking: ", url, " with: ", url.url)
			if sc.currentLevel < sc.maxLevel {
				sc.ScrapeSite(url.url)
				sc.currentLevel += 1
			}
		}
		if url.aws {
            outputResultHeader()
			logger.Info("AWS: ", url.hostname)
            outputResult(hostname, "hostname", url.hostname, url.awsService)
		}
	}
    if sc.displayIps {
        for _, ip := range s.discoveredIps {
            outputResultHeader()
            logger.Info("IP: ", ip)
            outputResult(hostname, "ip", ip, "N/A")
        }
    }
	sc.currentLevel = 0
	return nil
}

func outputResultHeader() {
    headerLine := "hostname, resultType, result, awsService"
    if !OutputHeaderPrinted {
        fmt.Println(headerLine)
        OutputHeaderPrinted = true
    }
}

func outputResult(hostname string, resultType string, result string, awsService string) {
    resultLine := hostname + ", " + resultType + ", " + result + ", " + awsService
    fmt.Println(resultLine)
}

func downloadFile(filepath string, url string) error {
	fullUrl := url
	if httpMatch, _ := regexp.MatchString(`^http`, url); !httpMatch {
		fullUrl = "http://" + url
	}
	// Get the data
	resp, err := http.Get(fullUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
