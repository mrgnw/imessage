package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTightenFileModes(t *testing.T) {
	dir := t.TempDir()
	dataDir := filepath.Join(dir, "corten-matrix")
	if err := os.Mkdir(dataDir, 0755); err != nil {
		t.Fatal(err)
	}
	dbPath := filepath.Join(dataDir, "corten-matrix.db")
	for _, p := range []string{dbPath, dbPath + "-wal"} {
		if err := os.WriteFile(p, nil, 0644); err != nil {
			t.Fatal(err)
		}
	}

	tightenFileModes("sqlite3-fk-wal", "file:"+dbPath+"?_txlock=immediate", dataDir)

	wantModes := map[string]os.FileMode{dataDir: 0700, dbPath: 0600, dbPath + "-wal": 0600}
	for p, want := range wantModes {
		info, err := os.Stat(p)
		if err != nil {
			t.Fatal(err)
		}
		if got := info.Mode().Perm(); got != want {
			t.Fatalf("%s mode = %o, want %o", filepath.Base(p), got, want)
		}
	}
	if _, err := os.Stat(dbPath + "-shm"); !os.IsNotExist(err) {
		t.Fatalf("missing -shm must stay missing, stat err = %v", err)
	}
}

func TestTightenFileModesIgnoresNonFileDatabases(t *testing.T) {
	dir := t.TempDir()
	otherFile := filepath.Join(dir, "bridge.db")
	if err := os.WriteFile(otherFile, nil, 0644); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ dbType, uri string }{
		{"postgres", "postgres://user@localhost/" + otherFile},
		{"sqlite3-fk-wal", ":memory:"},
		{"sqlite3-fk-wal", "file::memory:?cache=shared"},
		{"sqlite3-fk-wal", "file:" + otherFile + "?mode=memory"},
		{"sqlite3-fk-wal", ""},
		{"sqlite3-fk-wal", "file:" + filepath.Join(dir, "missing.db")},
	} {
		tightenFileModes(tc.dbType, tc.uri, filepath.Join(dir, "missing-data-dir"))
	}
	info, err := os.Stat(otherFile)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0644 {
		t.Fatalf("non-file DSN touched %s: mode = %o", otherFile, got)
	}
}
