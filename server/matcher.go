package server

import "strings"

// Matcher is a function that determines whether an app record should be included in search results.
type Matcher func(App) bool

// And produces a new matcher that returns true only when both input matchers return true.
func (m Matcher) And(next Matcher) Matcher {
	return func(app App) bool {
		return m(app) && next(app)
	}
}

// MatchAny matches every app metadata record.
var MatchAny Matcher = func(App) bool {
	return true
}

// MatchExactTitle matches apps that have the exact title (case-sensitive).
func MatchExactTitle(s string) Matcher {
	return func(app App) bool {
		return app.Title == s
	}
}

// MatchExactVersion matches apps that have the exact version (case-sensitive).
func MatchExactVersion(s string) Matcher {
	return func(app App) bool {
		return app.Version == s
	}
}

// MatchDescriptionContains matches apps whose description contains the specified string (case-sensitive).
func MatchDescriptionContains(s string) Matcher {
	return func(app App) bool {
		return strings.Contains(app.Description, s)
	}
}
