package trakt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetListPath(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		description string
	}{
		{
			name:        "Current user watchlist",
			input:       "watchlist",
			expected:    "/users/me/watchlist",
			description: "Default watchlist should use 'me' as username",
		},
		{
			name:        "Current user custom list",
			input:       "mylist",
			expected:    "/users/me/lists/mylist/items",
			description: "Custom list without username should use 'me' as username",
		},
		{
			name:        "Specific user watchlist",
			input:       "shootex/watchlist",
			expected:    "/users/shootex/watchlist",
			description: "Watchlist for specific user should use provided username",
		},
		{
			name:        "Specific user custom list",
			input:       "shootex/favorites",
			expected:    "/users/shootex/lists/favorites/items",
			description: "Custom list for specific user should use provided username",
		},
		{
			name:        "Empty list name",
			input:       "",
			expected:    "/users/me/lists//items",
			description: "Empty string should result in empty list name",
		},
		{
			name:        "List name with special characters",
			input:       "my-special_list",
			expected:    "/users/me/lists/my-special_list/items",
			description: "List names with special characters should be preserved",
		},
		{
			name:        "Username with special characters",
			input:       "user-name_123/my-list",
			expected:    "/users/user-name_123/lists/my-list/items",
			description: "Usernames and list names with special characters should be preserved",
		},
		{
			name:        "Multiple slashes in input",
			input:       "user/list/extra",
			expected:    "/users/user/lists/list/extra/items",
			description: "Multiple slashes split on first slash only, rest becomes list name",
		},
		{
			name:        "Slash at beginning",
			input:       "/watchlist",
			expected:    "/users//watchlist",
			description: "Leading slash results in empty username, watchlist recognized",
		},
		{
			name:        "Slash at end",
			input:       "user/",
			expected:    "/users/user/lists//items",
			description: "Trailing slash results in empty list name",
		},
		{
			name:        "Case sensitive watchlist",
			input:       "Watchlist",
			expected:    "/users/me/lists/Watchlist/items",
			description: "Watchlist check is case sensitive - should treat as regular list",
		},
		{
			name:        "Case sensitive user watchlist",
			input:       "user/Watchlist",
			expected:    "/users/user/lists/Watchlist/items",
			description: "User watchlist check is case sensitive - should treat as regular list",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getListPath(tt.input)
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}