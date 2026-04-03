package mlbapi

import (
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
)

func newTestClient(t *testing.T, handler func(*http.Request) (int, string)) *Client {
	t.Helper()

	return NewClient(
		WithBaseURL("https://example.test/api"),
		WithHTTPClient(&http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				t.Helper()
				status, body := handler(r)
				if status == 0 {
					status = http.StatusOK
				}

				return &http.Response{
					StatusCode: status,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body:       io.NopCloser(strings.NewReader(body)),
				}, nil
			}),
		}),
	)
}

type statusRecorderTransport struct {
	base     http.RoundTripper
	mu       sync.Mutex
	statuses []int
}

func (s *statusRecorderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	transport := s.base
	if transport == nil {
		transport = http.DefaultTransport
	}

	resp, err := transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	s.statuses = append(s.statuses, resp.StatusCode)
	s.mu.Unlock()

	return resp, nil
}

func (s *statusRecorderTransport) Statuses() []int {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]int, len(s.statuses))
	copy(out, s.statuses)
	return out
}
