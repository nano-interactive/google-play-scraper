package google_play_scraper

import "net/http"

type (
	HTTPClient interface {
		Do(req *http.Request) (*http.Response, error)
	}
)
