package vfsutil

import (
	"bytes"
	"encoding/gob"
	"testing"
)

func TestGobFS(t *testing.T) {
	var buf bytes.Buffer
	var gfs GobFS

	err := gob.NewEncoder(&buf).Encode(NewGob(fs2))
	if err != nil {
		t.Fatalf("Got unexpected error when encoding: %v", err)
	}

	err = gob.NewDecoder(&buf).Decode(&gfs)
	if err != nil {
		t.Fatalf("Got unexpected error when decoding: %v", err)
	}
}
