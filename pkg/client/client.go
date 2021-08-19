package client

import (
	"errors"
	"flibusta-go/internal/env"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
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
	env.Load()
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

func (c *FlibustaClient) Search(searchQuery string) (result string, err error) {
	searchUrl, err := buildSearchUrl(searchQuery)
	if err != nil {
		return
	}
	req, err := buildRequest(searchUrl)
	if err != nil {
		return
	}

	log.Printf("Search Flibusta for `%s`", searchUrl)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	// TODO: return parsed struct
	return string(body), nil
}

func (c *FlibustaClient) Download(id string, bookFormat string) (result *DownloadResult, err error) {
	err = validateBookFormat(bookFormat)
	if err != nil {
		return
	}
	bookUrl, err := buildDownloadUrl(id, bookFormat)
	if err != nil {
		return
	}
	req, err := buildRequest(bookUrl)
	if err != nil {
		return
	}

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

	result = &DownloadResult{Name: getFileNameFromHeader(resp.Header), File: file}
	return result, nil
}

func getHost() (host string) {
	host = os.Getenv("FLIBUSTA_HOST")
	if host == "" {
		log.Fatal("Missing FLIBUSTA_HOST in environment")
	}
	return host
}

func getBaseUrl() (u *url.URL) {
	return &url.URL{
		Host:   getHost(),
		Scheme: scheme,
	}
}

func buildSearchUrl(searchQuery string) (searchUrl string, err error) {
	u := getBaseUrl()
	u.Path = searchPath
	q := u.Query()
	q.Set("ask", searchQuery)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func buildDownloadUrl(bookId string, bookFormat string) (bookUrl string, err error) {
	u := getBaseUrl()
	u.Path = path.Join(downloadPath, bookId, bookFormat)
	return u.String(), nil
}

func buildRequest(url string) (request *http.Request, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", browserUserAgent)
	return req, nil
}

func getFileNameFromHeader(h http.Header) string {
	disposition := h.Get("Content-Disposition")
	if disposition == "" {
		return ""
	}
	splitted := strings.Split(disposition, "filename=")
	if len(splitted) > 1 {
		return splitted[1]
	} else {
		return ""
	}
}
