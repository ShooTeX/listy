package trakt

import (
	"fmt"
	"time"
)

type traktHistoryResponse []struct {
	ID        int       `json:"id"`
	WatchedAt time.Time `json:"watched_at"`
	Action    string    `json:"action"`
	Type      string    `json:"type"`
	Movie     struct {
		Title string `json:"title"`
		Year  int    `json:"year"`
		Ids   struct {
			Trakt int    `json:"trakt"`
			Slug  string `json:"slug"`
			Imdb  string `json:"imdb"`
			Tmdb  int    `json:"tmdb"`
		} `json:"ids"`
	} `json:"movie,omitempty"`
	Episode struct {
		Season int    `json:"season"`
		Number int    `json:"number"`
		Title  string `json:"title"`
		Ids    struct {
			Trakt int `json:"trakt"`
			Tvdb  int `json:"tvdb"`
			Imdb  any `json:"imdb"`
			Tmdb  int `json:"tmdb"`
		} `json:"ids"`
	} `json:"episode,omitempty"`
	Show struct {
		Title string `json:"title"`
		Year  int    `json:"year"`
		Ids   struct {
			Trakt int    `json:"trakt"`
			Slug  string `json:"slug"`
			Tvdb  int    `json:"tvdb"`
			Imdb  string `json:"imdb"`
			Tmdb  int    `json:"tmdb"`
		} `json:"ids"`
	} `json:"show,omitempty"`
}

func (t *Trakt) isWatched(itemType string, itemId int) (bool, error) {
	path := "/sync/history/" + itemType + "s/" + fmt.Sprint(itemId)

	var response traktHistoryResponse
	if _, err := t.client.R().
		SetResult(&response).
		Get(path); err != nil {
		return false, fmt.Errorf("failed to check if %s with ID %d is watched: %w", itemType, itemId, err)
	}

	return len(response) > 0, nil
}
