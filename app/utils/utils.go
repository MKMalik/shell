package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func ScanPath(paths, arg string) (*string, *string) {
	for _, dir := range strings.Split(paths, ":") {
		full := filepath.Join(dir, arg)

		info, err := os.Stat(full)
		if err != nil || info.IsDir() {
			continue
		}

		if info.Mode()&0111 != 0 {
			found := arg
			return &found, &full
		}
	}

	return nil, nil
}

func IsDirExists(path string) bool {
	info, err := os.Lstat(path)
	return err == nil && info.IsDir()
}
