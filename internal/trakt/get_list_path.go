package trakt

func getListPath(list string) string {
	if list == "watchlist" {
		return "/users/me/watchlist"
	}

	return "/users/me/lists/" + list + "/items"
}
