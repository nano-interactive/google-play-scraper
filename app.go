package google_play_scraper

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"time"

	"github.com/nano-interactive/google-play-scraper/internal/parse"
	"github.com/nano-interactive/google-play-scraper/internal/util"
)

type (
	Price struct {
		Currency string
		Value    float64
	}

	App struct {
		AdSupported              bool
		AndroidVersion           string
		AndroidVersionMin        float64
		Available                bool
		ContentRating            string
		ContentRatingDescription string
		Description              string
		DescriptionHTML          string
		Developer                string
		DeveloperID              string
		DeveloperInternalID      string
		DeveloperURL             string
		DeveloperWebsite         string
		FamilyGenre              string
		FamilyGenreID            string
		Free                     bool
		Genre                    string
		GenreID                  string
		HeaderImage              string
		IAPOffers                bool
		IAPRange                 string
		Icon                     string
		ID                       string
		Installs                 string
		InstallsMin              int
		InstallsMax              int
		Permissions              map[string][]string
		Price                    Price
		PriceFull                Price
		PrivacyPolicy            string
		Ratings                  int
		RatingsHistogram         map[int]int
		RecentChanges            string
		RecentChangesHTML        string
		Released                 string
		Score                    float64
		ScoreText                string
		Screenshots              []string
		SimilarURL               string
		Summary                  string
		Title                    string
		Updated                  time.Time
		URL                      string
		Version                  string
		Video                    string
		VideoImage               string

		options *Options
		client  HTTPClient
	}

	Options struct {
		Country  string
		Language string
	}
)

const (
	detailURL = "https://play.google.com/store/apps/details?id="
	playURL   = "https://play.google.com"

	GeographicLocation = "gl"
	HostLanguage       = "hl"
)

var (
	ErrRequiredFieldsMissing = errors.New("app URL or ID must be provided")
	ErrNotFound              = errors.New("app by ID doesn't exist on app store")
)

func New(id string, options Options, client HTTPClient) *App {
	return &App{
		ID:      id,
		URL:     detailURL + id,
		options: &options,
		client:  client,
	}
}

// LoadDetails loads page and maps app details into App
func (app *App) LoadDetails() error {
	if app.URL == "" && app.ID == "" {
		return ErrRequiredFieldsMissing
	}

	if app.URL == "" {
		app.URL = detailURL + app.ID
	}

	appData, err := app.fetchAndExtractData()
	if err != nil {
		return err
	}

	app.mapResponseToApp(appData)

	return nil
}

func (app *App) mapResponseToApp(appData map[string]string) {
	if app.ID == "" {
		app.ID = parse.ID(app.URL)
	}

	for dsAppInfo := range appData {
		relativeDevURL := util.GetJSONValue(appData[dsAppInfo], "1.2.68.1.4.2")
		if relativeDevURL == "" {
			continue
		}

		app.AdSupported = util.GetJSONValue(appData[dsAppInfo], "1.2.48") != ""

		app.AndroidVersion = util.GetJSONValue(appData[dsAppInfo], "1.2.140.1.1.0.0.1", "1.2.112.141.1.1.0.0.1")
		app.AndroidVersionMin = parse.Float(app.AndroidVersion)

		app.Available = util.GetJSONValue(appData[dsAppInfo], "1.2.18.0") != ""

		app.ContentRating = util.GetJSONValue(appData[dsAppInfo], "1.2.9.0")
		app.ContentRatingDescription = util.GetJSONValue(appData[dsAppInfo], "1.2.9.2.1")

		app.DescriptionHTML = util.GetJSONValue(appData[dsAppInfo], "1.2.72.0.1")
		app.Description = util.HTMLToText(app.DescriptionHTML)

		devURL, _ := util.AbsoluteURL(playURL, relativeDevURL)
		app.Developer = util.GetJSONValue(appData[dsAppInfo], "1.2.68.0")
		app.DeveloperID = parse.ID(util.GetJSONValue(appData[dsAppInfo], "1.2.68.1.4.2"))
		app.DeveloperInternalID = util.GetJSONValue(appData[dsAppInfo], "1.2.68.2")
		app.DeveloperURL = devURL
		app.DeveloperWebsite = util.GetJSONValue(appData[dsAppInfo], "1.2.69.0.5.2")

		app.Genre = util.GetJSONValue(appData[dsAppInfo], "1.2.79.0.0.0")
		app.GenreID = util.GetJSONValue(appData[dsAppInfo], "1.2.79.0.0.2")
		app.FamilyGenre = util.GetJSONValue(appData[dsAppInfo], "1.12.13.1.0")
		app.FamilyGenreID = util.GetJSONValue(appData[dsAppInfo], "1.12.13.1.2")

		app.HeaderImage = util.GetJSONValue(appData[dsAppInfo], "1.2.96.0.3.2")

		app.IAPRange = util.GetJSONValue(appData[dsAppInfo], "1.2.19.0")
		app.IAPOffers = app.IAPRange != ""

		app.Icon = util.GetJSONValue(appData[dsAppInfo], "1.2.95.0.3.2")

		app.Installs = util.GetJSONValue(appData[dsAppInfo], "1.2.13.0")
		app.InstallsMin = parse.Int(util.GetJSONValue(appData[dsAppInfo], "1.2.13.1"))
		app.InstallsMax = parse.Int(util.GetJSONValue(appData[dsAppInfo], "1.2.13.2"))

		price := Price{
			Currency: util.GetJSONValue(appData[dsAppInfo], "1.2.57.0.0.0.0.1.0.1"),
			Value:    parse.Float(util.GetJSONValue(appData[dsAppInfo], "1.2.57.0.0.0.0.1.0.2")),
		}
		app.Free = price.Value == 0
		app.Price = price
		app.PriceFull = Price{
			Currency: util.GetJSONValue(appData[dsAppInfo], "1.2.57.0.0.0.0.1.1.1"),
			Value:    parse.Float(util.GetJSONValue(appData[dsAppInfo], "1.2.57.0.0.0.0.1.1.2")),
		}

		app.PrivacyPolicy = util.GetJSONValue(appData[dsAppInfo], "1.2.99.0.5.2")

		app.Ratings = parse.Int(util.GetJSONValue(appData[dsAppInfo], "1.2.51.2.1"))
		app.RatingsHistogram = map[int]int{
			1: parse.Int(util.GetJSONValue(appData[dsAppInfo], "1.2.51.1.1.1")),
			2: parse.Int(util.GetJSONValue(appData[dsAppInfo], "1.2.51.1.2.1")),
			3: parse.Int(util.GetJSONValue(appData[dsAppInfo], "1.2.51.1.3.1")),
			4: parse.Int(util.GetJSONValue(appData[dsAppInfo], "1.2.51.1.4.1")),
			5: parse.Int(util.GetJSONValue(appData[dsAppInfo], "1.2.51.1.5.1")),
		}

		screenshots := util.GetJSONArray(appData[dsAppInfo], "1.2.78.0")
		for _, screen := range screenshots {
			app.Screenshots = append(app.Screenshots, util.GetJSONValue(screen.String(), "3.2"))
		}

		for dsAppSimilar := range appData {
			similarURL := util.GetJSONValue(appData[dsAppSimilar], "1.1.1.21.1.2.4.2")
			if similarURL != "" {
				app.SimilarURL, _ = util.AbsoluteURL(playURL, similarURL)
				break
			}
		}

		app.RecentChangesHTML = util.GetJSONValue(appData[dsAppInfo], "1.2.144.1.1", "1.2.145.0.0", "1.2.112.146.0.0")
		app.RecentChanges = util.HTMLToText(app.RecentChangesHTML)
		app.Released = util.GetJSONValue(appData[dsAppInfo], "1.2.10.0")
		app.Score = parse.Float(util.GetJSONValue(appData[dsAppInfo], "1.2.51.0.1"))
		app.ScoreText = util.GetJSONValue(appData[dsAppInfo], "1.2.51.0.0")
		app.Summary = util.GetJSONValue(appData[dsAppInfo], "1.2.73.0.1")
		app.Title = util.GetJSONValue(appData[dsAppInfo], "1.2.0.0")
		app.Updated = time.Unix(parse.Int64(util.GetJSONValue(appData[dsAppInfo], "1.2.145.0.1.0", "1.2.112.146.0.1.0")), 0)
		app.Version = util.GetJSONValue(appData[dsAppInfo], "1.2.140.0.0.0", "1.2.112.141.0.0.0")
		app.Video = util.GetJSONValue(appData[dsAppInfo], "1.2.100.0.0.3.2")
		app.VideoImage = util.GetJSONValue(appData[dsAppInfo], "1.2.100.1.0.3.2")
	}
}

func (app *App) fetchAndExtractData() (map[string]string, error) {
	req, err := http.NewRequest("GET", app.URL, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add(GeographicLocation, app.options.Country)
	q.Add(HostLanguage, app.options.Language)
	req.URL.RawQuery = q.Encode()

	resp, err := app.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close response body")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, ErrNotFound
		}
		
		return nil, fmt.Errorf("request error: %s", resp.Status)
	}

	html, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	appData := util.ExtractInitData(html)
	return appData, nil
}
