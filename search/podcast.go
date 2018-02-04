// Package search contains all the clients we use for searching for feeds
package search

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

// Result represent iTunes response outer most structure
type Result struct {
	ResultCount int     `json:"resultCount"`
	Results     []Entry `json:"results"`
}

// Entry holds iTunes track item
type Entry struct {
	TrackID                            int64    `json:"trackId"` // Track
	TrackName                          string   `json:"trackName"`
	TrackCensoredName                  string   `json:"trackCensoredName"`
	TrackViewURL                       string   `json:"trackViewUrl"`
	BundleID                           string   `json:"bundleId"` // App bundle
	ArtistID                           int64    `json:"artistId"` // Artist
	ArtistName                         string   `json:"artistName"`
	ArtistViewURL                      string   `json:"artistViewUrl"`
	SellerName                         string   `json:"sellerName"` // Seller
	SellerURL                          string   `json:"sellerUrl"`
	PrimaryGenreID                     int64    `json:"primaryGenreId"` // Genre
	GenreIDs                           []string `json:"genreIds"`
	PrimaryGenreName                   string   `json:"primaryGenreName"`
	Genres                             []string `json:"genres"`
	ArtworkURL60                       string   `json:"artworkUrl60"` // Icon
	ArtworkURL100                      string   `json:"artworkUrl100"`
	ArtworkURL512                      string   `json:"artworkUrl512"`
	Price                              float64  `json:"price"` // Price
	Currency                           string   `json:"currency"`
	FormattedPrice                     string   `json:"formattedPrice"`
	LanguageCodesISO2A                 []string `json:"languageCodesISO2A"` // Platform
	Features                           []string `json:"features"`
	SupportedDevices                   []string `json:"supportedDevices"`
	MinimumOsVersion                   string   `json:"minimumOsVersion"`
	TrackContentRating                 string   `json:"trackContentRating"`
	ContentAdvisoryRating              string   `json:"contentAdvisoryRating"` // Rating
	Advisories                         []string `json:"advisories"`
	UserRatingCount                    int64    `json:"userRatingCount"` // Ranking
	AverageUserRating                  float64  `json:"averageUserRating"`
	UserRatingCountForCurrentVersion   int64    `json:"userRatingCountForCurrentVersion"`
	AverageUserRatingForCurrentVersion float64  `json:"averageUserRatingForCurrentVersion"`
	Kind                               string   `json:"kind"` // Type
	WrapperType                        string   `json:"wrapperType"`
	ScreenshotURLs                     []string `json:"screenshotUrls"` // Screenshots
	IpadScreenshotURLs                 []string `json:"ipadScreenshotUrls"`
	AppletvScreenshotURLs              []string `json:"appletvScreenshotUrls"`
	IsGameCenterEnabled                bool     `json:"isGameCenterEnabled"` // Flags
	IsVppDeviceBasedLicensingEnabled   bool     `json:"isVppDeviceBasedLicensingEnabled"`
	FileSizeBytes                      string   `json:"fileSizeBytes"` // Attribute
	Version                            string   `json:"version"`
	Description                        string   `json:"description"`
	ReleaseNotes                       string   `json:"releaseNotes"`
	ReleaseDate                        string   `json:"releaseDate"`
	CurrentVersionReleaseDate          string   `json:"currentVersionReleaseDate"`
	FeedURL                            string   `json:"feedUrl"`
}

// Params represents the parameters in a
type Params struct {
	url.Values
}

// Finder discovers podcasts based on parameters provided
type Finder interface {
}

type finderImpl struct {
}

const searchURL = "https://itunes.apple.com/search?"

// NewFinder builds a new object for searching for data
func NewFinder() Finder {
	return &finderImpl{}
}

// AddTerm adds a term to our list of search terms
func (params Params) AddTerm(term string) Params {
	params.Values.Add("term", url.QueryEscape(term))
	return params
}

// Country sets the country for this search
func (params Params) Country(country string) Params {
	params.Values.Set("country", country)
	return params
}

// AddEntity adds an entity to the list of entities we are searching
func (params Params) AddEntity(entity string) Params {
	params.Values.Add("entity", entity)
	return params
}

// AddMedia adds the media to the list of entities we are searching
func (params Params) AddMedia(media string) Params {
	params.Values.Add("media", media)
	return params
}

// Limit sets the number of elements to return
func (params Params) Limit(n int) Params {
	if n > 200 {
		n = 200
	}

	if n < 1 {
		n = 1
	}
	params.Values.Set("limit", strconv.Itoa(n))
	return params
}

// Results returns the list of the results from a query
func (s *finderImpl) Search(params Params) (Result, error) {
	res, err := http.Get(searchURL + params.Encode())
	if err != nil {
		return Result{}, err
	}

	var result Result
	defer res.Body.Close()
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return Result{}, err
	}

	return result, nil
}
