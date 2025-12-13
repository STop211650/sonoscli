package sonos

import "testing"

func TestParseServiceDescriptorListXML(t *testing.T) {
	xmlPayload := `
<Services SchemaVersion="1">
  <Service Id="9" Name="Spotify" Version="1.1" Uri="http://example.com/svc" SecureUri="https://example.com/svc" ContainerType="MService" Capabilities="513">
    <Policy Auth="DeviceLink" PollInterval="30" />
    <Presentation>
      <Strings Version="1" Uri="https://example.com/strings.xml" />
      <PresentationMap Version="2" Uri="https://example.com/pmap.xml" />
    </Presentation>
  </Service>
  <Service Id="163" Name="Spreaker" Version="1.1" Uri="http://example.com/2" SecureUri="https://example.com/2" ContainerType="MService" Capabilities="0">
    <Policy Auth="Anonymous" />
  </Service>
  <Service Id="999" Name="Foo" Version="1.0" Uri="http://example.com/3" SecureUri="https://example.com/3" ContainerType="MService" Capabilities="0">
    <Policy Auth="AppLink" />
    <Manifest Uri="https://example.com/manifest.json" />
  </Service>
</Services>`

	services, err := parseServiceDescriptorListXML(xmlPayload)
	if err != nil {
		t.Fatalf("parseServiceDescriptorListXML: %v", err)
	}
	if len(services) != 3 {
		t.Fatalf("expected 3 services, got %d", len(services))
	}

	if services[0].Name != "Spotify" {
		t.Fatalf("expected Spotify first, got %q", services[0].Name)
	}
	if services[0].Auth != MusicServiceAuthDeviceLink {
		t.Fatalf("expected DeviceLink auth, got %q", services[0].Auth)
	}
	if services[0].ServiceType != "2311" {
		t.Fatalf("expected serviceType 2311, got %q", services[0].ServiceType)
	}
	if services[0].PresentationMapURI == "" {
		t.Fatalf("expected presentationMapUri")
	}

	if services[1].Auth != MusicServiceAuthAnonymous {
		t.Fatalf("expected Anonymous auth, got %q", services[1].Auth)
	}
	if services[1].ServiceType != "41735" {
		t.Fatalf("expected serviceType 41735, got %q", services[1].ServiceType)
	}

	if services[2].Auth != MusicServiceAuthAppLink {
		t.Fatalf("expected AppLink auth, got %q", services[2].Auth)
	}
	if services[2].ManifestURI == "" {
		t.Fatalf("expected manifestUri")
	}
}
