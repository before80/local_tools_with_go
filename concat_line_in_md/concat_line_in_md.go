package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func main() {
	//reLine, err := regexp.Compile("^> 原文: ")

	// 从命令行参数中获取 -d 参数，表示要处理的文件夹名称
	dir := flag.String("d", "", "要处理的文件夹名称")
	flag.Parse()

	// 获取当前文件夹中的所有markdown文件列表
	files, err := os.ReadDir(*dir)
	if err != nil {
		log.Fatal(fmt.Errorf("读取当前文件夹中的所有markdown文件列表出现错误：%v\n", err))
	}

	// 遍历所有Markdown文件，并打开每个markdown文件，遍历该文件中的所有行
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if filepath.Ext(file.Name()) != ".md" {
			continue
		}
		filePath := filepath.Join(dir, file.Name())
		fs.PathError{}
	}

}
