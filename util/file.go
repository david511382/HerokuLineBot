package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func GetFilePath(filename string) (string, error) {
	if strings.HasPrefix(filename, "~") {
		u, err := user.Current()
		if err != nil {
			return filename, err
		}

		filename = strings.Replace(filename, "~", u.HomeDir, 1)
	} else {
		base, err := filepath.Abs("")
		if err != nil {
			return filename, err
		}

		sep := "/"
		base = strings.ReplaceAll(base, "\\", sep)
		filename = strings.ReplaceAll(filename, "\\", sep)

		dirs := strings.Split(filename, sep)
		for i := 0; i < len(dirs); i++ {
			if dirs[i] != ".." {
				filename = sep + strings.Join(dirs[i:], sep)
				break
			}

			base = filepath.Dir(base)
		}

		filename, err = filepath.Abs(base + filename)
		if err != nil {
			return "", err
		}
	}

	return filename, nil
}

// ReadConfig read config from filepath
func ReadFile(filename string) ([]byte, error) {
	filepath, err := GetFilePath(filename)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func MakeFolderOn(folderPath string) error {
	return os.MkdirAll(folderPath, os.ModePerm)
}

func GetRootOf(
	folder string,
) (
	path string,
	resultErr error,
) {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	for filepath.Base(path) != folder {
		p := filepath.Dir(path)
		if path == p {
			resultErr = fmt.Errorf("No Dir %s", folder)
			return
		}
		path = p
	}

	return
}
