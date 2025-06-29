package trakt

import (
	"context"
	"fmt"

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
