package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	err := processFilesInCurrentDir()
	if err != nil {
		fmt.Println("Error:", err)
	}
	time.Sleep(999 * time.Second)
}

func processFilesInCurrentDir() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	err = filepath.Walk(currentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if info.Name() == "_index.md" {
			fmt.Println("current file path = ", path)
			err = processFile(path)
			if err != nil {
				fmt.Printf("Error processing file %s: %v\n", path, err)
			}
		}
		return nil
	})

	return err
}

func processFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(lines) < 11 && strings.Contains(line, "linkTitle") {
			continue
		}
		if len(lines) < 11 && strings.Contains(line, "[menu.main]") {
			continue
		}

		if len(lines) < 11 && strings.Contains(line, "weight") {
			line = strings.TrimSpace(line)
		}

		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	file.Close()

	err = os.Remove(filePath)
	if err != nil {
		return err
	}

	newFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer newFile.Close()

	writer := bufio.NewWriter(newFile)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	writer.Flush()

	return nil
}
