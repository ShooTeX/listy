package trakt

import mapset "github.com/deckarep/golang-set/v2"

type ListItem struct {
	Id       int
	Name     string
	Type     string
	EntityId int
}

type ListItems []ListItem

type comparableListItem struct {
	Type     string
	EntityId int
}

type compareableListItems []comparableListItem

func (li ListItem) Identity() comparableListItem {
	return comparableListItem{
		Type:     li.Type,
		EntityId: li.EntityId,
	}
}

func (l ListItems) Identities() []comparableListItem {
	ids := make([]comparableListItem, 0, len(l))
	for _, item := range l {
		ids = append(ids, item.Identity())
	}
	return ids
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

	exclude := make(map[comparableListItem]struct{})
	for _, other := range otherLists {
		for _, item := range other {
			exclude[item.Identity()] = struct{}{}
		}
	}

	result := make(ListItems, 0, len(l))
	for _, item := range l {
		if _, found := exclude[item.Identity()]; !found {
			result = append(result, item)
		}
	}

	return result
}

func (l ListItems) Intersection(otherLists ...ListItems) ListItems {
	if len(otherLists) == 0 {
		return l.Clone()
	}

	otherMaps := make([]map[comparableListItem]struct{}, len(otherLists))
	for i, other := range otherLists {
		m := make(map[comparableListItem]struct{}, len(other))
		for _, item := range other {
			m[item.Identity()] = struct{}{}
		}
		otherMaps[i] = m
	}

	result := make(ListItems, 0, len(l))
	for _, item := range l {
		id := item.Identity()
		inAll := true
		for _, m := range otherMaps {
			if _, found := m[id]; !found {
				inAll = false
				break
			}
		}
		if inAll {
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
