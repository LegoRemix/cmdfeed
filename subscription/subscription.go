// Package subscription handles the state of a subscription to a given feed
package subscription

import (
	"github.com/pkg/errors"
	"encoding/hex"
	"encoding/json"
	"sort"
	"time"

	"github.com/LegoRemix/rss-cli/rss"
)

// State represents the current state of a subscription to a feed
type State interface {
	EntryList() []Entry
}

type impl struct {
	Entries []Entry   `json:"entries"`
	Feed    rss.State `json:"feed"`
	Opts    Options   `json:"options"`
}

// Options lets one control exactly how a subscription state is managed
type Options struct {
	// IncludeRemovedEntries controls whether we keep entries not in the current FeedState in the Subscription
	IncludeRemovedEntries bool `json:"include_removed_entries,omitempty"`
}

// Entry represents a single entry in a subscription feed, it has slightly different semantics from rss.Item
type Entry struct {
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Content     string    `json:"content,omitempty"`
	Link        string    `json:"link,omitempty"`
	Updated     time.Time `json:"updated,omitempty"`
	Published   time.Time `json:"published,omitempty"`
	GUID        string    `json:"guid,omitempty"`
	Categories  []string  `json:"categories,omitempty"`
	ImageTitle  string    `json:"image_title,omitempty"`
	ImageURL    string    `json:"image_url, omitempty"`
	AuthorName  string    `json:"author_name,omitempty"`
	AuthorEmail string    `json:"author_email,omitempty"`
}

// EntryList returns the list of entries in this sub State
func (s *impl) EntryList() []Entry {
	return s.Entries
}

// ID creates a unique ID for the entry
func (e Entry) ID() (string, error) {
	if e.GUID != "" {
		return e.GUID, nil
	}

	hashed, err := json.Marshal(e)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hashed), nil
}

// feedStateToEntries extracts all of the Entrys from a FeedState
func feedStateToEntries(rssState rss.State) []Entry {
	items := rssState.Feed().Items
	result := make([]Entry, 0, len(items))
	for _, item := range items {

		// grab a valid update time for this sub entry
		updated := rssState.FetchTime()
		if item.Updated != nil {
			// if we have a valid update time, we hash that this so we can easily see updates
			updated = item.Updated.UTC()
		}
		// grab a valid publish time for this sub entry
		published := rssState.FetchTime()
		if item.Published != nil {
			published = *item.Published
		}

		entry := Entry{
			Title:       item.Title,
			Description: item.Description,
			Content:     item.Content,
			Link:        item.Link,
			Updated:     updated,
			Published:   published,
			Categories:  item.Categories,
			GUID:        item.GUID,
		}

		if item.Author != nil {
			entry.AuthorName = item.Author.Name
			entry.AuthorEmail = item.Author.Email
		}

		if item.Image != nil {
			entry.ImageURL = item.Image.URL
			entry.ImageTitle = item.Image.Title
		}

		result = append(result, entry)
	}

	return result
}

// mergeEntries merges two lists of entries and then sorts by updateTime
func mergeEntries(left, right []Entry) ([]Entry, error) {
	// deduplicate the two entry lists
	entrySet := make(map[string]Entry)
	for _, lst := range [][]Entry{left, right} {
		for _, entry := range lst {
			id, err := entry.ID()
			if err != nil {
				return nil, err
			}
			entrySet[id] = entry
		}
	}
	// unpack the set back into an array
	result := make([]Entry, 0, len(entrySet))
	for _, v := range entrySet {
		result = append(result, v)
	}

	// sort the entries by updated time
	sort.Slice(result, func(i, j int) bool { return result[i].Updated.Before(result[j].Updated) })

	return result, nil

}

// NewState creates a new subscription.State instance against a given URL
func NewState(url string, opt Options) (State, error) {
	feedState, err := rss.NewState(url)
	if err != nil {
		return nil, err
	}

	entries := feedStateToEntries(feedState)

	return &impl{
		Feed:    feedState,
		Entries: entries,
		Opts:    opt,
	}, nil
}

// Update creates a newly State of the subscription with updated entries
func (s *impl) Update() (State, error) {
	updated, err := s.Feed.UpdatedState()
	if err != nil {
		return nil, err
	}

	//create a newly updated list of entries
	newEntries := feedStateToEntries(updated)
	if s.Opts.IncludeRemovedEntries {
		newEntries, err = mergeEntries(newEntries, s.Entries)
		if err != nil {
			return nil, errors.Wrap(err, "failed to merge subscription entries")
		}
	}

	return &impl{
		Feed:    updated,
		Entries: newEntries,
		Opts:    s.Opts,
	}, nil
}
