# Go Scraper

Scrapes hostnames from website HTML and JS and checks for AWS-related hostnames.

Requires:
 - go1.14.3

Usage:
 - `./scraper scrape --help`

### Known issues

To handle large minified js files where the entire file is a single line, each source file is read in predetermined chunk sizes instead of reading line by line. This means there is a chance that a url is split between chunk n and chunk n+1. The larger the chunk size the lower the probability of this occurring.
