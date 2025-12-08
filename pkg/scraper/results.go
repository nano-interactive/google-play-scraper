package scraper

import (
	"github.com/nano-interactive/google-play-scraper"
)

// Results of operation
type Results []*google_play_scraper.App

// Append result
func (results *Results) Append(res ...google_play_scraper.App) {
	for _, result := range res {
		if !results.searchDuplicate(result.ID) {
			results.append(result)
		}
	}
}

func (results *Results) append(result google_play_scraper.App) {
	*results = append(*results, &result)
}

func (results *Results) searchDuplicate(id string) bool {
	for _, result := range *results {
		if id == result.ID {
			return true
		}
	}
	return false
}
