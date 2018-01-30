// Package subscription handles the state of a subscription to a given feed
package subscription

// Options lets one control exactly how
type Options struct {
	// IncludeRemovedEntries controls whether we keep entries not in the current FeedState in the Subscription
	IncludeRemovedEntries bool
	// 
}
