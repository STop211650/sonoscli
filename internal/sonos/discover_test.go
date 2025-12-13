package sonos

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestPreferDeviceSet(t *testing.T) {
	best := map[string]Device{
		"1.1.1.1": {IP: "1.1.1.1", Name: "A"},
	}

	// Smaller candidate should be ignored.
	best2 := preferDeviceSet(best, map[string]Device{})
	if len(best2) != 1 {
		t.Fatalf("expected best unchanged, got %d", len(best2))
	}

	// Larger candidate should replace.
	candidateLarger := map[string]Device{
		"2.2.2.2": {IP: "2.2.2.2", Name: "B"},
		"3.3.3.3": {IP: "3.3.3.3", Name: "C"},
	}
	best3 := preferDeviceSet(best, candidateLarger)
	if len(best3) != 2 || best3["2.2.2.2"].Name != "B" {
		t.Fatalf("expected replace with candidate, got %#v", best3)
	}

	// Equal-size candidate should merge missing keys.
	bestEqual := map[string]Device{
		"10.0.0.1": {IP: "10.0.0.1", Name: "X"},
		"10.0.0.2": {IP: "10.0.0.2", Name: "Y"},
	}
	candidateEqual := map[string]Device{
		"10.0.0.2": {IP: "10.0.0.2", Name: "Y2"}, // existing key should not overwrite
		"10.0.0.3": {IP: "10.0.0.3", Name: "Z"},
	}
	merged := preferDeviceSet(bestEqual, candidateEqual)
	if len(merged) != 3 {
		t.Fatalf("expected merge size 3, got %d", len(merged))
	}
	if merged["10.0.0.2"].Name != "Y" {
		t.Fatalf("expected existing key preserved, got %q", merged["10.0.0.2"].Name)
	}
	if merged["10.0.0.3"].Name != "Z" {
		t.Fatalf("expected new key added, got %#v", merged["10.0.0.3"])
	}
}

func TestDiscoverFallsBackWhenSSDPDeadlineExceeded(t *testing.T) {
	origSSDP := ssdpDiscoverFunc
	origScan := scanAnySpeakerIPFunc
	origTop := discoverViaTopologyFunc
	origTopFromIP := discoverViaTopologyFromIPFunc
	t.Cleanup(func() {
		ssdpDiscoverFunc = origSSDP
		scanAnySpeakerIPFunc = origScan
		discoverViaTopologyFunc = origTop
		discoverViaTopologyFromIPFunc = origTopFromIP
	})

	ssdpDiscoverFunc = func(ctx context.Context, timeout time.Duration) ([]ssdpResult, error) {
		return nil, context.DeadlineExceeded
	}
	discoverViaTopologyFunc = func(ctx context.Context, timeout time.Duration, results []ssdpResult, includeInvisible bool) ([]Device, error) {
		return nil, errors.New("no ssdp candidates")
	}
	scanAnySpeakerIPFunc = func(ctx context.Context, timeout time.Duration) (string, error) {
		return "192.168.1.10", nil
	}
	discoverViaTopologyFromIPFunc = func(ctx context.Context, timeout time.Duration, ip string, includeInvisible bool) ([]Device, error) {
		return []Device{{IP: ip, Name: "Office", UDN: "RINCON_x"}}, nil
	}

	devs, err := Discover(context.Background(), DiscoverOptions{Timeout: 3 * time.Second})
	if err != nil {
		t.Fatalf("Discover: %v", err)
	}
	if len(devs) != 1 || devs[0].Name != "Office" {
		t.Fatalf("unexpected devices: %#v", devs)
	}
}
