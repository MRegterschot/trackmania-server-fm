package utils

import (
	"io/fs"
	"path"
	"path/filepath"
	"strings"
)

func GetFilesRecursively(path string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(path, func(currentPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		relativePath, err := filepath.Rel(path, currentPath)
		if err != nil {
			return err
		}

		files = append(files, relativePath)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func GetSizeIfFile(info fs.FileInfo) int64 {
	if info.IsDir() {
		return 0
	}
	return info.Size()
}

func IsProbablyDirectory(p string) bool {
	ext := path.Ext(p)
	return strings.HasSuffix(p, "/") || ext == ""
}
