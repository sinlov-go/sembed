package sembed

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func (s *innerSebmeb) calculatePath(src string) string {
	base := filepath.Join(s.basedir, src)

	return filepath.ToSlash(base)
}

// Open opens the named file for reading and returns it as an fs.File.
// name style is path style not filepath style.
func (s *innerSebmeb) Open(name string) (fs.File, error) {
	path := s.calculatePath(name)

	return s.embedFS.Open(path)
}

// ReadDir reads and returns the entire named directory.
// name style is path style not filepath style.
func (s *innerSebmeb) ReadDir(name string) ([]fs.DirEntry, error) {
	path := s.calculatePath(name)

	return s.embedFS.ReadDir(path)
}

// ReadFile reads and returns the content of the named file.
// name style is path style not filepath style.
func (s *innerSebmeb) ReadFile(name string) ([]byte, error) {
	path := s.calculatePath(name)

	return s.embedFS.ReadFile(path)
}

// FS returns a new Sembed anchored at the given subdirectory.
// subDir style is path style not filepath style.
func (s *innerSebmeb) FS(subDir string) (Sembed, error) {
	path := s.calculatePath(subDir)

	return FS(s.embedFS, path)
}

// CopyFile copies a file from the embedded filesystem to the target path.
//
//	sourcePath is the path to the file in the embedded filesystem.
//	target is the target path where the file will be copied.
//	perm is the permission mode of the target file. os.FileMode(0o644) or os.FileMode(0o666)
//	coverage true will coverage old
func (s *innerSebmeb) CopyFile(
	sourcePath string,
	target string,
	perm os.FileMode,
	coverage bool,
) error {
	sourceFile, err := s.Open(sourcePath)
	if err != nil {
		return err
	}

	if !coverage {
		exists, errExist := pathExists(target)
		if errExist != nil {
			return errExist
		}

		if exists {
			return fmt.Errorf("not coverage, which target path exist %v", target)
		}
	}

	parentPath := filepath.Dir(target)
	if !pathExistsFast(parentPath) {
		errMkParentPath := os.MkdirAll(parentPath, fetchDefaultFolderFileMode())
		if errMkParentPath != nil {
			return fmt.Errorf(
				"can not WriteFileByB at new dir at mode: %v , at parent path: %v, why: %v",
				perm,
				parentPath,
				errMkParentPath,
			)
		}
	}

	targetFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, perm)
	if err != nil {
		return err
	}

	_, err = io.Copy(targetFile, sourceFile)

	return err
}

// pathExists
//
//	path exists
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// pathExistsFast
//
//	path exists fast
func pathExistsFast(path string) bool {
	exists, _ := pathExists(path)

	return exists
}

// fetchDefaultFolderFileMode
// not support umask will use os.FileMode(0o777).
func fetchDefaultFolderFileMode() fs.FileMode {
	switch runtime.GOOS {
	case "windows":
		return os.FileMode(0o766)
	default:
		umaskCode, err := getUmask()
		if err != nil {
			return os.FileMode(0o777)
		}

		if len(umaskCode) > 3 {
			umaskCode = umaskCode[len(umaskCode)-3:]
		}

		umaskOct, errParseUmask := strconv.ParseInt(umaskCode, 8, 64)
		if errParseUmask != nil {
			return os.FileMode(0o777)
		}

		defaultFOlderCode := 0o777
		nowOct := defaultFOlderCode - int(umaskOct)

		return os.FileMode(nowOct)
	}
}

func getUmask() (string, error) {
	cmd := exec.Command("sh", "-c", "umask")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(&out)
	scanner.Split(bufio.ScanWords)

	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text()), nil
	}

	return "", errors.New("no output from umask command")
}
