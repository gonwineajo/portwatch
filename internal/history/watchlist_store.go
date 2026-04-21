package history

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// SaveWatchlist persists a Watchlist to a JSON file at the given path.
func SaveWatchlist(path string, wl *Watchlist) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(wl.Entries())
}

// LoadWatchlist reads a Watchlist from a JSON file.
// Returns an empty Watchlist if the file does not exist.
func LoadWatchlist(path string) (*Watchlist, error) {
	wl := NewWatchlist()
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return wl, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []WatchlistEntry
	if err := json.NewDecoder(f).Decode(&entries); err != nil {
		return nil, err
	}
	for _, e := range entries {
		wl.Add(e.Host, e.Port, e.Label)
	}
	return wl, nil
}
