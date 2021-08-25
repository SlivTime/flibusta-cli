package client

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

const (
	defaultHost        = "flibustahezeous3.onion"
	FlibustaHostEnvKey = "FLIBUSTA_HOST"
)

func getHost() (host string) {
	host = os.Getenv(FlibustaHostEnvKey)
	if host == "" {
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

func buildSearchUrl(searchQuery string) string {
	u := getBaseUrl()
	u.Path = searchPath
	q := u.Query()
	q.Set("ask", searchQuery)
	q.Set("chb", "on") // Search only books
	u.RawQuery = q.Encode()
	return u.String()
}

func buildDownloadUrl(bookId string, bookFormat string) string {
	u := getBaseUrl()
	u.Path = path.Join(downloadPath, bookId, bookFormat)
	return u.String()
}

func buildInfoUrl(bookId string) string {
	u := getBaseUrl()
	u.Path = path.Join(downloadPath, bookId)
	return u.String()
}

func buildRequest(url string) *http.Request {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", browserUserAgent)
	return req
}

func getFileNameFromHeader(h *http.Header) string {
	disposition := h.Get("Content-Disposition")
	if disposition == "" {
		return ""
	}
	splitted := strings.Split(disposition, "filename=")
	if len(splitted) > 1 {
		return strings.ReplaceAll(splitted[1], "\"", "")
	} else {
		return ""
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
