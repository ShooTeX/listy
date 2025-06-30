package trakt

import "fmt"

func (t *Trakt) updateListOrder(listId string, items []ListItem) error {
	path := getListPath(listId) + "/reorder"

	var body updateListOrderBody
	for _, item := range items {
		body.Rank = append(body.Rank, item.Id)
	}

	_, err := t.client.R().
		SetBody(body).
		Post(path)
	if err != nil {
		return fmt.Errorf("failed to update order of list %s: %w", listId, err)
	}

	return nil
}
