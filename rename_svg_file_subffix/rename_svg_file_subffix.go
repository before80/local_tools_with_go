package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	err := renameSvgFileSuffix(".")
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("已完成")
	time.Sleep(999 * time.Second)
}

func renameSvgFileSuffix(rootDir string) error {
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.ToLower(filepath.Ext(info.Name())) == ".svg+xml" {
			newPath := path[:strings.Index(path, "+")]
			err := os.Rename(path, newPath)
			if err != nil {
				return err
			}
			fmt.Println("已经重命名:", path, "->", newPath)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
