package sonos

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFileSMAPITokenStore_SaveLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tokens.json")
	store, err := NewFileSMAPITokenStore(path)
	if err != nil {
		t.Fatalf("NewFileSMAPITokenStore: %v", err)
	}

	pair := SMAPITokenPair{
		AuthToken:  "token",
		PrivateKey: "key",
		UpdatedAt:  time.Now().UTC().Truncate(time.Second),
	}
	if err := store.Save("9", "Sonos_ABC", pair); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, ok, err := store.Load("9", "Sonos_ABC")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !ok {
		t.Fatalf("expected token to exist")
	}
	if got.AuthToken != "token" || got.PrivateKey != "key" {
		t.Fatalf("unexpected token pair: %#v", got)
	}

	fi, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if fi.Mode().Perm() != 0o600 {
		t.Fatalf("expected perms 0600, got %o", fi.Mode().Perm())
	}
}
