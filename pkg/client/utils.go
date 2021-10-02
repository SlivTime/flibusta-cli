package client

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
)

var HostRe = regexp.MustCompile(`(?P<Scheme>https?)?(://)?(?P<Host>[0-8a-z.]+):?(?P<Port>[0-9]+)?/?`)

type Headers map[string]string

func getEnvHost() string {
	return os.Getenv(FlibustaHostEnvKey)
}

func getBaseUrl() (u *url.URL) {
	return &url.URL{
		Host:   getEnvHost(),
		Scheme: defaultScheme,
	}
}

func buildSearchPath(searchQuery string) string {
	u := url.URL{}
	u.Path = searchPath
	q := u.Query()
	q.Set("ask", searchQuery)
	q.Set("chb", "on") // Search only books
	u.RawQuery = q.Encode()
	return u.String()
}

func buildDownloadPath(bookId string, bookFormat string) string {
	u := url.URL{}
	u.Path = path.Join(downloadPath, bookId, bookFormat)
	return u.String()
}

func buildInfoPath(bookId string) string {
	u := url.URL{}
	u.Path = path.Join(downloadPath, bookId)
	return u.String()
}

func buildRequest(host string, path string, headers Headers) (*http.Request, error) {
	match := HostRe.FindStringSubmatch(host)
	if match == nil {
		return nil, errors.New("Cannot parse host")
	}
	scheme := match[1]
	if scheme == "" {
		scheme = defaultScheme
	}

	u := url.URL{
		Scheme: scheme,
		Host:   match[3],
		Path:   path,
	}
	req, _ := http.NewRequest("GET", u.String(), nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return req, nil
}

func getHeaders() Headers {
	headers := Headers{
		"User-Agent": browserUserAgent,
	}
	return headers
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
