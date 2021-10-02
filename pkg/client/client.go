package client

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	defaultScheme      = "http"
	searchPath         = "/booksearch"
	downloadPath       = "/b/"
	browserUserAgent   = "Safari: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36"
	defaultProxyScheme = "http"
	defaultProxyUrl    = "http://localhost:8118"
)

type FlibustaClient struct {
	httpClient *http.Client
	proxyUrl   *url.URL
}

type DownloadResult struct {
	Name string
	File []byte
}

type InfoResult struct {
	ID         string
	Title      string
	Genre      string
	Annotation string
	Size       string
	Formats    []string
}

func validateBookFormat(format string) (err error) {

	for _, valid := range validFormats {
		if format == valid {
			return nil
		}
	}
	return errors.New("invalid book format")
}

func isHttpProxy(url string) bool {
	return strings.HasPrefix(url, defaultProxyScheme)
}

func FromEnv() (*FlibustaClient, error) {
	proxyUrlString := os.Getenv("FLIBUSTA_PROXY_URL")
	if proxyUrlString == "" {
		proxyUrlString = defaultProxyUrl
	}
	if !isHttpProxy(proxyUrlString) {
		return nil, fmt.Errorf("%s does not contain defaultScheme (http or https)", proxyUrlString)
	}

	proxyUrl, err := url.Parse(proxyUrlString)
	if err != nil {
		err = errors.New("invalid FLIBUSTA_PROXY_URL")
		return nil, err
	}

	client := FlibustaClient{}
	myClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	client.httpClient = myClient
	client.proxyUrl = proxyUrl

	return &client, nil
}

type ResponseResult struct {
	Host     string
	Response *http.Response
	Error    error
}

// Fetch all known mirrors and return first response
func executeRequest(client *http.Client, path string, headers Headers) (*http.Response, error) {
	mirrors := FlibustaMirrors
	envHost := getEnvHost()
	if envHost != "" {
		mirrors = append(mirrors, envHost)
	}
	result := make(chan *ResponseResult)
	for _, host := range mirrors {
		req, err := buildRequest(host, path, headers)
		if err != nil {
			continue
		}
		go func(r *http.Request, h string, out chan *ResponseResult) {
			resp, err := client.Do(r)
			out <- &ResponseResult{
				Host:     h,
				Response: resp,
				Error:    err,
			}
		}(req, host, result)
	}
	for i := 0; i < len(mirrors); i++ {
		rr := <-result
		if rr.Error != nil {
			log.Println(rr.Error)
		} else if rr.Response.StatusCode != 200 {
			// TODO: should handle this?
			bodyBytes, _ := io.ReadAll(rr.Response.Body)
			body := string(bodyBytes)
			log.Println(body)
			defer rr.Response.Body.Close()
		} else {
			return rr.Response, nil
		}
	}
	return nil, errors.New("All request attempts failed")
}

func (c *FlibustaClient) Search(searchQuery string, respProcessor func(stream io.Reader) (*[]ListItem, error)) (result *[]ListItem, err error) {
	searchPath := buildSearchPath(searchQuery)
	headers := getHeaders()
	log.Printf("Search Flibusta for `%s`", searchPath)

	resp, err := executeRequest(c.httpClient, searchPath, headers)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	return respProcessor(resp.Body)
}

func (c *FlibustaClient) Download(id string, bookFormat string) (result *DownloadResult, err error) {
	err = validateBookFormat(bookFormat)
	if err != nil {
		return
	}
	bookPath := buildDownloadPath(id, bookFormat)
	headers := getHeaders()

	log.Printf("Download file by id: `%s`", bookPath)

	resp, err := executeRequest(c.httpClient, bookPath, headers)

	if err != nil {
		return
	}
	defer resp.Body.Close()

	file, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return &DownloadResult{Name: getFileNameFromHeader(&resp.Header), File: file}, nil
}

func (c *FlibustaClient) Info(id string, respProcessor func(stream io.Reader) (result *InfoResult, err error)) (result *InfoResult, err error) {
	infoPath := buildInfoPath(id)
	headers := getHeaders()

	log.Printf("Download file by id: `%s`", infoPath)

	resp, err := executeRequest(c.httpClient, infoPath, headers)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	return respProcessor(resp.Body)
}
