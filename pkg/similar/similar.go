package similar

import (
	"github.com/nano-interactive/google-play-scraper"
	"github.com/nano-interactive/google-play-scraper/pkg/scraper"
)

// Options type alias
type Options = scraper.Options

// New return similar list instance
func New(appID string, options Options) *scraper.Scraper {
	a := google_play_scraper.New(appID, google_play_scraper.Options{
		Country:  options.Country,
		Language: options.Language,
	})
	err := a.LoadDetails()
	if err != nil || a.SimilarURL == "" {
		return nil
	}
	return scraper.New(a.SimilarURL, &options)
}
