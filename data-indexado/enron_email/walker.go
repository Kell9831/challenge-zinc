package enron_email

import (
	"os"
	"path/filepath"
)

func Walk(root string, files chan string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files <- path
		}
		return nil
	})
}
