package sonos

import "testing"

func TestParseDIDLItems_Minimal(t *testing.T) {
	t.Parallel()

	xml := `<?xml version="1.0"?>
<DIDL-Lite xmlns="urn:schemas-upnp-org:metadata-1-0/DIDL-Lite/"
  xmlns:dc="http://purl.org/dc/elements/1.1/"
  xmlns:upnp="urn:schemas-upnp-org:metadata-1-0/upnp/">
  <item id="Q:0/1" parentID="Q:0" restricted="true">
    <dc:title>Hello</dc:title>
    <upnp:class>object.item.audioItem.musicTrack</upnp:class>
    <res protocolInfo="x-rincon-playlist:*:*:*">x-sonos-spotify:spotify%3atrack%3a123</res>
  </item>
</DIDL-Lite>`

	items, err := ParseDIDLItems(xml)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("unexpected len: %d", len(items))
	}
	if items[0].Title != "Hello" {
		t.Fatalf("unexpected title: %q", items[0].Title)
	}
	if items[0].URI == "" {
		t.Fatalf("expected uri")
	}
	if items[0].Class != "object.item.audioItem.musicTrack" {
		t.Fatalf("unexpected class: %q", items[0].Class)
	}
	if items[0].ID != "Q:0/1" {
		t.Fatalf("unexpected id: %q", items[0].ID)
	}
}
