package server

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateApp(t *testing.T) {
	testCases := []struct {
		name     string
		app      App
		expected ValidationErrors
	}{
		{
			name: "valid app",
			app: App{
				Title:   "title",
				Version: "version",
				Maintainers: []Maintainer{
					{Name: "name", Email: "name@example.com"},
				},
				Company:     "company",
				Website:     "http://example.com",
				Source:      "https://git.example.com/repo",
				License:     "license",
				Description: "description",
			},
			expected: nil,
		},
		{
			name: "app missing fields",
			app:  App{},
			expected: ValidationErrors{
				errors.New("app.title: Missing required field"),
				errors.New("app.version: Missing required field"),
				errors.New("app.company: Missing required field"),
				errors.New("app.website: URL cannot be empty"),
				errors.New("app.source: URL cannot be empty"),
				errors.New("app.license: Missing required field"),
				errors.New("app.description: Missing required field"),
				errors.New("app.maintainers: At least one maintainer must be specified"),
			},
		},
		{
			name: "maintainer missing fields",
			app: App{
				Title:   "title",
				Version: "version",
				Maintainers: []Maintainer{
					{Name: "", Email: ""},
				},
				Company:     "company",
				Website:     "http://example.com",
				Source:      "https://git.example.com/repo",
				License:     "license",
				Description: "description",
			},
			expected: ValidationErrors{
				errors.New("maintainer.name: Missing required field"),
				errors.New("maintainer.email: Email address cannot be empty"),
			},
		},
		{
			name: "app invalid URLs",
			app: App{
				Title:   "title",
				Version: "version",
				Maintainers: []Maintainer{
					{Name: "name", Email: "name@example.com"},
				},
				Company:     "company",
				Website:     "http://    invalid url",
				Source:      "http://    invalid url",
				License:     "license",
				Description: "description",
			},
			expected: ValidationErrors{
				errors.New("app.website: Invalid URL"),
				errors.New("app.source: Invalid URL"),
			},
		},
		{
			name: "maintainer invalid email",
			app: App{
				Title:   "title",
				Version: "version",
				Maintainers: []Maintainer{
					{Name: "name", Email: "invalid.com"},
				},
				Company:     "company",
				Website:     "http://example.com",
				Source:      "https://git.example.com/repo",
				License:     "license",
				Description: "description",
			},
			expected: ValidationErrors{
				errors.New("maintainer.email: Invalid email address"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs ValidationErrors
			tc.app.Validate(&errs)
			assert.Equal(t, tc.expected, errs)
		})
	}
}
