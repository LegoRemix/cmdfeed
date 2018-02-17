// Package logic implements the core logic for podcmdr
package logic

import (
	"github.com/LegoRemix/cmdfeed/store"
	"github.com/LegoRemix/cmdfeed/subscription"
	"github.com/pkg/errors"
	"github.com/vmihailenco/msgpack"
)

// podcastNameSpace is where our data about podcasts lives
var podcastNameSpace = []byte("podcast")

type impl struct {
	storage store.Backend
}

// PodcastOptions refers to the options we can have associated with a podcast
type PodcastOptions struct {
	DownloadDirectory *string `msgpack:"download_directory"`
	RecentEntries     *int    `msg:"recent_entries"`
}

// Implementation implements the logic of this
type Implementation interface {
}

// Podcast is the in memory representation of a podcast
type Podcast struct {
	Slug         string             `msgpack:"slug"`
	Subscription subscription.State `msgpack:"sub_state"`
	Downloaded   map[string]string  `msgpack:"downloaded"`
	Options      PodcastOptions     `msgpack:"podcast_options"`
}

// NewImplementation returns a implementation of our logic layer
func NewImplementation() (Implementation, error) {
	backend, err := store.NewLocalBackend()
	if err != nil {
		return nil, errors.Wrap(err, "creating implementation of podcast logic")
	}
	err = backend.CreateNamespace(podcastNameSpace)
	return &impl{storage: backend}, err
}

// NewPodcast adds a new podcast to our set of podcasts
func (backend *impl) NewPodcast(slug string, url string, opts PodcastOptions, subOpts subscription.Options) (Podcast, error) {
	subState, err := subscription.NewState(url, subOpts)
	if err != nil {
		return Podcast{}, errors.Wrap(err, "adding new podcast")
	}

	pod := Podcast{
		Slug:         slug,
		Subscription: subState,
		Downloaded:   make(map[string]string),
		Options:      opts,
	}

	// marshall the object to msgpack
	payload, err := msgpack.Marshal(pod)
	if err != nil {
		return Podcast{}, errors.Wrap(err, "adding new podcast")
	}

	err = backend.storage.Put(podcastNameSpace, []byte(slug), payload)
	if err != nil {
		return Podcast{}, errors.Wrap(err, "adding new podcast")
	}

	return pod, nil

}

// Podcast gets a podcast by slug
func (backend *impl) Podcast(slug string) (Podcast, error) {
	payload, err := backend.storage.Get(podcastNameSpace, []byte(slug))
	if err != nil {
		return Podcast{}, errors.Wrap(err, "getting podcast")
	}

	// unmarshal the podcast
	var pod Podcast
	err = msgpack.Unmarshal(payload, &pod)
	if err != nil {
		return Podcast{}, errors.Wrap(err, "getting podcast")
	}

	return pod, nil
}

// AllPodcasts gets all the podcasts in the store
func (backend *impl) AllPodcasts() ([]Podcast, error) {
	var podcasts []Podcast
	err := backend.storage.ForEach(podcastNameSpace, func(key []byte, value []byte) error {

		var pod Podcast
		err := msgpack.Unmarshal(value, &pod)
		if err != nil {
			return err
		}

		podcasts = append(podcasts, pod)
		return nil
	})

	return podcasts, err
}

// WritePodcast writes an updated podcast to the store
func (backend *impl) WritePodcast(pod Podcast) error {
	payload, err := msgpack.Marshal(pod)
	if err != nil {
		return errors.Wrap(err, "writing podcast - marshall")
	}

	return backend.storage.Put(podcastNameSpace, []byte(pod.Slug), payload)
}
