package example

import (
	"embed"
	"fmt"
	"github.com/sinlov-go/sembed"
	"io/fs"
	"os"
)

// Example Filesystem:
//
//	example/fixtures
//	├── dir1
//	│   └── onefile.txt
//	└── dir2
//	  └── inner
//	      ├── deeper
//	      │   └── three.txt
//	      ├── one.txt
//	      └── two.txt
//
//go:embed fixtures
var fixtures embed.FS

// To use this lib package, just
// the SetOutput function when your application starts.
func Example() {
	root, _ := sembed.FS(fixtures, "fixtures")

	// Anchor to "fixtures/dir1"
	test1, _ := root.FS("dir1")
	files1, _ := test1.ReadDir(".")

	fmt.Printf("dir1 len %v\n", len(files1))
	fmt.Printf(files1[0].Name()) // "onefile.txt"

	// Anchor to "fixtures/dir2/inner"
	inner, _ := root.FS("dir2/inner")
	one, _ := inner.ReadFile("one.txt")
	fmt.Printf("dir2/inner/one.txt file content: %s\n", one)

	// Fully compatible FS
	_ = fs.WalkDir(inner, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fmt.Printf("Path: %s Name: %s\n", path, d.Name())
		return nil
	})

	// Copy files
	deeper, _ := inner.FS("deeper")
	err := deeper.CopyFile("three.txt", "target.txt", os.FileMode(0o644), true)
	if err != nil {
		fmt.Printf("CopyFile error: %v\n", err)
	}

	// Output:
	// dir1 len 1
	// onefile.txtdir2/inner/one.txt file content: 1
	// Path: . Name: inner
	// Path: deeper Name: deeper
	// Path: deeper/three.txt Name: three.txt
	// Path: one.txt Name: one.txt
	// Path: two.txt Name: two.txt
}
