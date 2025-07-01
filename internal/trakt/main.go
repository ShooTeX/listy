package trakt

import (
	"context"
	"fmt"
	"slices"
	"sync"

	mapset "github.com/deckarep/golang-set/v2"
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

func (t *Trakt) AddIntersectToList(ctx context.Context, lists []string, destination string, clean bool) error {
	allLists, err := t.getLists(ctx, lists)
	if err != nil {
		return fmt.Errorf("failed to get lists: %w", err)
	}
	destinationList, err := t.getList(destination)
	if err != nil {
		return fmt.Errorf("failed to get destination list %s: %w", destination, err)
	}

	unorderedIntersection := allLists[0].Intersection(allLists[1:]...)

	var intersection []ListItem
	for _, item := range allLists[0] {
		set := mapset.NewSet(unorderedIntersection...)
		if set.Contains(item) {
			intersection = append(intersection, item)
		}
	}

	if clean && len(destinationList) > 0 {
		if err := t.removeFromList(destination, destinationList); err != nil {
			return fmt.Errorf("failed to remove unknown items from list %s: %w", destination, err)
		}
	}

	if err := t.addToList(destination, intersection); err != nil {
		return fmt.Errorf("failed to add intersection to list %s: %w", destination, err)
	}

	return nil
}

func (t *Trakt) AddDifferenceToList(ctx context.Context, lists []string, destination string, clean bool) error {
	allLists, err := t.getLists(ctx, lists)
	if err != nil {
		return fmt.Errorf("failed to get lists: %w", err)
	}
	destinationList, err := t.getList(destination)
	if err != nil {
		return fmt.Errorf("failed to get destination list %s: %w", destination, err)
	}

	unorderedDifference := allLists[0].Difference(allLists[1:]...)

	var difference ListItems
	for _, item := range allLists[0] {
		set := mapset.NewSet(unorderedDifference...)
		if set.Contains(item) {
			difference = append(difference, item)
		}
	}

	if clean && len(destinationList) > 0 {
		if err := t.removeFromList(destination, destinationList); err != nil {
			return fmt.Errorf("failed to remove unknown items from list %s: %w", destination, err)
		}
	}

	if err := t.addToList(destination, difference); err != nil {
		return fmt.Errorf("failed to add difference to list %s: %w", destination, err)
	}

	return nil
}

func (t *Trakt) CopyListOrder(list, destination string) error {
	fromList, err := t.getListDeprecated(list)
	if err != nil {
		return fmt.Errorf("failed to get list %s: %w", list, err)
	}

	destinationList, err := t.getListDeprecated(destination)
	if err != nil {
		return fmt.Errorf("failed to get destination list %s: %w", destination, err)
	}

	orderMap := make(map[int]int)
	for i, item := range fromList.ToListItems() {
		orderMap[item.EntityId] = i
	}

	orderedItems := destinationList.ToListItems()

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
	listResponse, err := t.getListDeprecated(list)
	if err != nil {
		return fmt.Errorf("failed to get list %s: %w", list, err)
	}

	listItems := listResponse.ToListItems()

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
