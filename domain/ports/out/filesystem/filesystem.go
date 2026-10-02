package filesystem

type Entry struct {
	Name  string
	IsDir bool
}

type Snapshot struct {
	Path    string
	Base    string
	IsDir   bool
	Entries []Entry
}

type Inspector interface {
	Inspect(path string) (Snapshot, error)
}
