package filesystem

import (
	"GameManager/domain/ports/out/filesystem"
	"os"
	"path/filepath"
)

type Inspector struct{}

func NewInspector() *Inspector {
	return &Inspector{}
}

func (Inspector) Inspect(path string) (filesystem.Snapshot, error) {
	info, err := os.Stat(path)
	if err != nil {
		return filesystem.Snapshot{}, err
	}
	snapshot := filesystem.Snapshot{Path: path, Base: filepath.Base(path), IsDir: info.IsDir()}
	if !info.IsDir() {
		return snapshot, nil
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return filesystem.Snapshot{}, err
	}
	snapshot.Entries = make([]filesystem.Entry, 0, len(entries))
	for _, entry := range entries {
		snapshot.Entries = append(snapshot.Entries, filesystem.Entry{Name: entry.Name(), IsDir: entry.IsDir()})
	}
	return snapshot, nil
}
