package trakt

import (
	"fmt"
)

func (t *Trakt) addToList(listId string, items []ListItem) error {
	path := getListPath(listId)
	var body addListItemsBody
	for _, item := range items {
		switch item.Type {
		case "movie":
			body.Movies = append(body.Movies,
				addListItemsBodyItem{
					Ids: addListItemsBodyItemIds{
						Trakt: int64(item.EntityId),
					},
				})
		case "show":
			body.Shows = append(body.Shows,
				addListItemsBodyItem{
					Ids: addListItemsBodyItemIds{
						Trakt: int64(item.EntityId),
					},
				})
		}
	}

	_, err := t.client.R().SetBody(body).Post(path)
	if err != nil {
		return fmt.Errorf("failed to add items to list %s: %w", listId, err)
	}

	return nil
}
