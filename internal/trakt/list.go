package trakt

import mapset "github.com/deckarep/golang-set/v2"

type ListItem struct {
	Id       int
	Name     string
	Type     string
	EntityId int
}

type ListItems []ListItem

type compareableListItem struct {
	Type     string
	EntityId int
}

func (li ListItem) Identity() compareableListItem {
	return compareableListItem{
		Type:     li.Type,
		EntityId: li.EntityId,
	}
}

func (l ListItems) Clone() ListItems {
	cloned := make(ListItems, len(l))
	copy(cloned, l)
	return cloned
}

func (l ListItems) Difference(otherLists ...ListItems) ListItems {
	if len(otherLists) == 0 {
		return l.Clone()
	}

	excludeSet := mapset.NewSet[compareableListItem]()
	for _, other := range otherLists {
		for _, item := range other {
			excludeSet.Add(item.Identity())
		}
	}

	result := make(ListItems, 0)
	for _, item := range l {
		if !excludeSet.Contains(item.Identity()) {
			result = append(result, item)
		}
	}

	return result
}

func (e *traktListEntriesResponse) ToListItems() ListItems {
	items := make(ListItems, 0, len(*e))
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

func (e *traktListEntriesResponses) ToListItems() []ListItems {
	if e == nil || len(*e) == 0 {
		return nil
	}

	items := make([]ListItems, 0, len(*e))
	for _, entry := range *e {
		items = append(items, entry.ToListItems())
	}

	return items
}

func (e *traktListEntriesResponse) ListItemsSet() mapset.Set[ListItem] {
	return mapset.NewSet(e.ToListItems()...)
}
