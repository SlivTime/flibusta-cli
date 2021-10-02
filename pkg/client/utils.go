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

func getEnvUrl() (u *url.URL) {
	return &url.URL{
		Host:   getEnvHost(),
		Scheme: defaultScheme,
	}
}

func getBaseUrl() (u *url.URL) {
	return &url.URL{
		Host:   "flibusta",
		Scheme: defaultScheme,
	}
}

func buildSearchUrl(searchQuery string) *url.URL {
	u := getBaseUrl()
	u.Path = searchPath
	q := u.Query()
	q.Set("ask", searchQuery)
	q.Set("chb", "on") // Search only books
	u.RawQuery = q.Encode()
	return u
}

func buildDownloadUrl(bookId string, bookFormat string) *url.URL {
	u := getBaseUrl()
	u.Path = path.Join(downloadPath, bookId, bookFormat)
	return u
}

func buildInfoUrl(bookId string) *url.URL {
	u := getBaseUrl()
	u.Path = path.Join(downloadPath, bookId)
	return u
}

func buildRequest(host string, url *url.URL, headers Headers) (*http.Request, error) {
	match := HostRe.FindStringSubmatch(host)
	if match == nil {
		return nil, errors.New("Cannot parse host")
	}
	scheme := match[1]
	if scheme == "" {
		scheme = defaultScheme
	}
	cleanHost := match[3]

	url.Scheme = scheme
	url.Host = cleanHost
	req, _ := http.NewRequest("GET", url.String(), nil)
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
