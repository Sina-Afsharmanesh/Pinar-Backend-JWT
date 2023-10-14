package httpLocal

import (
	"net/http"
	"time"
)

var (
	HttpClient = createHttpClient()
)

func createHttpClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        0,
			MaxIdleConnsPerHost: 10000,
			ForceAttemptHTTP2:   true,
			DisableKeepAlives:   false,
		},
	}
}
