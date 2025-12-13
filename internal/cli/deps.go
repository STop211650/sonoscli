package cli

import "github.com/steipete/sonoscli/internal/sonos"

// Dependency injection points for tests.
var newSMAPITokenStore = func() (sonos.SMAPITokenStore, error) {
	return sonos.NewDefaultSMAPITokenStore()
}
