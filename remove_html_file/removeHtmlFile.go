package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	err := deleteHTMLFiles(".")
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func deleteHTMLFiles(rootDir string) error {
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.ToLower(filepath.Ext(info.Name())) == ".html" {
			err := os.Remove(path)
			if err != nil {
				return err
			}
			fmt.Println("Deleted:", path)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
