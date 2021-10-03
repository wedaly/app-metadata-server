package server

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func withTestServer(f func(url string)) {
	server := httptest.NewServer(Handler())
	defer server.Close()
	f(server.URL)
}

func loadTestData(t *testing.T, filename string) io.Reader {
	path := filepath.Join("testdata", filename)
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	return bytes.NewReader(data)
}

func postApp(t *testing.T, url string, data io.Reader) *http.Response {
	resp, err := http.Post(url+"/apps", "application/yaml", data)
	require.NoError(t, err)
	return resp
}

func assertResponseBody(t *testing.T, resp *http.Response, expected string) {
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, expected, string(body))
}

func TestServerInsertValid(t *testing.T) {
	withTestServer(func(url string) {
		validYaml := loadTestData(t, "valid.yaml")
		resp := postApp(t, url, validYaml)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestServerInsertEmpty(t *testing.T) {
	withTestServer(func(url string) {
		emptyYaml := strings.NewReader("")
		resp := postApp(t, url, emptyYaml)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assertResponseBody(t, resp, "invalid YAML")
	})
}

func TestServerInsertNonAsciiUnicode(t *testing.T) {
	withTestServer(func(url string) {
		unicodeYaml := loadTestData(t, "unicode.yaml")
		resp := postApp(t, url, unicodeYaml)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestServerInsertInvalidYaml(t *testing.T) {
	withTestServer(func(url string) {
		invalidYaml := strings.NewReader("{")
		resp := postApp(t, url, invalidYaml)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assertResponseBody(t, resp, "invalid YAML")
	})
}

func TestServerInsertInvalidApp(t *testing.T) {
	withTestServer(func(url string) {
		invalidYaml := loadTestData(t, "invalid.yaml")
		resp := postApp(t, url, invalidYaml)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		expectedBody := strings.Join([]string{
			"app.license: Missing required field",
			"maintainer.email: Invalid email address",
		}, "\n") + "\n"
		assertResponseBody(t, resp, expectedBody)
	})
}

func TestServerSearch(t *testing.T) {
	buildValidApp := func(title, version, description string) App {
		return App{
			Title:       title,
			Version:     version,
			Description: description,
			Maintainers: []Maintainer{
				{Name: "name", Email: "name@example.com"},
			},
			Company: "company",
			Website: "http://example.com",
			Source:  "https://git.example.com/repo",
			License: "license",
		}
	}

	testCases := []struct {
		name     string
		apps     []App
		params   url.Values
		expected []App
	}{
		{
			name: "no parameters retrieves everything",
			apps: []App{
				buildValidApp("foo", "0.0.1", "foo v0.0.1"),
				buildValidApp("foo", "0.0.2", "foo v0.0.2"),
				buildValidApp("bar", "1.2.3", "bar v1.2.3"),
			},
			params: url.Values{},
			expected: []App{
				buildValidApp("foo", "0.0.1", "foo v0.0.1"),
				buildValidApp("foo", "0.0.2", "foo v0.0.2"),
				buildValidApp("bar", "1.2.3", "bar v1.2.3"),
			},
		},
		{
			name: "filter by exact title match",
			apps: []App{
				buildValidApp("foo", "0.0.1", "foo v0.0.1"),
				buildValidApp("foo", "0.0.2", "foo v0.0.2"),
				buildValidApp("bar", "1.2.3", "bar v1.2.3"),
				buildValidApp("bar", "4.5.6", "bar v4.5.6"),
				buildValidApp("baz", "0.0.0.", "baz v0.0.0"),
			},
			params: url.Values{
				"title": []string{"bar"},
			},
			expected: []App{
				buildValidApp("bar", "1.2.3", "bar v1.2.3"),
				buildValidApp("bar", "4.5.6", "bar v4.5.6"),
			},
		},
		{
			name: "filter by exact version match",
			apps: []App{
				buildValidApp("foo", "0.0.1", "foo v0.0.1"),
				buildValidApp("foo", "0.0.2", "foo v0.0.2"),
				buildValidApp("bar", "0.0.1", "bar v0.0.1"),
			},
			params: url.Values{
				"version": []string{"0.0.1"},
			},
			expected: []App{
				buildValidApp("foo", "0.0.1", "foo v0.0.1"),
				buildValidApp("bar", "0.0.1", "bar v0.0.1"),
			},
		},
		{
			name: "filter by exact title AND version match",
			apps: []App{
				buildValidApp("foo", "0.0.1", "foo v0.0.1"),
				buildValidApp("foo", "0.0.2", "foo v0.0.2"),
				buildValidApp("bar", "0.0.1", "bar v0.0.1"),
			},
			params: url.Values{
				"title":   []string{"foo"},
				"version": []string{"0.0.1"},
			},
			expected: []App{
				buildValidApp("foo", "0.0.1", "foo v0.0.1"),
			},
		},
		{
			name: "filter by description contains",
			apps: []App{
				buildValidApp("foobar", "0", "foobar"),
				buildValidApp("baz", "1", "baz"),
				buildValidApp("foobat", "2", "foobat"),
				buildValidApp("bazbat", "3", "bazbat"),
			},
			params: url.Values{
				"descriptionContains": []string{"foo"},
			},
			expected: []App{
				buildValidApp("foobar", "0", "foobar"),
				buildValidApp("foobat", "2", "foobat"),
			},
		},
		{
			name: "filter no matches",
			apps: []App{
				buildValidApp("foo", "0.0.1", "foo v0.0.1"),
			},
			params: url.Values{
				"title": []string{"nomatch"},
			},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			withTestServer(func(url string) {
				// Insert apps for this test case.
				for _, app := range tc.apps {
					postData, err := yaml.Marshal(&app)
					require.NoError(t, err)
					resp := postApp(t, url, bytes.NewReader(postData))
					require.Equal(t, http.StatusOK, resp.StatusCode)
				}

				// Search using the query params for this test case.
				searchUrl := url + "/apps?" + tc.params.Encode()
				res, err := http.Get(searchUrl)
				require.NoError(t, err)
				assert.Equal(t, http.StatusOK, res.StatusCode)

				// Parse results back into apps.
				var results []App
				decoder := yaml.NewDecoder(res.Body)
				for {
					var app App
					err = decoder.Decode(&app)
					if err == io.EOF {
						break
					} else if err != nil {
						require.NoError(t, err)
					}
					results = append(results, app)
				}

				// Check that the expected apps were returned.
				assert.Equal(t, tc.expected, results)
			})
		})
	}
}
