package trakt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListItems_Difference(t *testing.T) {
	listA := ListItems{
		{EntityId: 1, Name: "Item1"},
		{EntityId: 2, Name: "Item2"},
		{EntityId: 3, Name: "Item3"},
	}

	tests := []struct {
		name        string
		base        ListItems
		others      []ListItems
		expected    ListItems
		description string
	}{
		{
			name:        "Basic difference",
			base:        listA,
			others:      []ListItems{{{EntityId: 2}, {EntityId: 3}}},
			expected:    ListItems{{EntityId: 1, Name: "Item1"}},
			description: "Removes Item2 and Item3 from base list",
		},
		{
			name:        "No other lists returns clone",
			base:        listA,
			others:      nil,
			expected:    listA,
			description: "No input means return the same items (clone)",
		},
		{
			name:        "Disjoint lists",
			base:        listA,
			others:      []ListItems{{{EntityId: 100, Name: "Different"}}},
			expected:    listA,
			description: "Disjoint items should not be removed",
		},
		{
			name:        "Full overlap",
			base:        listA,
			others:      []ListItems{{{EntityId: 1}, {EntityId: 2}, {EntityId: 3}}},
			expected:    ListItems{},
			description: "All items are removed",
		},
		{
			name:        "Multiple exclusion lists",
			base:        listA,
			others:      []ListItems{{{EntityId: 2}}, {{EntityId: 3}}},
			expected:    ListItems{{EntityId: 1, Name: "Item1"}},
			description: "Removes Item2 and Item3 with separate lists",
		},
		{
			name: "Ignores non-identity fields",
			base: ListItems{
				{EntityId: 1, Name: "Item1", Type: "show"},
				{EntityId: 2, Name: "Item2", Type: "show"},
			},
			others: []ListItems{
				{
					// Same EntityId and Type, but different Name and Id
					{EntityId: 2, Name: "CompletelyDifferentName", Type: "show", Id: 999},
				},
			},
			expected: ListItems{
				{EntityId: 1, Name: "Item1", Type: "show"},
			},
			description: "Item2 should be removed regardless of Name and Id mismatches",
		},
		{
			name: "Does not remove if type mismatches",
			base: ListItems{
				{EntityId: 1, Name: "Item1", Type: "show"},
			},
			others: []ListItems{
				{
					{EntityId: 1, Name: "Item1", Type: "movie"},
				},
			},
			expected: ListItems{
				{EntityId: 1, Name: "Item1", Type: "show"},
			},
			description: "Does not remove if Type does not match (show != movie)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.base.Difference(tt.others...)

			assert.ElementsMatch(t, tt.expected, result, tt.description)

			// Optional: Check for cloning behavior when no others provided
			if tt.others == nil {
				if len(result) > 0 {
					original := tt.base[0].Name
					result[0].Name = "Modified"
					assert.NotEqual(t, original, result[0].Name, "Should be a clone, not the same reference")
				}
			}
		})
	}
}
