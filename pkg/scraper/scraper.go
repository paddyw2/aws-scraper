package scraper

import (
	"errors"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/paddyw2/urlextract"

	"github.com/paddyw2/aws-scraper/pkg/logging"
)

type DiscoveredUrl struct {
	url        string
    hostname   string
    tld        string
	follow     bool
	aws        bool
	awsService string
}

type Scraper struct {
	logger                 *logging.Logger
	rootHostname           string
	localFilename          string
	targetSiteUrl          string
	targetSiteListFilename string

	discoveredUrls []*DiscoveredUrl
	discoveredIps  []string
}

func (scraper *Scraper) markUrlsAsAwsService() error {
	for _, url := range scraper.discoveredUrls {
        url.aws = true
        if strings.Contains(url.url, "cloudfront") {
            url.awsService = "cloudfront"
        } else if strings.Contains(url.url, "execute-api") {
            url.awsService = "apigateway"
        } else if strings.Contains(url.url, "s3") {
            url.awsService = "s3"
        } else if strings.Contains(url.url, "elb") {
            url.awsService = "elb"
        } else if strings.Contains(url.url, "amazon") || strings.Contains(url.url, "aws") {
            url.awsService = "unknown"
        } else {
            scraper.logger.Debug("Non-AWS hostname: ", url.url)
            url.awsService = "N/A"
            url.aws = false
        }
	}
	return nil
}

func (scraper *Scraper) markUrlsToCheck() error {
	jsExtensionPattern := `(cloudflare|cloudfront).*\.js(\?){0,1}`
	jsRegex := regexp.MustCompile(jsExtensionPattern)

	for _, url := range scraper.discoveredUrls {
		regexMatch := jsRegex.MatchString(url.url)
		rootHostnameMatch, _ := regexp.MatchString(`(`+url.hostname+`).*\.js(\?){0,1}`, url.url)
		if regexMatch || rootHostnameMatch {
			url.follow = true
		} else {
			url.follow = false
		}
	}
	return nil
}

func (scraper *Scraper) scrapeLocalFile() error {
	if scraper.localFilename == "" {
		return errors.New("localFilename cannot be empty")
	}

	file, err := os.Open(scraper.localFilename)

	scraper.check(err, "File could not be opened")
	defer file.Close()

	// declare chunk size
	const maxSzBytes = 1024 * 1024

	// create buffer
	buffer := make([]byte, maxSzBytes)

	for {
		// read content to buffer
		readTotal, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				scraper.check(err, "Erro during file read")
			}
			break
		}
		stringChunk := string(buffer[:readTotal])
		scraper.extractHostnamesIps(stringChunk)
	}
	return nil
}

func (scraper *Scraper) check(e error, msg string) {
	if e != nil {
		scraper.logger.Fatal(msg, e)
		panic(e)
	}
}

func (scraper *Scraper) extractHostnamesIps(line string) {
    urlextractor := urlextract.NewExtractor()
    urlextractor.ExtractHostnamesIps(line)
    for _, ip := range urlextractor.Ips {
	    scraper.discoveredIps = append(scraper.discoveredIps, ip)
    }
    for _, url := range urlextractor.Urls {
        discoveredUrl := DiscoveredUrl{url: url.Url, hostname: url.Hostname, tld: url.Tld}
	    scraper.discoveredUrls = append(scraper.discoveredUrls, &discoveredUrl)
    }
}
