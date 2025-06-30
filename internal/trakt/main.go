package trakt

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/shootex/listy/internal/auth"
	"resty.dev/v3"
)

type Trakt struct {
	client *resty.Client
	ctx    context.Context
}

func New(ctx context.Context) (*Trakt, error) {
	client, err := auth.NewClient(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &Trakt{
		client: client,
		ctx:    ctx,
	}, nil
}

func (t *Trakt) AddIntersectToList(lists []string, destination string, clean bool) error {
	allLists, err := t.getLists(lists)
	if err != nil {
		return fmt.Errorf("failed to get lists: %w", err)
	}
	destinationList, err := t.getList(destination)
	if err != nil {
		return fmt.Errorf("failed to get destination list %s: %w", destination, err)
	}

	intersection := allLists[0].ListItemsSet().Clone()
	for _, entries := range allLists[1:] {
		intersection = intersection.Intersect(entries.ListItemsSet())
	}

	intersection = intersection.Difference(destinationList.ListItemsSet())

	if intersection.IsEmpty() {
		return nil
	}

	var orderedIntersection []ListItem
	for _, item := range allLists[0].ListItemsSlice() {
		if intersection.Contains(item) {
			orderedIntersection = append(orderedIntersection, item)
		}
	}

	if clean {
		if err := t.removeFromList(destination, destinationList.ListItemsSlice()); err != nil {
			return fmt.Errorf("failed to remove unknown items from list %s: %w", destination, err)
		}
	}

	if err := t.addToList(destination, orderedIntersection); err != nil {
		return fmt.Errorf("failed to add intersection to list %s: %w", destination, err)
	}

	return nil
}

func (t *Trakt) AddDifferenceToList(lists []string, destination string, clean bool) error {
	allLists, err := t.getLists(lists)
	if err != nil {
		return fmt.Errorf("failed to get lists: %w", err)
	}
	destinationList, err := t.getList(destination)
	if err != nil {
		return fmt.Errorf("failed to get destination list %s: %w", destination, err)
	}

	difference := allLists[0].ListItemsSet().Clone()
	for _, entries := range allLists[1:] {
		difference = difference.Difference(entries.ListItemsSet())
	}

	difference = difference.Difference(destinationList.ListItemsSet())

	if difference.IsEmpty() {
		return nil
	}

	var orderedDifference []ListItem
	for _, item := range allLists[0].ListItemsSlice() {
		if difference.Contains(item) {
			orderedDifference = append(orderedDifference, item)
		}
	}

	if clean {
		if err := t.removeFromList(destination, destinationList.ListItemsSlice()); err != nil {
			return fmt.Errorf("failed to remove unknown items from list %s: %w", destination, err)
		}
	}

	if err := t.addToList(destination, orderedDifference); err != nil {
		return fmt.Errorf("failed to add difference to list %s: %w", destination, err)
	}

	return nil
}

func (t *Trakt) CopyListOrder(list, destination string) error {
	fromList, err := t.getList(list)
	if err != nil {
		return fmt.Errorf("failed to get list %s: %w", list, err)
	}

	destinationList, err := t.getList(destination)
	if err != nil {
		return fmt.Errorf("failed to get destination list %s: %w", destination, err)
	}

	orderMap := make(map[int]int)
	for i, item := range fromList.ListItemsSlice() {
		orderMap[item.EntityId] = i
	}

	orderedItems := destinationList.ListItemsSlice()

	slices.SortFunc(orderedItems, func(a, b ListItem) int {
		aIdx, aOk := orderMap[a.EntityId]
		bIdx, bOk := orderMap[b.EntityId]

		switch {
		case aOk && bOk:
			return aIdx - bIdx
		case aOk && !bOk:
			return -1
		case !aOk && bOk:
			return 1
		default:
			return 0
		}
	})

	if err := t.updateListOrder(destination, orderedItems); err != nil {
		return fmt.Errorf("failed to update order of list %s: %w", destination, err)
	}

	return nil
}

type CleanOptions struct {
	Watched bool
}

func (t *Trakt) Clean(list string, options *CleanOptions) error {
	listResponse, err := t.getList(list)
	if err != nil {
		return fmt.Errorf("failed to get list %s: %w", list, err)
	}

	listItems := listResponse.ListItemsSlice()

	var itemsToRemove []ListItem
	if options.Watched {
		sem := make(chan struct{}, 100)
		var wg sync.WaitGroup
		var mu sync.Mutex

		for _, item := range listItems {
			sem <- struct{}{}
			wg.Add(1)

			go func(item ListItem) {
				defer wg.Done()
				defer func() { <-sem }()

				isItemWatched, err := t.isWatched(item.Type, item.EntityId)
				if err != nil {
					// NOTE: silent skip, implement error handling maybe
					return
				}

				if isItemWatched {
					mu.Lock()
					itemsToRemove = append(itemsToRemove, item)
					mu.Unlock()
				}
			}(item)
		}

		wg.Wait()
	}

	if err := t.removeFromList(list, itemsToRemove); err != nil {
		return fmt.Errorf("failed to remove watched items from list %s: %w", list, err)
	}

	return nil
}
