package client

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

const (
	scheme           = "http"
	searchPath       = "/booksearch"
	downloadPath     = "/b/"
	browserUserAgent = "Safari: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36"
	Fb2              = "fb2"
	Epub             = "epub"
	Mobi             = "mobi"
)

type FlibustaClient struct {
	httpClient *http.Client
	host       string
}

type DownloadResult struct {
	Name string
	File []byte
}

func validateBookFormat(format string) (err error) {
	validFormats := []string{Fb2, Epub, Mobi}
	for _, valid := range validFormats {
		if format == valid {
			return nil
		}
	}
	return errors.New("invalid book format")
}

func FromEnv() (*FlibustaClient, error) {
	proxyUrlString := os.Getenv("FLIBUSTA_PROXY_URL")

	proxyUrl, err := url.Parse(proxyUrlString)
	if err != nil {
		log.Fatal("Invalid FLIBUSTA_PROXY_URL")
	}

	client := FlibustaClient{}
	myClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	client.httpClient = myClient

	return &client, nil
}

func (c *FlibustaClient) Search(searchQuery string) (result *[]ListItem, err error) {
	searchUrl := buildSearchUrl(searchQuery)
	req := buildRequest(searchUrl)

	log.Printf("Search Flibusta for `%s`", searchUrl)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	result = ParseSearch(resp.Body)

	return result, nil
}

func (c *FlibustaClient) Download(id string, bookFormat string) (result *DownloadResult, err error) {
	err = validateBookFormat(bookFormat)
	if err != nil {
		return
	}
	bookUrl := buildDownloadUrl(id, bookFormat)
	req := buildRequest(bookUrl)

	log.Printf("Download file by id: `%s`", bookUrl)

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return
	}
	defer resp.Body.Close()

	file, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	result = &DownloadResult{Name: getFileNameFromHeader(&resp.Header), File: file}
	return result, nil
}
