package main

import (
	"flag"
	"testing"
)

func TestSkipGramophileShuffleDefaultFlag(t *testing.T) {
	f := flag.Lookup("skip_gramophile_shuffle")
	if f == nil {
		t.Fatal("expected skip_gramophile_shuffle flag to be registered")
	}
	if f.DefValue != "true" {
		t.Errorf("expected skip_gramophile_shuffle default to be true, got %q", f.DefValue)
	}
}

func TestShouldSkipGramophileShuffle(t *testing.T) {
	tests := []struct {
		name        string
		lines       []string
		skipEnabled bool
		want        bool
	}{
		{
			name:        "exact match on line",
			lines:       []string{"Gramophile Shuffle", "Artist - Album"},
			skipEnabled: true,
			want:        true,
		},
		{
			name:        "substring match within line",
			lines:       []string{"=== Gramophile Shuffle: Daily Pick ===", "Track 1"},
			skipEnabled: true,
			want:        true,
		},
		{
			name:        "case-insensitive uppercase match",
			lines:       []string{"GRAMOPHILE SHUFFLE", "Some Record"},
			skipEnabled: true,
			want:        true,
		},
		{
			name:        "case-insensitive lowercase match on later line",
			lines:       []string{"Header", "queued via gramophile shuffle today"},
			skipEnabled: true,
			want:        true,
		},
		{
			name:        "flag disabled allows printing",
			lines:       []string{"Gramophile Shuffle", "Artist - Album"},
			skipEnabled: false,
			want:        false,
		},
		{
			name:        "non-matching job",
			lines:       []string{"Gramophile Wantlist", "Shopping List"},
			skipEnabled: true,
			want:        false,
		},
		{
			name:        "empty lines",
			lines:       nil,
			skipEnabled: true,
			want:        false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := shouldSkipGramophileShuffle(tc.lines, tc.skipEnabled)
			if got != tc.want {
				t.Errorf("shouldSkipGramophileShuffle(%q, %v) = %v, want %v", tc.lines, tc.skipEnabled, got, tc.want)
			}
		})
	}
}
