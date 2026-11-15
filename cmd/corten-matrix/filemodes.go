package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"maunium.net/go/mautrix/bridgev2/matrix/mxmain"
)

// tightenDataFileModes makes the bridge DB (and its WAL/SHM side files) and
// the corten-matrix data directory private to the owning user. SQLite creates
// the DB under the process umask (usually 0644) and copies the DB's mode onto
// new -wal/-shm files, so fixing the DB once keeps the side files private too.
func tightenDataFileModes(br *mxmain.BridgeMain) {
	if br.Config == nil {
		return
	}
	dataDir, err := cortenDataDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[filemodes] cannot determine data dir: %v\n", err)
		dataDir = ""
	}
	tightenFileModes(br.Config.Database.Type, br.Config.Database.URI, dataDir)
}

func tightenFileModes(dbType, dbURI, dataDir string) {
	if dataDir != "" {
		chmodIfExists(dataDir, 0700)
	}
	if !isSQLiteDatabaseType(dbType) {
		return
	}
	dbPath := sqliteFilePath(dbURI)
	if dbPath == "" {
		return
	}
	for _, p := range []string{dbPath, dbPath + "-wal", dbPath + "-shm"} {
		chmodIfExists(p, 0600)
	}
}

// sqliteFilePath returns the on-disk path of a SQLite DSN, or "" for an
// in-memory or empty DSN.
func sqliteFilePath(uri string) string {
	base, query, _ := strings.Cut(strings.TrimSpace(uri), "?")
	base = strings.TrimPrefix(base, "file:")
	if base == "" || strings.HasPrefix(base, ":memory:") || strings.Contains(query, "mode=memory") {
		return ""
	}
	return base
}

func chmodIfExists(path string, mode os.FileMode) {
	info, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "[filemodes] stat %s: %v\n", path, err)
		}
		return
	}
	if info.Mode().Perm() == mode {
		return
	}
	if err = os.Chmod(path, mode); err != nil {
		fmt.Fprintf(os.Stderr, "[filemodes] chmod %s: %v\n", path, err)
	}
}

func cortenDataDir() (string, error) {
	dataDir := os.Getenv("XDG_DATA_HOME")
	if dataDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dataDir = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(dataDir, "corten-matrix"), nil
}
