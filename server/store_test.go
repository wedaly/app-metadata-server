package server

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStoreInsertAndSearch(t *testing.T) {
	testCases := []struct {
		name     string
		apps     []App
		matcher  Matcher
		expected []App
	}{
		{
			name:     "no records, no results",
			apps:     nil,
			matcher:  MatchAny,
			expected: nil,
		},
		{
			name: "single record, matches",
			apps: []App{
				{Title: "test"},
			},
			matcher: MatchExactTitle("test"),
			expected: []App{
				{Title: "test"},
			},
		},
		{
			name: "single record, no match",
			apps: []App{
				{Title: "test"},
			},
			matcher:  MatchExactTitle("other"),
			expected: nil,
		},
		{
			name: "many records, some match",
			apps: []App{
				{Title: "foo", Description: "a"},
				{Title: "bar", Description: "b"},
				{Title: "baz", Description: "c"},
				{Title: "bat", Description: "ab"},
			},
			matcher: MatchDescriptionContains("b"),
			expected: []App{
				{Title: "bar", Description: "b"},
				{Title: "bat", Description: "ab"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var store Store
			for _, app := range tc.apps {
				store.Insert(app)
			}
			var results []App
			store.Search(tc.matcher, func(app App) {
				results = append(results, app)
			})
			assert.Equal(t, tc.expected, results)
		})
	}
}

func BenchmarkSearch(b *testing.B) {
	benchmarks := []struct {
		numApps     int
		numSelected int
	}{
		{
			numApps:     10,
			numSelected: 0,
		},
		{
			numApps:     10,
			numSelected: 5,
		},
		{
			numApps:     10,
			numSelected: 10,
		},
		{
			numApps:     1000,
			numSelected: 0,
		},
		{
			numApps:     1000,
			numSelected: 500,
		},
		{
			numApps:     1000,
			numSelected: 1000,
		},
		{
			numApps:     10000,
			numSelected: 0,
		},
		{
			numApps:     10000,
			numSelected: 5000,
		},
		{
			numApps:     10000,
			numSelected: 10000,
		},
	}

	for _, bm := range benchmarks {
		name := fmt.Sprintf("%d apps, %d selected", bm.numApps, bm.numSelected)
		b.Run(name, func(b *testing.B) {
			var store Store

			// Setup: insert apps into the store
			for i := 0; i < bm.numApps; i++ {
				var title string
				if i < bm.numSelected {
					title = "selected"
				} else {
					title = "notselected"
				}

				store.Insert(App{
					Title:       title,
					Version:     "0.0.0",
					Description: "description",
					Maintainers: []Maintainer{
						{Name: "name", Email: "name@example.com"},
					},
					Company: "company",
					Website: "http://example.com",
					Source:  "https://git.example.com/repo",
					License: "license",
				})
			}

			// Benchmark: time search
			m := MatchExactTitle("selected")
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				store.Search(m, func(App) {})
			}
		})
	}
}
