package trakt

import mapset "github.com/deckarep/golang-set/v2"

type ListItem struct {
	Id       int
	Name     string
	Type     string
	EntityId int
}

// func (l *[]ListItem) computeDifference(*ListItem...)

func (e *traktListEntriesResponse) toListItems() []ListItem {
	items := make([]ListItem, 0, len(*e))
	for _, entry := range *e {
		switch entry.Type {
		case "movie":
			item := ListItem{
				Id:       entry.Id,
				Name:     entry.Movie.Title,
				Type:     entry.Type,
				EntityId: entry.Movie.Ids.Trakt,
			}
			items = append(items, item)
		case "show":
			item := ListItem{
				Id:       entry.Id,
				Name:     entry.Show.Title,
				Type:     entry.Type,
				EntityId: entry.Show.Ids.Trakt,
			}
			items = append(items, item)
		}
	}

	return items
}

func (e *traktListEntriesResponse) ListItemsSet() mapset.Set[ListItem] {
	return mapset.NewSet(e.toListItems()...)
}
