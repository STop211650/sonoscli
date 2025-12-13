package sonos

import "testing"

func TestParsePresentationMapXML(t *testing.T) {
	raw := []byte(`
<PresentationMap>
  <SearchCategories>
    <Category id="tracks" mappedId="search:track"/>
    <Category id="albums" mappedId="search:album"/>
    <CustomCategory stringId="Blogs" mappedId="SBLG"/>
  </SearchCategories>
</PresentationMap>`)

	m, err := parsePresentationMapXML(raw)
	if err != nil {
		t.Fatalf("parsePresentationMapXML: %v", err)
	}
	if m["tracks"] != "search:track" {
		t.Fatalf("tracks mapping wrong: %q", m["tracks"])
	}
	if m["albums"] != "search:album" {
		t.Fatalf("albums mapping wrong: %q", m["albums"])
	}
	if m["Blogs"] != "SBLG" {
		t.Fatalf("custom mapping wrong: %q", m["Blogs"])
	}
}
