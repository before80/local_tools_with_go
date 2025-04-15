package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// 遍历指定目录及其子目录，查找所有 .md 文件
func findMarkdownFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// 替换文件中的特定内容
func replaceInFile(filePath string) (bool, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false, err
	}
	replacements := []struct {
		pattern     *regexp.Regexp
		replacement string
	}{
		{regexp.MustCompile("输出：\\s*?\n```\\s*?\n"), "输出：\n\n```txt\n"},
		{regexp.MustCompile("```\\s*?\n#include"), "```c\n#include"},
		{regexp.MustCompile("```\\s*?\ntypedef"), "```c\ntypedef"},
		{regexp.MustCompile("```\\s*?\nvoid"), "```c\nvoid"},
		{regexp.MustCompile("```\\s*?\n#define"), "```c\n#define"},
		{regexp.MustCompile("```\\s*?\nchar"), "```c\nchar"},
		{regexp.MustCompile("```\\s*?\nint"), "```c\nint"},
		{regexp.MustCompile("```\\s*?\nfloat"), "```c\nfloat"},
		{regexp.MustCompile("```\\s*?\ndouble"), "```c\ndouble"},
		{regexp.MustCompile("@!br /!@"), "<br />"},
		{regexp.MustCompile("!@"), ">"},
		{regexp.MustCompile("@!"), "<"},
		{regexp.MustCompile("### 返回值"), "**返回值**"},
		{regexp.MustCompile("### 注解"), "**注解**"},
		{regexp.MustCompile("### 示例"), "**示例**"},
		{regexp.MustCompile("### 参数"), "**参数**"},
		{regexp.MustCompile("&zeroWidthSpace;"), "​\t"},
		{regexp.MustCompile(`### ([a-zA-Z_]+) <- \(C(\d+) 起\)`), "### $1 <- $2+"},
		{regexp.MustCompile(`原址：([a-zA-Z0-9_:/?.#=&-]+)`), "原址：[$1]($1)"},
		{regexp.MustCompile(`运行此代码`), ""},
	}

	modified := false
	newContent := string(content)
	for _, r := range replacements {
		if r.pattern.MatchString(newContent) {
			newContent = r.pattern.ReplaceAllString(newContent, r.replacement)
			if !modified {
				modified = true
			}
		}
	}

	if modified {
		err = os.WriteFile(filePath, []byte(newContent), 0644)
		if err != nil {
			return false, err
		}
	}
	return modified, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("用法: replace_md_of_c_doc_something.exe <目录路径>")
		return
	}
	dir := os.Args[1]
	for {
		files, err := findMarkdownFiles(dir)
		if err != nil {
			fmt.Printf("查找 .md 文件时出错: %v\n", err)
			return
		}
		replacedCount := 0
		for _, file := range files {
			fmt.Println("file=", file)
			replaced, err := replaceInFile(file)
			if err != nil {
				fmt.Printf("处理文件 %s 时出错: %v\n", file, err)
				continue
			}
			if replaced {
				replacedCount++
			}
		}
		fmt.Printf("本次替换了 %d 个文件中的内容。\n", replacedCount)
		fmt.Println("是否再次执行？输入 1 再次执行，输入 2 退出：")
	labelForJudge:
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("读取输入时出错: %v\n", err)
			return
		}
		input = strings.TrimSpace(input)
		if input == "2" {
			break
		} else if input != "1" {
			fmt.Println("输入无效，请输入 1 或 2。")
			goto labelForJudge
		}
	}
}
