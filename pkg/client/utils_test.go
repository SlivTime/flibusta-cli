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

func Test_buildSearchPath(t *testing.T) {
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
			"http://flibusta/booksearch?ask=&chb=on",
		},
		{
			"Empty",
			args{"book"},
			"http://flibusta/booksearch?ask=book&chb=on",
		},
		{
			"Empty",
			args{"my book"},
			"http://flibusta/booksearch?ask=my+book&chb=on",
		},
		{
			"Empty",
			args{"The book#that^shoud%be&escaped"},
			"http://flibusta/booksearch?ask=The+book%23that%5Eshoud%25be%26escaped&chb=on",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSearchUrl := buildSearchUrl(tt.args.searchQuery)
			if gotSearchUrl.String() != tt.wantSearchUrl {
				t.Errorf("buildSearchPath() gotSearchUrl = %v, want %v", gotSearchUrl, tt.wantSearchUrl)
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
			if gotU := getEnvUrl(); !reflect.DeepEqual(gotU, tt.wantU) {
				t.Errorf("getEnvUrl() = %v, want %v", gotU, tt.wantU)
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
			want: "http://flibusta/b",
		},
		{
			name: "numbers",
			args: args{
				"1",
				"1",
			},
			want: "http://flibusta/b/1/1",
		},
		{
			name: "most common",
			args: args{
				"123",
				"mobi",
			},
			want: "http://flibusta/b/123/mobi",
		},
		{
			name: "Foobar",
			args: args{
				"foo",
				"bar",
			},
			want: "http://flibusta/b/foo/bar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildDownloadUrl(tt.args.bookId, tt.args.bookFormat); got.String() != tt.want {
				t.Errorf("buildDownloadPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildRequest(t *testing.T) {
	type args struct {
		host    string
		url     *url.URL
		headers Headers
	}
	type want struct {
		userAgent string
		method    string
		url       string
	}
	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{
			"Foo",
			args{
				"",
				&url.URL{Host: "foo"},
				getHeaders(),
			},
			want{
				browserUserAgent,
				"GET",
				"foo",
			},
			true,
		},
		{
			"Full host",
			args{
				"example.com",
				&url.URL{},
				getHeaders(),
			},
			want{
				browserUserAgent,
				"GET",
				"http://example.com",
			},
			false,
		},
		{
			"Keep protocol ",
			args{
				"https://example.com/",
				&url.URL{},
				getHeaders(),
			},
			want{
				browserUserAgent,
				"GET",
				"https://example.com",
			},
			false,
		},
		{
			"Full url",
			args{
				"flibustahezeous3.onion",
				&url.URL{Path: "b/175105"},
				getHeaders(),
			},
			want{
				browserUserAgent,
				"GET",
				"http://flibustahezeous3.onion/b/175105",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildRequest(tt.args.host, tt.args.url, tt.args.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if got.URL.String() != tt.want.url {
					t.Errorf("Got %v, want %v", got.URL, tt.want.url)
				}
				if got.UserAgent() != tt.want.userAgent {
					t.Errorf("Got %v, want %v", got.UserAgent(), tt.want.userAgent)
				}
				if got.Method != tt.want.method {
					t.Errorf("Got %v, want %v", got.Method, tt.want.method)
				}
			}
		})
	}
}

func Test_buildInfoUrl(t *testing.T) {
	type args struct {
		bookId string
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
			},
			want: "http://flibusta/b",
		},
		{
			name: "numbers",
			args: args{
				"1",
			},
			want: "http://flibusta/b/1",
		},
		{
			name: "most common",
			args: args{
				"123",
			},
			want: "http://flibusta/b/123",
		},
		{
			name: "Foobar",
			args: args{
				"foo",
			},
			want: "http://flibusta/b/foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildInfoUrl(tt.args.bookId).String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildInfoUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
