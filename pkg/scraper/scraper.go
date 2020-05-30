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

func ScrapeLocalFile(hostname string, localFilename string, verbose bool) error {
    logger := logging.NewLogger(verbose)
    logger.Debug("Scraping local file...")
    s := Scraper{logger: logger, rootHostname: hostname}
    s.localFilename = localFilename
    s.scrapeLocalFile()
    s.markUrlsToCheck()
    for _, el := range s.discoveredUrls {

        if el.follow {
            logger.Debug("Will check: ", el, " with: ", el.url)
        }
    }
    return nil
}

type DiscoveredUrl struct {
    url string
    hostname string
    follow bool
    aws bool
    awsService string
}

type Scraper struct {
    logger *zap.SugaredLogger
    rootHostname string
    localFilename string
    targetSiteUrl string
    targetSiteListFilename string

    discoveredUrls []*DiscoveredUrl
    discoveredIps []string
}

func (scraper *Scraper) markUrlsAsAwsService() error {
    return nil
}

func (scraper *Scraper) markUrlsToCheck() error {
    jsExtensionPattern := `(cloudflare|cloudfront).*\.js(\?){0,1}`
	jsRegex := regexp.MustCompile(jsExtensionPattern)

	for _, url := range scraper.discoveredUrls {
        regexMatch := jsRegex.MatchString(url.url)
        rootHostnameMatch, _ := regexp.MatchString(`(` + url.hostname + `).*\.js(\?){0,1}`, url.url)
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
    // hostnamePattern explanation:
    // 1. in a web page, we only care about matches starting with one of "'/ (i.e http://, or "www...)
    // 2. next comes the sub domains pattern, which allows 1 or more valid subdomains (i.e. www.my-site, or just my-site)
    // 3. all hostnames must end with a period followed by valid tld characters (i.e. www.my-site.com)
    hostnamePattern := `["'/]([a-zA-Z0-9]+(-[a-zA-Z0-9]+)*)+(\.([a-zA-Z]+(-[a-zA-Z]+)*))+`
	hostnameRegex := regexp.MustCompile(hostnamePattern)
    // urlPattern explanation:
    // 1. url must start with hostname
    // 2. a url can then have optional paths or options, so allow a /followed by any combination of
    // legal url characters
    // 3. the urls we care about must target some file, so this pattern of legal url characters must
    // end with . followed by lower case alphabet characters to mark a file extension (i.e. .js)
    // 4. lastly, the url may have 0 or 1 options after it (i.e. ?v=3411234)
    urlPattern := hostnamePattern + `(/[a-zA-Z0-9-_&=\.%?/]*)*(\.[a-z]+)(\?[a-zA-Z0-9-_&=\.%?]*){0,1}`
	urlRegex := regexp.MustCompile(urlPattern)
    submatchall = urlRegex.FindAllString(line, -1)
	for _, raw_url := range submatchall {
        hostname := hostnameRegex.FindString(raw_url)
        tldPattern := regexp.MustCompile(`\.([a-zA-Z]+(-[a-zA-Z]+)*)$`)
        tldMatch := tldPattern.FindString(hostname)
        if _, ok := tlds.TLDS[tldMatch[1:]]; ok {
            // remove leading /'"
            var url string
            if raw_url[0] == '"' || raw_url[0] == '\'' || raw_url[0] == '/' {
                url = raw_url[1:]
            } else {
                url = raw_url
            }
            scraper.logger.Debug("---> url: ", url)
            scraper.logger.Debug("--> hostname: ", hostname)
            newUrl := DiscoveredUrl{url: url, hostname: hostname}
            scraper.discoveredUrls = append(scraper.discoveredUrls, &newUrl)
        }
	}
}
