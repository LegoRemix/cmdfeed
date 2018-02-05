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
	ArtworkURL        string   `json:"artworkUrl600"`
	UserRatingCount   int64    `json:"userRatingCount"`
	AverageUserRating float64  `json:"averageUserRating"`
	FeedURL           string   `json:"feedUrl"`
}

// Search returns the list of the results from a query
func Search(param string) (Result, error) {
	//set up the search parameters
	params := make(url.Values)
	params.Set("limit", searchLimit)
	params.Set("media", "podcast")
	params.Set("term", param)

	// make the GET request to the API
	res, err := http.Get(searchURL + params.Encode())
	if err != nil {
		return Result{}, err
	}

	// read the JSON body into our struct
	var result Result
	defer res.Body.Close()
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return Result{}, err
	}

	return result, nil
}
