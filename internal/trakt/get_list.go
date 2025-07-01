package trakt

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

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

func (t *Trakt) getLists(ctx context.Context, lists []string) ([]ListItems, error) {
	g, _ := errgroup.WithContext(ctx)

	g.SetLimit(100)

	var mu sync.Mutex

	allLists := make([]ListItems, len(lists))
	for i, list := range lists {
		g.Go(func() error {
			listItems, err := t.getList(list)
			if err != nil {
				return fmt.Errorf("failed to get list %s: %w", list, err)
			}
			mu.Lock()
			allLists[i] = listItems
			mu.Unlock()
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("failed to get lists: %w", err)
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
