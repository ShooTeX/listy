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

	intersectionDestinationDifference := fullIntersection.Difference(destinationList.ListItemsSet())

	if intersectionDestinationDifference.IsEmpty() {
		return nil
	}

	if err := t.addToList(destination, intersectionDestinationDifference); err != nil {
		return fmt.Errorf("failed to add intersection to list %s: %w", destination, err)
	}

	return nil
}

func (t *Trakt) AddDifferenceToList(lists []string, destination string) error {
	allLists, err := t.getLists(lists)
	if err != nil {
		return fmt.Errorf("failed to get lists: %w", err)
	}

	difference := allLists[0].ListItemsSet().Clone()
	for _, entries := range allLists[1:] {
		difference = difference.Difference(entries.ListItemsSet())
	}

	if difference.IsEmpty() {
		return nil
	}

	if err := t.addToList(destination, difference); err != nil {
		return fmt.Errorf("failed to add intersection to list %s: %w", destination, err)
	}

	return nil
}
