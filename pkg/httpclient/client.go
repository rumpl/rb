package httpclient

import (
	"fmt"
	"maps"
	"net/http"
	"runtime"

	"github.com/rumpl/rb/pkg/version"
)

type HTTPOptions struct {
	Header http.Header
}

type Opt func(*HTTPOptions)

func NewHTTPClient(opts ...Opt) *http.Client {
	httpOptions := HTTPOptions{
		Header: make(http.Header),
	}

	for _, opt := range opts {
		opt(&httpOptions)
	}

	// Enforce a consistent User-Agent header
	httpOptions.Header.Set("User-Agent", fmt.Sprintf("RB/%s (%s; %s)", version.Version, runtime.GOOS, runtime.GOARCH))

	return &http.Client{
		Transport: &userAgentTransport{
			httpOptions: httpOptions,
			rt:          http.DefaultTransport,
		},
	}
}

func WithHeader(key, value string) Opt {
	return func(o *HTTPOptions) {
		o.Header.Set(key, value)
	}
}

func WithProxiedBaseURL(value string) Opt {
	return func(o *HTTPOptions) {
		o.Header.Set("X-RB-Forward", value)

		// Enforce consistent headers (Anthropic client sets similar header already)
		o.Header.Set("X-RB-Lang", "go")
		o.Header.Set("X-RB-OS", runtime.GOOS)
		o.Header.Set("X-RB-Arch", runtime.GOARCH)
		o.Header.Set("X-RB-Runtime", "rb")
		o.Header.Set("X-RB-Runtime-Version", version.Version)
	}
}

func WithProvider(provider string) Opt {
	return func(o *HTTPOptions) {
		o.Header.Set("X-RB-Provider", provider)
	}
}

func WithModel(model string) Opt {
	return func(o *HTTPOptions) {
		o.Header.Set("X-RB-Model", model)
	}
}

type userAgentTransport struct {
	httpOptions HTTPOptions
	rt          http.RoundTripper
}

func (u *userAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	r2 := req.Clone(req.Context())
	maps.Copy(r2.Header, u.httpOptions.Header)
	return u.rt.RoundTrip(r2)
}
