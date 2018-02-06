// Package feed contains functions for searching for general rss feeds
package feed

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

// searchURL is the host from which we'll be getting our results
const searchURL = "https://cloud.feedly.com/v3/search/feeds?"

// feedIDPrefix in front of the rss feed
const feedIDPrefix = "feed/"

// Result contains our search result for search
type Result struct {
	Hint    string   `json:"hint"`
	Results []Entry  `json:"results"`
	Related []string `json:"related"`
}

// Entry contains all the information for a search result
type Entry struct {
	Title   string `json:"title"`
	Website string `json:"website"`
	FeedID  string `json:"feedId"`
}

// FeedURL gets the url of the feed
func (e Entry) FeedURL() string {
	return strings.Replace(e.FeedID, feedIDPrefix, "", 1)
}

// Search fetches a list of feeds based on entries
func Search(query string) (Result, error) {

	//create the params
	params := make(url.Values)
	params.Set("query", query)

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
