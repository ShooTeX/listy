package trakt

import "fmt"

func (t *Trakt) getList(listId string) (traktListEntriesResponse, error) {
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

func (t *Trakt) getLists(lists []string) ([]traktListEntriesResponse, error) {
	var result []traktListEntriesResponse
	for _, list := range lists {
		entries, err := t.getList(list)
		if err != nil {
			return nil, fmt.Errorf("failed to get list %s: %w", list, err)
		}
		result = append(result, entries)
	}
	return result, nil
}
