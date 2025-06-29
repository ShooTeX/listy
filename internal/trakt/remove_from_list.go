package trakt

import "fmt"

type removeListItemsBodyItemIds struct {
	Trakt int64 `json:"trakt"`
}

type removeListItemsBodyItem struct {
	Ids removeListItemsBodyItemIds `json:"ids"`
}

type removeListItemsBody struct {
	Movies []removeListItemsBodyItem `json:"movies,omitempty"`
	Shows  []removeListItemsBodyItem `json:"shows,omitempty"`
}

func (t *Trakt) removeFromList(list string, items []ListItem) error {
	path := getListPath(list) + "/remove"
	var body removeListItemsBody
	for _, item := range items {
		switch item.Type {
		case "movie":
			body.Movies = append(body.Movies,
				removeListItemsBodyItem{
					Ids: removeListItemsBodyItemIds{
						Trakt: int64(item.EntityId),
					},
				})
		case "show":
			body.Shows = append(body.Shows,
				removeListItemsBodyItem{
					Ids: removeListItemsBodyItemIds{
						Trakt: int64(item.EntityId),
					},
				})
		}
	}
	if _, err := t.client.R().
		SetBody(body).
		Post(path); err != nil {
		return fmt.Errorf("failed to remove items from list %s: %w", list, err)
	}

	return nil
}
