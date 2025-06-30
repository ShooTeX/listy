package trakt

import "time"

type traktListEntriesResponse []struct {
	Rank     int       `json:"rank"`
	Id       int       `json:"id"`
	ListedAt time.Time `json:"listed_at"`
	Notes    any       `json:"notes"`
	Type     string    `json:"type"`
	Movie    struct {
		Title string `json:"title"`
		Year  int    `json:"year"`
		Ids   struct {
			Trakt int    `json:"trakt"`
			Slug  string `json:"slug"`
			Imdb  string `json:"imdb"`
			Tmdb  int    `json:"tmdb"`
		} `json:"ids"`
	} `json:"movie"`
	Show struct {
		Title string `json:"title"`
		Year  int    `json:"year"`
		Ids   struct {
			Trakt  int    `json:"trakt"`
			Slug   string `json:"slug"`
			Tvdb   int    `json:"tvdb"`
			Imdb   string `json:"imdb"`
			Tmdb   int    `json:"tmdb"`
			Tvrage any    `json:"tvrage"`
		} `json:"ids"`
		AiredEpisodes int `json:"aired_episodes"`
	} `json:"show"`
}

type addListItemsBodyItemIds struct {
	Trakt int64 `json:"trakt"`
}

type addListItemsBodyItem struct {
	Ids addListItemsBodyItemIds `json:"ids"`
}

type addListItemsBody struct {
	Movies []addListItemsBodyItem `json:"movies,omitempty"`
	Shows  []addListItemsBodyItem `json:"shows,omitempty"`
}

type updateListOrderBody struct {
	Rank []int `json:"rank"`
}
