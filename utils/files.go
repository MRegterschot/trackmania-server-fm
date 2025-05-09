package utils

import (
	"io/fs"
	"path/filepath"
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