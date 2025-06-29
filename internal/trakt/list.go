package trakt

import (
	"time"

	mapset "github.com/deckarep/golang-set/v2"
)

type traktListEntries []struct {
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

type ListItem struct {
	Name string
	Type string
	Id   int
}

func (e *traktListEntries) ListItemsSlice() []ListItem {
	items := make([]ListItem, 0, len(*e))
	for _, entry := range *e {
		switch entry.Type {
		case "movie":
			item := ListItem{
				Name: entry.Movie.Title,
				Type: entry.Type,
				Id:   entry.Movie.Ids.Trakt,
			}
			items = append(items, item)
		case "show":
			item := ListItem{
				Name: entry.Show.Title,
				Type: entry.Type,
				Id:   entry.Show.Ids.Trakt,
			}
			items = append(items, item)
		}
	}

	return items
}

func (e *traktListEntries) ListItemsSet() mapset.Set[ListItem] {
	return mapset.NewSet(e.ListItemsSlice()...)
}
