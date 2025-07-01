package trakt

import "fmt"

func (t *Trakt) getList(listId string) (ListItems, error) {
	path := getListPath(listId)
	var response traktListEntriesResponse
	_, err := t.client.R().
		SetResult(&response).
		Get(path)
	if err != nil {
		return nil, err
	}

	listItems := response.ToListItems()

	return listItems, nil
}

func (t *Trakt) getLists(lists []string) ([]ListItems, error) {
	var allLists []ListItems
	for _, list := range lists {
		listItems, err := t.getList(list)
		if err != nil {
			return nil, fmt.Errorf("failed to get list %s: %w", list, err)
		}
		allLists = append(allLists, listItems)
	}

	return allLists, nil
}

// Deprecated: returns api response instead of ListItems
func (t *Trakt) getListDeprecated(listId string) (traktListEntriesResponse, error) {
	path := getListPath(listId)
	var listEntries traktListEntriesResponse
	_, err := t.client.R().
		SetResult(&listEntries).
		Get(path)
	if err != nil {
		return nil, err
	}

	return listEntries, nil
}

// Deprecated: returns api response instead of ListItems
func (t *Trakt) getListsDeprecated(lists []string) ([]traktListEntriesResponse, error) {
	var result []traktListEntriesResponse
	for _, list := range lists {
		entries, err := t.getListDeprecated(list)
		if err != nil {
			return nil, fmt.Errorf("failed to get list %s: %w", list, err)
		}
		result = append(result, entries)
	}
	return result, nil
}
