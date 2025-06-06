package main

import (
	"bufio"
	"fmt"
	"github.com/eiannone/keyboard"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// findSectionStart 查找指定类型的起始行
func findSectionStart(lines []string, sectionType string) int {
	sectionHeader := ""
	switch sectionType {
	case "f":
		sectionHeader = "## 函数"
	case "m":
		sectionHeader = "## 宏"
	case "t":
		sectionHeader = "## 类型"
	default:
		return -1
	}
	for i, line := range lines {
		if line == sectionHeader {
			return i
		}
	}
	return -1
}

// insertSubmenu 逐个插入新的菜单名
func insertSubmenu(lines []string, start int, newSubmenu string) []string {
	l := len(lines)
	insertIndex := start + 1
	hadExisted := false

	for i := start + 1; i < l; i++ {
		line := lines[i]
		if !strings.HasPrefix(line, "### ") && !strings.HasPrefix(line, "## ") {
			insertIndex = i + 1
			continue
		}

		if strings.HasPrefix(line, "### ") {
			line = strings.TrimPrefix(line, "### ")
			line = strings.TrimSpace(line)
			line = strings.TrimSuffix(line, " <- 99+")
			line = strings.TrimSuffix(line, " <- 11+")
			line = strings.TrimSuffix(line, " <- 95+")
			line = strings.TrimSuffix(line, " <- 17+")
			line = strings.TrimSuffix(line, " <- 23+")
			if newSubmenu == line {
				hadExisted = true
				break
			} else if newSubmenu < line {
				insertIndex = i
				break
			}
			continue
		} else if strings.HasPrefix(line, "## ") {
			insertIndex = i
			break
		}
	}

	if !hadExisted {
		newLines := make([]string, 0, len(lines)+8)
		if insertIndex <= l {
			newLines = append(newLines, lines[:insertIndex]...)
		} else {
			newLines = append(newLines, lines[:insertIndex-1]...)
		}

		newLines = append(newLines, fmt.Sprintf("### %s\n", newSubmenu))
		newLines = append(newLines, fmt.Sprintf("原址：\n"))
		newLines = append(newLines, fmt.Sprintf("```c\n"))
		newLines = append(newLines, fmt.Sprintf("```\n"))
		newLines = append(newLines, fmt.Sprintf("\n"))
		newLines = append(newLines, fmt.Sprintf("\n"))
		newLines = append(newLines, fmt.Sprintf("\n"))
		newLines = append(newLines, fmt.Sprintf("\n"))
		if insertIndex < l {
			newLines = append(newLines, lines[insertIndex:]...)
		}
		return newLines
	} else {
		return lines
	}

}

func insertSubmenus(lines []string, start int, newSubmenus []string) []string {
	sort.Strings(newSubmenus)
	for _, newSubmenu := range newSubmenus {
		lines = insertSubmenu(lines, start, newSubmenu)
	}
	return lines
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("用法: go run main.go <目录路径>")
		return
	}
	dir := os.Args[1]
	// 初始化键盘监听
	if err := keyboard.Open(); err != nil {
		log.Fatalf("无法初始化键盘监听: %v", err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	for {

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("请一次性输入文件名（无需输入 .md 后缀）、插入类型（f: 函数, m: 宏, t: 类型）和要插入的菜单名（用空格分隔），如：your_file f menu1 menu2 ... : ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("读取输入时出错: %v\n", err)
			return
		}
		input = strings.TrimSpace(input)
		parts := strings.Fields(input)
		if len(parts) < 3 {
			fmt.Println("输入参数不足，请提供文件名、插入类型和至少一个菜单名。")
			return
		}

		filename := parts[0] + ".md"
		sectionType := parts[1]
		menuNames := parts[2:]

		filePath := filepath.Join(dir, filename)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			fmt.Printf("文件 %s 不存在\n", filePath)
			return
		}

		file, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("打开文件 %s 时出错: %v\n", filePath, err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		var lines []string
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("读取文件 %s 时出错: %v\n", filePath, err)
			return
		}

		start := findSectionStart(lines, sectionType)
		if start == -1 {
			fmt.Printf("未找到类型为 %s 的起始行\n", sectionType)
			return
		}

		newLines := insertSubmenus(lines, start, menuNames)

		outputFile, err := os.Create(filePath)
		if err != nil {
			fmt.Printf("创建文件 %s 时出错: %v\n", filePath, err)
			return
		}
		defer outputFile.Close()

		writer := bufio.NewWriter(outputFile)
		for _, line := range newLines {
			_, err = writer.WriteString(line + "\n")
			if err != nil {
				fmt.Printf("写入文件 %s 时出错: %v\n", filePath, err)
				return
			}
		}
		writer.Flush()

		fmt.Println("插入成功！")
		//labelForJudge:
		//	fmt.Println("是否再次执行？输入 1 再次执行，输入 2 退出：")
		//	reader1 := bufio.NewReader(os.Stdin)
		//	input1, err := reader1.ReadString('\n')
		//	if err != nil {
		//		fmt.Printf("读取输入时出错: %v\n", err)
		//		return
		//	}
		//	input1 = strings.TrimSpace(input1)
		//	if input1 == "2" {
		//		break
		//	} else if input1 != "1" {
		//		fmt.Println("输入无效，请输入 1 或 2。")
		//		goto labelForJudge
		//	}
	labelForJudge:
		fmt.Println("按 CTRL+2 继续，CTRL+3 退出")
		// 捕获按键事件
		_, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatalf("读取键盘输入时出错: %v", err)
		}

		// 检查是否按下了 Ctrl 键组合
		if key == keyboard.KeyCtrl2 { // Ctrl + 2
			continue
		} else if key == keyboard.KeyCtrl3 { // Ctrl + 3
			fmt.Println("\n退出程序...")
			break // 跳出循环
		} else {
			fmt.Println("输入无效，请按下 Ctrl + 2 继续，或者按下 Ctrl + 3 退出。")
			goto labelForJudge
		}
	}
}
