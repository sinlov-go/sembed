package sembed

import (
	"embed"
	"io/fs"
	"os"
)

type embedFS = embed.FS

// Sembed is an embed.FS compatible wrapper, providing Sub() functionality.
type Sembed interface {
	// Open opens the named file for reading and returns it as an fs.File.
	// name style is path style not filepath style.
	Open(name string) (fs.File, error)

	// ReadDir reads and returns the entire named directory.
	// name style is path style not filepath style.
	ReadDir(name string) ([]fs.DirEntry, error)

	// ReadFile reads and returns the content of the named file.
	// name style is path style not filepath style.
	ReadFile(name string) ([]byte, error)

	// FS returns a new Sembed anchored at the given subdirectory.
	// subDir style is path style not filepath style.
	FS(subDir string) (Sembed, error)

	// CopyFile copies a file from the embedded filesystem to the target path.
	//
	//	sourcePath is the path to the file in the embedded filesystem.
	//	target is the target path where the file will be copied.
	//	perm is the permission mode of the target file. os.FileMode(0o644) or os.FileMode(0o666)
	//	coverage true will coverage old
	CopyFile(sourcePath string, target string, perm os.FileMode, coverage bool) error
}

type innerSebmeb struct {
	embedFS

	basedir string
}

// FS creates an embed.FS compatible struct, anchored to the given basedir.
// basedir style is path style not filepath style.
func FS(fs embed.FS, basedir string) (Sembed, error) {
	result := &innerSebmeb{
		embedFS: fs,
		basedir: basedir,
	}

	_, err := result.ReadDir(".")
	if err != nil {
		return nil, err
	}

	return result, nil
}
