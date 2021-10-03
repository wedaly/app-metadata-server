package server

import (
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

// Handler constructs an http.Handler for the server.
func Handler() http.Handler {
	var store Store
	r := mux.NewRouter()
	r.HandleFunc("/apps", appMetaPostHandlerFunc(&store)).Methods(http.MethodPost)
	r.HandleFunc("/apps", appMetaGetHandlerFunc(&store)).Methods(http.MethodGet)
	return r
}

// appMetaPostHandlerFunc constructs a handler that writes an app metadata record to storage.
// If multiple YAML documents are submitted in a single request, only the first will be processed.
func appMetaPostHandlerFunc(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/text")

		var app App
		if err := yaml.NewDecoder(r.Body).Decode(&app); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, `invalid YAML`)
			return
		}

		var errs ValidationErrors
		app.Validate(&errs)
		if len(errs) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, errs.Error())
			return
		}

		store.Insert(app)

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `OK`)
	}
}

// appMetaGetHandlerFunc constructs a handler that searches storage for matching app metadata records.
func appMetaGetHandlerFunc(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		w.WriteHeader(http.StatusOK)

		encoder := yaml.NewEncoder(w)
		defer encoder.Close()

		matcher := matcherFromParams(r.URL.Query())
		store.Search(matcher, func(app App) {
			if err := encoder.Encode(&app); err != nil {
				// YAML serialization should always succeed,
				// so this is very likely an error from http.ResponseWriter.
				log.Printf("Error writing YAML response: %s\n", err)
			}
		})
	}
}

// matcherFromParams constructs a matcher function from the request's URL query params.
func matcherFromParams(params url.Values) Matcher {
	m := MatchAny

	if s := params.Get("title"); s != "" {
		m = m.And(MatchExactTitle(s))
	}

	if s := params.Get("version"); s != "" {
		m = m.And(MatchExactVersion(s))
	}

	if s := params.Get("descriptionContains"); s != "" {
		m = m.And(MatchDescriptionContains(s))
	}

	// This can be easily extended to support other filters such as:
	// * Exact or substring match on any of the other app fields (company, website, etc.)
	// * Case-insensitive string comparison.
	// * Version is compatible with a semver range.
	// * Has at least one maintainer with a given name or email address.

	return m
}
