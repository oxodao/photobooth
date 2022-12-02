package utils

import (
	"os"
	"path/filepath"

	"github.com/oxodao/photomaton/config"
)

func GetPath(path string) string {
	return filepath.Join(config.GET.RootPath, path)
}

func MakeOrCreateFolder(path string) error {
	if _, err := os.Stat(GetPath(path)); os.IsNotExist(err) {
		err := os.MkdirAll(GetPath(path), os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}
