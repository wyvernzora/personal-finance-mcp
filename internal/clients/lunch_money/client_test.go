package lunchmoney

import "net/http"

// fakeTransport lets us stub out HTTP responses.
type fakeTransport struct {
	fn func(req *http.Request) (*http.Response, error)
}

func (f *fakeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return f.fn(req)
}

// newTestClient returns a *client with its HTTP transport replaced by fn.
func newTestClient(token string, fn func(req *http.Request) (*http.Response, error)) *client {
	cli := NewClient(token).(*client)
	cli.Client = &http.Client{Transport: &fakeTransport{fn}}
	return cli
}
