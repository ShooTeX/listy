package trakt

import "strings"

func getListPath(list string) string {
	username := "me"
	listName := list

	// Split into at most 2 parts to handle user/list format
	parts := strings.SplitN(list, "/", 2)
	if len(parts) == 2 {
		username = parts[0]
		listName = parts[1]
	}

	if listName == "watchlist" {
		return "/users/" + username + "/watchlist"
	}

	return "/users/" + username + "/lists/" + listName + "/items"
}
