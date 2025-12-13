package sonos

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// fetchAndParsePresentationMap downloads and parses a Sonos SMAPI presentation map.
// It returns a map from human-facing category IDs (e.g. "tracks") to the "mappedId"
// values required for SMAPI search() calls.
func fetchAndParsePresentationMap(ctx context.Context, httpClient *http.Client, uri string) (map[string]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("presentation map: http %s", resp.Status)
	}
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, err
	}
	return parsePresentationMapXML(raw)
}

type pmapEnvelope struct {
	SearchCategories *struct {
		Categories       []pmapCategory       `xml:"Category"`
		CustomCategories []pmapCustomCategory `xml:"CustomCategory"`
	} `xml:"SearchCategories"`
}

type pmapCategory struct {
	ID       string `xml:"id,attr"`
	MappedID string `xml:"mappedId,attr"`
	AltID    string `xml:"mappedID,attr"` // some services are inconsistent
}

type pmapCustomCategory struct {
	StringID string `xml:"stringId,attr"`
	MappedID string `xml:"mappedId,attr"`
}

func parsePresentationMapXML(raw []byte) (map[string]string, error) {
	var env pmapEnvelope
	if err := xml.Unmarshal(raw, &env); err != nil {
		return nil, err
	}
	out := map[string]string{}
	if env.SearchCategories == nil {
		return out, nil
	}
	for _, c := range env.SearchCategories.Categories {
		id := strings.TrimSpace(c.ID)
		mapped := strings.TrimSpace(c.MappedID)
		if mapped == "" {
			mapped = strings.TrimSpace(c.AltID)
		}
		if id != "" && mapped != "" {
			out[id] = mapped
		}
	}
	for _, c := range env.SearchCategories.CustomCategories {
		id := strings.TrimSpace(c.StringID)
		mapped := strings.TrimSpace(c.MappedID)
		if id != "" && mapped != "" {
			out[id] = mapped
		}
	}
	return out, nil
}
