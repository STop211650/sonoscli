package spotify

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestSearch_Tracks_ClientCredentials(t *testing.T) {
	t.Parallel()

	var tokenCalls int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/token":
			atomic.AddInt32(&tokenCalls, 1)
			auth := r.Header.Get("Authorization")
			want := "Basic " + base64.StdEncoding.EncodeToString([]byte("id:secret"))
			if auth != want {
				t.Fatalf("unexpected auth header: got %q want %q", auth, want)
			}
			if ct := r.Header.Get("Content-Type"); !strings.Contains(ct, "application/x-www-form-urlencoded") {
				t.Fatalf("unexpected content-type: %q", ct)
			}
			_ = r.ParseForm()
			if r.Form.Get("grant_type") != "client_credentials" {
				t.Fatalf("unexpected grant_type: %q", r.Form.Get("grant_type"))
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"access_token": "tok",
				"token_type":   "Bearer",
				"expires_in":   3600,
			})
		case "/v1/search":
			if r.Header.Get("Authorization") != "Bearer tok" {
				t.Fatalf("unexpected bearer token: %q", r.Header.Get("Authorization"))
			}
			q, _ := url.ParseQuery(r.URL.RawQuery)
			if q.Get("type") != "track" {
				t.Fatalf("unexpected type: %q", q.Get("type"))
			}
			if q.Get("q") != "hello world" {
				t.Fatalf("unexpected q: %q", q.Get("q"))
			}
			if q.Get("limit") != "2" {
				t.Fatalf("unexpected limit: %q", q.Get("limit"))
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"tracks": map[string]any{
					"items": []any{
						map[string]any{
							"id":   "t1",
							"name": "Song 1",
							"uri":  "spotify:track:t1",
							"external_urls": map[string]any{
								"spotify": "https://open.spotify.com/track/t1",
							},
							"artists": []any{
								map[string]any{"name": "Artist A"},
							},
							"album": map[string]any{"name": "Album X"},
						},
					},
				},
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(srv.Close)

	c := New("id", "secret", &http.Client{Timeout: 2 * time.Second})
	c.AccountsBaseURL = srv.URL
	c.APIBaseURL = srv.URL

	ctx := context.Background()
	results, err := c.Search(ctx, "hello world", TypeTrack, 2, "")
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("unexpected results length: %d", len(results))
	}
	if results[0].URI != "spotify:track:t1" || results[0].ID != "t1" || results[0].Title != "Song 1" {
		t.Fatalf("unexpected result: %+v", results[0])
	}
	if results[0].Subtitle != "Artist A — Album X" {
		t.Fatalf("unexpected subtitle: %q", results[0].Subtitle)
	}

	// Second search should reuse cached token.
	_, err = c.Search(ctx, "hello world", TypeTrack, 2, "")
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if got := atomic.LoadInt32(&tokenCalls); got != 1 {
		t.Fatalf("expected 1 token call, got %d", got)
	}
}
