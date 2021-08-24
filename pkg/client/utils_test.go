package client

import (
	"net/http"
	"net/url"
	"os"
	"reflect"
	"testing"
)

const (
	testHost = "test.host"
)

func Test_buildSearchUrl(t *testing.T) {
	_ = os.Setenv(FlibustaHostEnvKey, testHost)
	type args struct {
		searchQuery string
	}
	tests := []struct {
		name          string
		args          args
		wantSearchUrl string
	}{
		{
			"Empty",
			args{""},
			"http://test.host/booksearch?ask=&chb=on",
		},
		{
			"Empty",
			args{"book"},
			"http://test.host/booksearch?ask=book&chb=on",
		},
		{
			"Empty",
			args{"my book"},
			"http://test.host/booksearch?ask=my+book&chb=on",
		},
		{
			"Empty",
			args{"The book#that^shoud%be&escaped"},
			"http://test.host/booksearch?ask=The+book%23that%5Eshoud%25be%26escaped&chb=on",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSearchUrl := buildSearchUrl(tt.args.searchQuery)
			if gotSearchUrl != tt.wantSearchUrl {
				t.Errorf("buildSearchUrl() gotSearchUrl = %v, want %v", gotSearchUrl, tt.wantSearchUrl)
			}
		})
	}
}

func Test_getBaseUrl(t *testing.T) {
	_ = os.Setenv(FlibustaHostEnvKey, testHost)
	tests := []struct {
		name  string
		wantU *url.URL
	}{
		{
			"Just add protocol",
			&url.URL{
				Scheme: "http",
				Host:   testHost,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotU := getBaseUrl(); !reflect.DeepEqual(gotU, tt.wantU) {
				t.Errorf("getBaseUrl() = %v, want %v", gotU, tt.wantU)
			}
		})
	}
}

func Test_getFileNameFromHeader(t *testing.T) {
	type args struct {
		h http.Header
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Empty",
			args{
				http.Header{"Content-Type": {"text/plain"}},
			},
			"",
		},
		{
			"Empty disposition",
			args{
				http.Header{
					"Content-Type":        {"text/plain"},
					"Content-Disposition": {},
				},
			},
			"",
		},
		{
			"Form data",
			args{
				http.Header{
					"Content-Type":        {"text/plain"},
					"Content-Disposition": {"form-data"},
				},
			},
			"",
		},
		{
			"Form data, field name",
			args{
				http.Header{
					"Content-Type":        {"text/plain"},
					"Content-Disposition": {"form-data; name=\"fieldName\""},
				},
			},
			"",
		},
		{
			"Form date with file unquoted",
			args{
				http.Header{
					"Content-Type":        {"text/plain"},
					"Content-Disposition": {"form-data; name=\"fieldName\"; filename=filename.jpg"},
				},
			},
			"filename.jpg",
		},
		{
			"Form date with file",
			args{
				http.Header{
					"Content-Type":        {"text/plain"},
					"Content-Disposition": {"form-data; name=\"fieldName\"; filename=\"filename.jpg\""},
				},
			},
			"filename.jpg",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFileNameFromHeader(&tt.args.h); got != tt.want {
				t.Errorf("getFileNameFromHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getHost(t *testing.T) {
	tests := []struct {
		name     string
		envHost  string
		wantHost string
	}{
		{
			"Empty env, return default",
			"",
			defaultHost,
		},
		{
			"Default case, take from env",
			"foo.bar",
			"foo.bar",
		},
		{
			"Default case",
			"flibustahezeous3.onion",
			"flibustahezeous3.onion",
		},
		{
			"With scheme",
			"http://example.com",
			"http://example.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envHost == "" {
				_ = os.Unsetenv(FlibustaHostEnvKey)
			} else {
				_ = os.Setenv(FlibustaHostEnvKey, tt.envHost)
			}
			if gotHost := getHost(); gotHost != tt.wantHost {
				t.Errorf("getHost() = %v, want %v", gotHost, tt.wantHost)
			}
		})
	}
}

func Test_buildDownloadUrl(t *testing.T) {
	_ = os.Setenv(FlibustaHostEnvKey, testHost)
	type args struct {
		bookId     string
		bookFormat string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			args: args{
				"",
				"",
			},
			want: "http://test.host/b",
		},
		{
			name: "numbers",
			args: args{
				"1",
				"1",
			},
			want: "http://test.host/b/1/1",
		},
		{
			name: "most common",
			args: args{
				"123",
				"mobi",
			},
			want: "http://test.host/b/123/mobi",
		},
		{
			name: "Foobar",
			args: args{
				"foo",
				"bar",
			},
			want: "http://test.host/b/foo/bar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildDownloadUrl(tt.args.bookId, tt.args.bookFormat); got != tt.want {
				t.Errorf("buildDownloadUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildRequest(t *testing.T) {
	type args struct {
		url string
	}
	type want struct {
		userAgent string
		method    string
		url       string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"Foo",
			args{
				"foo",
			},
			want{
				browserUserAgent,
				"GET",
				"foo",
			},
		},
		{
			"Full host",
			args{
				"http://example.com/",
			},
			want{
				browserUserAgent,
				"GET",
				"http://example.com/",
			},
		},
		{
			"Preserve protocol",
			args{
				"https://example.com/",
			},
			want{
				browserUserAgent,
				"GET",
				"https://example.com/",
			},
		},
		{
			"Do not validate protocol, actually",
			args{
				"robert://pike.com/",
			},
			want{
				browserUserAgent,
				"GET",
				"robert://pike.com/",
			},
		},
		{
			"Full url",
			args{
				"http://flibustahezeous3.onion/b/175105",
			},
			want{
				browserUserAgent,
				"GET",
				"http://flibustahezeous3.onion/b/175105",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildRequest(tt.args.url)
			if got.URL.String() != tt.want.url {
				t.Errorf("Got %v, want %v", got.URL, tt.want.url)
			}
			if got.UserAgent() != tt.want.userAgent {
				t.Errorf("Got %v, want %v", got.UserAgent(), tt.want.userAgent)
			}
			if got.Method != tt.want.method {
				t.Errorf("Got %v, want %v", got.Method, tt.want.method)
			}
		})
	}
}
