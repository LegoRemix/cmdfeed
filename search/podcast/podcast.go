// Package podcast contains all the clients we use for searching for feeds
package podcast

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// the url for searching itunes
const searchURL = "https://itunes.apple.com/search?"

// limit for the search
const searchLimit = "200"

// Result represent iTunes response outer most structure
type Result struct {
	ResultCount int     `json:"resultCount"`
	Results     []Entry `json:"results"`
}

// Entry holds iTunes track item
type Entry struct {
	Name              string   `json:"collectionName"`
	ArtistName        string   `json:"artistName"`
	Genres            []string `json:"genres"`
	ArtworkURL        string   `json:"artworkUrl512"`
	UserRatingCount   int64    `json:"userRatingCount"`
	AverageUserRating float64  `json:"averageUserRating"`
	Description       string   `json:"description"`
	FeedURL           string   `json:"feedUrl"`
}

// Search returns the list of the results from a query
func Search(param string) (Result, error) {
	//set up the search parameters
	var params url.Values
	params.Set("limit", searchLimit)
	params.Set("media", "podcast")
	params.Set("term", url.QueryEscape(param))

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
