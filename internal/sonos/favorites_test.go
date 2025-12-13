package sonos

import "testing"

func TestFavoriteURI(t *testing.T) {
	t.Parallel()

	t.Run("direct", func(t *testing.T) {
		if got := favoriteURI(DIDLItem{URI: "x://1"}); got != "x://1" {
			t.Fatalf("got %q", got)
		}
	})

	t.Run("fromResMD", func(t *testing.T) {
		f := DIDLItem{
			ResMD: `<DIDL-Lite xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns="urn:schemas-upnp-org:metadata-1-0/DIDL-Lite/">
  <item id="R:0/0/0" parentID="R:0/0" restricted="true">
    <dc:title>X</dc:title>
    <res>http://example.com/stream</res>
  </item>
</DIDL-Lite>`,
		}
		if got := favoriteURI(f); got != "http://example.com/stream" {
			t.Fatalf("got %q", got)
		}
	})

	t.Run("none", func(t *testing.T) {
		if got := favoriteURI(DIDLItem{}); got != "" {
			t.Fatalf("got %q", got)
		}
	})
}
