// Package rss implements a wrapper around an rss feed that exposes
// useful properties about the feed including elements, hash, etc.
package rss

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
)

// State is an abstract representation of an RSS/Atom Feed state with additional properties and methods
type State interface {
	// FetchTime is when this instance of the feed was fetched
	FetchTime() time.Time
	// Hash returns an MD5 hash of the feed from when it was fetched
	Hash() string
	// UpdatedState goes over the network and fetches the current state of the RSS feed
	UpdatedState() (State, error)
	// Feed returns the underlying Feed object
	Feed() Feed
}

// impl is the actual implementation of a Feed
type impl struct {
	URL         string    `json:"url,omitempty"`
	HashString  string    `json:"hash,omitempty"`
	FetchedTime time.Time `json:"fetchTime,omitempty"`
	FeedInfo    Feed      `json:"feed,omitempty"`
}

// Image is an image that is the artwork for a given
type Image struct {
	URL   string `json:"url,omitempty"`
	Title string `json:"title,omitempty"`
}

// Feed is an RSS Feed
type Feed struct {
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	Link        string     `json:"link,omitempty"`
	FeedLink    string     `json:"feedLink,omitempty"`
	Updated     *time.Time `json:"updated,omitempty"`
	Published   *time.Time `json:"published,omitempty"`
	Author      *Person    `json:"author,omitempty"`
	Language    string     `json:"language,omitempty"`
	Image       *Image     `json:"image,omitempty"`
	Copyright   string     `json:"copyright,omitempty"`
	Categories  []string   `json:"categories,omitempty"`
	Items       []*Item    `json:"items"`
}

// Item is the universal Item type that atom.Entry
// and rss.Item gets translated to.  It represents
// a single entry in a given feed.
type Item struct {
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	Content     string     `json:"content,omitempty"`
	Link        string     `json:"link,omitempty"`
	Updated     *time.Time `json:"updated,omitempty"`
	Published   *time.Time `json:"published,omitempty"`
	Author      *Person    `json:"author,omitempty"`
	GUID        string     `json:"guid,omitempty"`
	Image       *Image     `json:"image,omitempty"`
	Categories  []string   `json:"categories,omitempty"`
}

// Person is an individual specified in a feed
type Person struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

// FetchTime returns the time at which this feed was fetched
func (feed *impl) FetchTime() time.Time {
	return feed.FetchedTime
}

// Hash returns the hash value of the feed state
func (feed *impl) Hash() string {
	return feed.HashString
}

// Feed returns the Feed object inside this state
func (feed *impl) Feed() Feed {
	return feed.FeedInfo
}

// UpdatedState returns an updated version of the rss feed
func (feed *impl) UpdatedState() (State, error) {
	return NewState(feed.URL)
}

// feedStateFromJSON constructs a State from JSON
func feedStateFromJSON(data []byte) (*impl, error) {
	state := new(impl)
	err := json.Unmarshal(data, state)
	return state, err
}

// NewState returns a new instance of the RichFeed type
func NewState(url string) (State, error) {
	// attempt to fetch the feed from a given url
	response, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch rss feed")
	}

	//read the body of the response into a byte slice
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read rss body")
	}

	// hash the contents
	hasher := md5.New()
	hasher.Write(contents)
	hash := hex.EncodeToString(hasher.Sum(nil))

	// parse the rss feed into intelligible data
	parser := gofeed.NewParser()
	parsedFeed, err := parser.Parse(bytes.NewReader(contents))
	if err != nil {
		return nil, err
	}
	fetchTime := time.Now().UTC()

	// Convert all of our items into the proper format
	items := make([]*Item, 0, len(parsedFeed.Items))
	for _, item := range parsedFeed.Items {
		if item == nil {
			continue
		}

		// Extract the author if it exists
		var author *Person
		if item.Author != nil {
			author = new(Person)
			*author = Person(*item.Author)
		}

		// Extract the Item Image, if it exists
		var image *Image
		if item.Image != nil {
			image = new(Image)
			*image = Image(*item.Image)
		}

		ourItem := &Item{
			Updated:     item.UpdatedParsed,
			Published:   item.PublishedParsed,
			Title:       item.Title,
			Description: item.Description,
			Content:     item.Content,
			Author:      author,
			Image:       image,
			Categories:  item.Categories,
			Link:        item.Link,
			GUID:        item.GUID,
		}

		items = append(items, ourItem)
	}

	// Extract the feed Image, if it exists
	var image *Image
	if parsedFeed.Image != nil {
		image = new(Image)
		*image = Image(*parsedFeed.Image)
	}

	// Extract the author if it exists
	var author *Person
	if parsedFeed.Author != nil {
		author = new(Person)
		*author = Person(*parsedFeed.Author)
	}

	feed := Feed{
		Author:      author,
		Image:       image,
		Items:       items,
		Title:       parsedFeed.Title,
		Description: parsedFeed.Description,
		Link:        parsedFeed.Link,
		FeedLink:    parsedFeed.Link,
		Updated:     parsedFeed.UpdatedParsed,
		Published:   parsedFeed.PublishedParsed,
		Language:    parsedFeed.Language,
		Categories:  parsedFeed.Categories,
		Copyright:   parsedFeed.Copyright,
	}

	return &impl{
		FeedInfo:    feed,
		URL:         url,
		HashString:  hash,
		FetchedTime: fetchTime,
	}, nil
}
