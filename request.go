package plow

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/google/go-cmp/cmp"
	"github.com/onsi/gomega"
)

func newClient() (*http.Client, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	if os.Getenv("HTTP_PROXY") != "" {
		proxyUrl, err := url.Parse(os.Getenv("HTTP_PROXY"))
		if err != nil {
			return nil, err
		}
		transport.Proxy = http.ProxyURL(proxyUrl)
	}

	return &http.Client{
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}, nil
}

type BasicAuth struct {
	User string
	Pass string
}

type RequestOptions struct {
	Headers   http.Header
	Cookies   map[string]string
	BasicAuth *BasicAuth
	Body      io.Reader
}

func newRequest(method string, uri string, options RequestOptions) *http.Request {
	req, err := http.NewRequest(method, uri, options.Body)
	if err != nil {
		panic(err)
	}

	if options.Headers != nil {
		if hostHeader, ok := options.Headers["Host"]; ok && len(options.Headers) > 0 {
			req.Host = hostHeader[0]
			delete(options.Headers, "Host")
		}
		req.Header = options.Headers
	}

	if options.BasicAuth != nil {
		req.SetBasicAuth(options.BasicAuth.User, options.BasicAuth.Pass)
	}

	if options.Cookies != nil {
		for k, v := range options.Cookies {
			req.AddCookie(&http.Cookie{Name: k, Value: v})
		}
	}

	return req
}

// SendGETRequest send http request
func SendGETRequest(uri string, options RequestOptions) (*http.Response, error) {
	cl, err := newClient()
	if err != nil {
		panic(err)
	}
	req := newRequest("GET", uri, options)
	return cl.Do(req)
}

// SendPOSTRequest send http request
func SendPOSTRequest(uri string, options RequestOptions) (*http.Response, error) {
	cl, err := newClient()
	if err != nil {
		panic(err)
	}
	req := newRequest("POST", uri, options)
	return cl.Do(req)
}

func ExpectJSON(resp *http.Response, err error, expectation, target interface{}) {
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(resp).ToNot(gomega.BeNil())
	jsonErr := json.NewDecoder(resp.Body).Decode(target)
	gomega.Expect(jsonErr).To(gomega.BeNil())
	gomega.Expect(cmp.Diff(expectation, target)).To(gomega.BeEmpty())
}
