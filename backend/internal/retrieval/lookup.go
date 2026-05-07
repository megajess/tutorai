package retrieval

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ColorLookup resolves guild, shard, and wedge names to color identity slices.
// It is loaded once at startup from data/color_identity_lookup.json.
type ColorLookup struct {
	table map[string][]string
}

// LoadColorLookup reads the color identity lookup JSON file and returns a
// ColorLookup ready for use. The file is read once; subsequent lookups are
// in-memory map reads.
func LoadColorLookup(path string) (*ColorLookup, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read color lookup file: %w", err)
	}

	var table map[string][]string
	if err := json.Unmarshal(data, &table); err != nil {
		return nil, fmt.Errorf("parse color lookup file: %w", err)
	}

	return &ColorLookup{table: table}, nil
}

// Resolve returns the color identity slice for the given term, or nil if the
// term is not a known guild, shard, or wedge name (including aliases).
// Lookup is case-insensitive.
func (cl *ColorLookup) Resolve(term string) []string {
	colors, ok := cl.table[strings.ToLower(strings.TrimSpace(term))]
	if !ok {
		return nil
	}
	return colors
}
