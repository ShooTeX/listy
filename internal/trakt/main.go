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

	intersection := allLists[0].ListItemsSet().Clone()
	for _, entries := range allLists[1:] {
		intersection = intersection.Intersect(entries.ListItemsSet())
	}

	if intersection.IsEmpty() {
		return nil
	}

	if err := t.addToList(destination, intersection); err != nil {
		return fmt.Errorf("failed to add intersection to list %s: %w", destination, err)
	}

	return nil
}
