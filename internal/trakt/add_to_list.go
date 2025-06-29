package trakt

import (
	"fmt"

	mapset "github.com/deckarep/golang-set/v2"
)

type addListItemsBodyItemIds struct {
	Trakt int64 `json:"trakt"`
}

type addListItemsBodyItem struct {
	Ids addListItemsBodyItemIds `json:"ids"`
}

type addListItemsBody struct {
	Movies []addListItemsBodyItem `json:"movies,omitempty"`
	Shows  []addListItemsBodyItem `json:"shows,omitempty"`
}

func (t *Trakt) addToList(listId string, items mapset.Set[ListItem]) error {
	path := getListPath(listId)
	var body addListItemsBody
	for item := range items.Iter() {
		switch item.Type {
		case "movie":
			body.Movies = append(body.Movies,
				addListItemsBodyItem{
					Ids: addListItemsBodyItemIds{
						Trakt: int64(item.Id),
					},
				})
		case "show":
			body.Shows = append(body.Shows,
				addListItemsBodyItem{
					Ids: addListItemsBodyItemIds{
						Trakt: int64(item.Id),
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
