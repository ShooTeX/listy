package trakt

import (
	"context"
	"fmt"
	"slices"

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

func (t *Trakt) AddIntersectToList(lists []string, destination string) error {
	allLists, err := t.getLists(lists)
	if err != nil {
		return fmt.Errorf("failed to get lists: %w", err)
	}
	destinationList, err := t.getList(destination)
	if err != nil {
		return fmt.Errorf("failed to get destination list %s: %w", destination, err)
	}

	fullIntersection := allLists[0].ListItemsSet().Clone()
	for _, entries := range allLists[1:] {
		fullIntersection = fullIntersection.Intersect(entries.ListItemsSet())
	}

	intersectionWithoutDestinationItems := fullIntersection.Difference(destinationList.ListItemsSet())

	if intersectionWithoutDestinationItems.IsEmpty() {
		return nil
	}

	var orderedIntersection []ListItem
	for _, item := range allLists[0].ListItemsSlice() {
		if intersectionWithoutDestinationItems.Contains(item) {
			orderedIntersection = append(orderedIntersection, item)
		}
	}

	if err := t.addToList(destination, orderedIntersection); err != nil {
		return fmt.Errorf("failed to add intersection to list %s: %w", destination, err)
	}

	return nil
}

func (t *Trakt) AddDifferenceToList(lists []string, destination string) error {
	allLists, err := t.getLists(lists)
	if err != nil {
		return fmt.Errorf("failed to get lists: %w", err)
	}
	destinationList, err := t.getList(destination)
	if err != nil {
		return fmt.Errorf("failed to get destination list %s: %w", destination, err)
	}

	fullDifference := allLists[0].ListItemsSet().Clone()
	for _, entries := range allLists[1:] {
		fullDifference = fullDifference.Difference(entries.ListItemsSet())
	}

	differenceWithoutDestinationItems := fullDifference.Difference(destinationList.ListItemsSet())

	if differenceWithoutDestinationItems.IsEmpty() {
		return nil
	}

	var orderedDifference []ListItem
	for _, item := range allLists[0].ListItemsSlice() {
		if differenceWithoutDestinationItems.Contains(item) {
			orderedDifference = append(orderedDifference, item)
		}
	}

	if err := t.addToList(destination, orderedDifference); err != nil {
		return fmt.Errorf("failed to add intersection to list %s: %w", destination, err)
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
		for _, item := range listItems {
			isItemWatched, err := t.isWatched(item.Type, item.EntityId)
			if err != nil {
				return fmt.Errorf("failed to check if item %s with ID %d is watched: %w", item.Type, item.EntityId, err)
			}

			if isItemWatched {
				itemsToRemove = append(itemsToRemove, item)
			}
		}
	}

	if err := t.removeFromList(list, itemsToRemove); err != nil {
		return fmt.Errorf("failed to remove watched items from list %s: %w", list, err)
	}

	return nil
}
