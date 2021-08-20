package client

import (
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

const defaultHost = "flibustahezeous3.onion"

func getHost() (host string) {
	host, exists := os.LookupEnv("FLIBUSTA_HOST")
	if !exists {
		return defaultHost
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
	q.Set("chb", "on") // Search only books
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
