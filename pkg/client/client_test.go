package client

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	_ "net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"testing"
)

var (
	successSearchTestResult = &[]ListItem{
		{
			Title:   "Ok",
			Authors: nil,
			ID:      "1",
		},
	}
	testServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, testResponse(r.URL))
	}))
)

func TestFlibustaClient_Download(t *testing.T) {
	tsURL, _ := url.Parse(testServer.URL)
	oldEnv := os.Getenv("FLIBUSTA_HOST")
	defer func() {
		_ = os.Setenv("FLIBUSTA_HOST", oldEnv)
	}()
	type env struct {
		envHost string
	}
	type args struct {
		id         string
		bookFormat string
	}
	tests := []struct {
		name       string
		env        env
		args       args
		wantResult *DownloadResult
		wantErr    bool
	}{
		{
			"Success story",
			env{
				tsURL.Host,
			},
			args{
				"123",
				"mobi",
			},

			&DownloadResult{
				Name: "",
				File: []byte("/b/123/mobi\n"),
			},
			false,
		},
		{
			"Wrong format",
			env{
				tsURL.Host,
			},
			args{
				"123",
				"docx",
			},

			nil,
			true,
		},
		{
			"Missing host - error",
			env{
				"missing.host",
			},
			args{
				"123",
				"mobi",
			},

			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &FlibustaClient{
				httpClient: testServer.Client(),
			}
			_ = os.Setenv("FLIBUSTA_HOST", tt.env.envHost)
			gotResult, err := c.Download(tt.args.id, tt.args.bookFormat)
			if (err != nil) != tt.wantErr {
				t.Errorf("Download() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Download() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestFlibustaClient_Search(t *testing.T) {
	tsURL, _ := url.Parse(testServer.URL)
	oldEnv := os.Getenv("FLIBUSTA_HOST")
	defer func() {
		_ = os.Setenv("FLIBUSTA_HOST", oldEnv)
	}()

	type env struct {
		envHost string
	}
	type args struct {
		searchQuery string
	}
	type want struct {
		result *[]ListItem
		url    string
	}
	tests := []struct {
		name    string
		env     env
		args    args
		want    want
		wantErr bool
	}{
		{
			"Success story",
			env{
				tsURL.Host,
			},
			args{"test"},

			want{
				successSearchTestResult,
				"/booksearch?ask=test&chb=on\n",
			},
			false,
		},
		{
			"Missing host - error",
			env{
				"missing.host",
			},
			args{"test"},
			want{
				nil,
				"/booksearch?ask=test&chb=on\n",
			},
			true,
		},
		{
			"Error from resp processor",
			env{
				tsURL.Host,
			},
			args{"test"},
			want{
				nil,
				"/does_not_match",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &FlibustaClient{
				httpClient: testServer.Client(),
			}
			_ = os.Setenv("FLIBUSTA_HOST", tt.env.envHost)
			gotResult, err := c.Search(tt.args.searchQuery, processorFuncFabric(tt.want.url))
			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.want.result) {
				t.Errorf("Search() gotResult = %v, want %v", gotResult, tt.want.result)
			}
		})
	}
}

func TestFromEnv(t *testing.T) {
	oldEnv := os.Getenv("FLIBUSTA_PROXY_URL")
	defer func() {
		_ = os.Setenv("FLIBUSTA_PROXY_URL", oldEnv)
	}()
	type env struct {
		proxyUrl string
	}
	tests := []struct {
		name    string
		env     env
		want    *FlibustaClient
		wantErr bool
	}{
		{
			"Empty env - use default",
			env{},
			&FlibustaClient{
				&http.Client{},
				&url.URL{
					Scheme: "http",
					Host:   "localhost:8118",
				},
			},
			false,
		},
		{
			"Invalid proxy url in env",
			env{
				"postgres://user:abc{DEf1=ghi@example.com:5432/db?sslmode=require",
			},
			nil,
			true,
		},
		{
			"Invalid http proxy url in env",
			env{
				"http://user:abc{DEf1=ghi@example.com:5432/db?sslmode=require",
			},
			nil,
			true,
		},
		{
			"Proxy from env",
			env{
				"http://test.proxy:123/",
			},
			&FlibustaClient{
				&http.Client{},
				&url.URL{
					Scheme: "http",
					Host:   "test.proxy:123",
					Path:   "/",
				},
			},
			false,
		},
		{
			"Proxy from env without slash",
			env{
				"http://test.proxy:123",
			},
			&FlibustaClient{
				&http.Client{},
				&url.URL{
					Scheme: "http",
					Host:   "test.proxy:123",
				},
			},
			false,
		},
		{
			"Proxy from env without scheme",
			env{
				"proxy.com:123",
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		_ = os.Setenv("FLIBUSTA_PROXY_URL", tt.env.proxyUrl)
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromEnv()
			if (err != nil) != tt.wantErr {
				t.Errorf("FromEnv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if !reflect.DeepEqual(got.proxyUrl, tt.want.proxyUrl) {
					t.Errorf("FromEnv() \ngot:  %v\nwant: %v", got.proxyUrl, tt.want.proxyUrl)
				}
			}
		})
	}
}

func Test_validateBookFormat(t *testing.T) {
	type args struct {
		format string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Empty",
			args{""},
			true,
		},
		{
			"Epub",
			args{"epub"},
			false,
		},
		{
			"Fb2",
			args{"fb2"},
			false,
		},
		{
			"Mobi",
			args{"mobi"},
			false,
		},
		{
			"pdf",
			args{"pdf"},
			true,
		},
		{
			"djvu",
			args{"djvu"},
			true,
		},
		{
			"txt",
			args{"txt"},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateBookFormat(tt.args.format); (err != nil) != tt.wantErr {
				t.Errorf("validateBookFormat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func testResponse(url *url.URL) string {
	return url.String()
}

func processorFuncFabric(wantUrl string) func(stream io.Reader) (*[]ListItem, error) {
	return func(stream io.Reader) (*[]ListItem, error) {
		bodyBytes, _ := io.ReadAll(stream)
		body := string(bodyBytes)
		if body != wantUrl {
			return nil, errors.New("fail")
		}
		return &[]ListItem{
			{
				ID:      "1",
				Title:   "Ok",
				Authors: nil,
			},
		}, nil
	}
}
