package scraper

import (
	"errors"
	"io"
	"os"
	"regexp"

	"github.com/go-scraper/pkg/logging"
	"github.com/go-scraper/pkg/tlds"
)

type DiscoveredUrl struct {
	url        string
	hostname   string
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
    awsPattern := `(aws|cloudfront|s3|amazon)`
	for _, url := range scraper.discoveredUrls {
		if awsMatch, _ := regexp.MatchString(awsPattern, url.url); awsMatch {
			scraper.logger.Debug("AWS: ", url.url)
			url.aws = true
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
    scraper.logger.Debug("#--> ", line)
	ipPattern := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
	submatchall := ipPattern.FindAllString(line, -1)
	for _, element := range submatchall {
		scraper.logger.Debug("--> ip: ", element)
		scraper.discoveredIps = append(scraper.discoveredIps, element)
	}
	// hostnamePattern explanation:
	// 1. in a web page, we only care about matches starting with one of "'/ (i.e http://, or "www...)
	// 2. next comes the sub domains pattern, which allows 1 or more valid subdomains (i.e. www.my-site, or just my-site)
	// 3. all hostnames must end with a period followed by valid tld characters (i.e. www.my-site.com)
	hostnamePattern := `["'/]([a-zA-Z0-9]+(-[a-zA-Z0-9]+)*\.)+[a-zA-Z]+`
	hostnameRegex := regexp.MustCompile(hostnamePattern)
	// urlPattern explanation:
	// 1. url must start with hostname
	// 2. a url can then have optional paths or options, so allow a /followed by any combination of
	// legal url characters
	// 3. the urls we care about must target some file, so this pattern of legal url characters must
	// end with . followed by lower case alphabet characters to mark a file extension (i.e. .js)
	// 4. lastly, the url may have 0 or 1 options after it (i.e. ?v=3411234)
	urlPattern := hostnamePattern + `((/[a-zA-Z0-9-_&=\.%?/]*)*(\.[a-z]+)(\?[a-zA-Z0-9-_&=\.%?]*){0,1}){0,1}`
	urlRegex := regexp.MustCompile(urlPattern)
	submatchall = urlRegex.FindAllString(line, -1)
	for _, rawUrl := range submatchall {
		rawHostname := hostnameRegex.FindString(rawUrl)
		tldPattern := regexp.MustCompile(`\.([a-zA-Z]+(-[a-zA-Z]+)*)$`)
		tldMatch := tldPattern.FindString(rawHostname)
		if _, ok := tlds.TLDS[tldMatch[1:]]; ok {
			// remove leading /'"
			var url string
			var hostname string
			if rawUrl[0] == '"' || rawUrl[0] == '\'' || rawUrl[0] == '/' {
				url = rawUrl[1:]
				hostname = rawHostname[1:]
			} else {
				url = rawUrl
				hostname = rawHostname
			}
			scraper.logger.Info("---> url: ", url)
			scraper.logger.Info("---> hostname: ", hostname)
			newUrl := DiscoveredUrl{url: url, hostname: hostname}
			scraper.discoveredUrls = append(scraper.discoveredUrls, &newUrl)
		}
	}
}
