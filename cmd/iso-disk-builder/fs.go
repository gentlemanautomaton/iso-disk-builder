package main

import (
	"io/fs"
	"path/filepath"
)

func calculateFileTreeSize(root string) (int64, error) {
	var size int64
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				return err
			}
			if info.Mode().IsRegular() {
				size += info.Size()
			}
		}
		return nil
	})
	return size, err
}
