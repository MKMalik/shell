package utils

import (
	"log"
	"os"
	"strings"
)

func ScanPath(paths, arg string) (*string, *string) {
	entries := strings.Split(paths, ":")
	for i := range entries {
		dirEntries, err := os.ReadDir(entries[i])
		if err != nil {
			// panic(err)
			continue
		}
		for _, dirEntry := range dirEntries {
			if dirEntry.IsDir() {
				ScanPath(entries[i]+"/"+dirEntry.Name(), arg)
			} else {
				info, err := dirEntry.Info()
				if err != nil {
					log.Printf("Error getting info for %s: %v", dirEntry.Name(), err)
					continue
				}
				var isExec bool = info.Mode().Perm()&0100 != 0
				if isExec && arg == dirEntry.Name() {
					found := dirEntry.Name()
					full := entries[i] + "/" + found
					return &found, &full
				}
			}
		}
	}
	return nil, nil
}
