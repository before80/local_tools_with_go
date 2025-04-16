package main

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"log"
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
		{regexp.MustCompile("\n\n```\\s*?\n&zeroWidthSpace;"), "```\n&zeroWidthSpace;"},
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
		{regexp.MustCompile("### 注意"), "**注意**"},
		{regexp.MustCompile("### 注解"), "**注解**"},
		{regexp.MustCompile("### 示例"), "**示例**"},
		{regexp.MustCompile("### 参数"), "**参数**"},
		{regexp.MustCompile("### 缺陷报告"), "**缺陷报告**"},
		{regexp.MustCompile("&zeroWidthSpace;"), "​\t"},
		{regexp.MustCompile(`### ([a-zA-Z_]+)\s*?\(C(\d+)\s*?起\)`), "### $1 <- $2+"},
		{regexp.MustCompile(`### ([a-zA-Z_]+)\s*?<-\s*?\(C(\d+)\s*?起\)`), "### $1 <- $2+"},
		{regexp.MustCompile(`### ([a-zA-Z_]+)\s*?<-\s*?(\d{2}\+)\s*?\(C(\d{2})\s*?移除\)`), "### $1 <- $2 $3 D"},
		{regexp.MustCompile(`### ([a-zA-Z_]+)\s*?<-\s*?(\d{2}\+)\s*?\(C(\d{2})\s*?前\)`), "### $1 <- $2 $3 F"},
		{regexp.MustCompile(`原址：([a-zA-Z0-9_:/?.#=&-]+)`), "原址：[$1]($1)"},
		{regexp.MustCompile(`运行此代码`), ""},
		{regexp.MustCompile("`\\*\\*A\\*\\*`"), "`A`"},
		{regexp.MustCompile("`\\*\\*a\\*\\*`"), "`a`"},
		{regexp.MustCompile("`\\*\\*c\\*\\*`"), "`c`"},
		{regexp.MustCompile("`\\*\\*d\\*\\*`"), "`d`"},
		{regexp.MustCompile("`\\*\\*F\\*\\*`"), "`F`"},
		{regexp.MustCompile("`\\*\\*f\\*\\*`"), "`f`"},
		{regexp.MustCompile("`\\*\\*E\\*\\*`"), "`E`"},
		{regexp.MustCompile("`\\*\\*e\\*\\*`"), "`e`"},
		{regexp.MustCompile("`\\*\\*G\\*\\*`"), "`G`"},
		{regexp.MustCompile("`\\*\\*g\\*\\*`"), "`g`"},
		{regexp.MustCompile("`\\*\\*i\\*\\*`"), "`i`"},
		{regexp.MustCompile("`\\*\\*n\\*\\*`"), "`n`"},
		{regexp.MustCompile("`\\*\\*P\\*\\*`"), "`P`"},
		{regexp.MustCompile("`\\*\\*O\\*\\*`"), "`O`"},
		{regexp.MustCompile("`\\*\\*o\\*\\*`"), "`o`"},
		{regexp.MustCompile("`\\*\\*p\\*\\*`"), "`p`"},
		{regexp.MustCompile("`\\*\\*s\\*\\*`"), "`s`"},
		{regexp.MustCompile("`\\*\\*U\\*\\*`"), "`U`"},
		{regexp.MustCompile("`\\*\\*u\\*\\*`"), "`u`"},
		{regexp.MustCompile("`\\*\\*X\\*\\*`"), "`X`"},
		{regexp.MustCompile("`\\*\\*x\\*\\*`"), "`x`"},
		{regexp.MustCompile("`\\*\\*Z\\*\\*`"), "`Z`"},
		{regexp.MustCompile("`\\*\\*z\\*\\*`"), "`z`"},
		{regexp.MustCompile("`\\*\\*\\+\\*\\*`"), "`+`"},
		{regexp.MustCompile("`\\*\\*-\\*\\*`"), "`-`"},
		{regexp.MustCompile("`\\*\\*%\\*\\*`"), "`%`"},
		{regexp.MustCompile("`\\*\\*\\*\\*\\*`"), "`*`"},
		{regexp.MustCompile("`\\*\\*\\.\\*\\*`"), "`.`"},
		{regexp.MustCompile("`\\*\\*#\\*\\*`"), "`#`"},
		{regexp.MustCompile("`\\*\\*0\\*\\*`"), "`0`"},
		{regexp.MustCompile("`\\*\\*0X\\*\\*`"), "`0X`"},
		{regexp.MustCompile("`\\*\\*0x\\*\\*`"), "`0x`"},
		{regexp.MustCompile("`\\*\\*%p\\*\\*`"), "`%p`"},
		{regexp.MustCompile("`\\*\\*%d\\*\\*`"), "`%d`"},
		{regexp.MustCompile("`\\*\\*%f\\*\\*`"), "`%f`"},
		{regexp.MustCompile("`\\*\\*%%\\*\\*`"), "`%%`"},
		{regexp.MustCompile(`"\*\*\\f\*\*"`), "`\"\\f\"`"},
		{regexp.MustCompile(`'\*\*\\f\*\*'`), "`'\\f'`"},
		{regexp.MustCompile(`"\*\*\\n\*\*"`), "`\"\\n\"`"},
		{regexp.MustCompile(`'\*\*\\n\*\*'`), "`'\\n'`"},
		{regexp.MustCompile(`"\*\*\\r\*\*"`), "`\"\\r\"`"},
		{regexp.MustCompile(`'\*\*\\r\*\*'`), "`'\\r'`"},
		{regexp.MustCompile(`"\*\*\\t\*\*"`), "`\"\\t\"`"},
		{regexp.MustCompile(`'\*\*\\t\*\*'`), "`'\\t'`"},
		{regexp.MustCompile(`"\*\*\\v\*\*"`), "`\"\\v\"`"},
		{regexp.MustCompile(`'\*\*\\v\*\*'`), "`'\\v'`"},
		{regexp.MustCompile("`\\*\\*\\f\\*\\*`"), "`\\f`"},
		{regexp.MustCompile("`\\*\\*\\n\\*\\*`"), "`\\n`"},
		{regexp.MustCompile("`\\*\\*\\r\\*\\*`"), "`\\r`"},
		{regexp.MustCompile("`\\*\\*\\t\\*\\*`"), "`\\t`"},
		{regexp.MustCompile("`\\*\\*\\v\\*\\*`"), "`\\v`"},
		{regexp.MustCompile("`\\*\\*INF\\*\\*`"), "`INF`"},
		{regexp.MustCompile("`\\*\\*INFINITY\\*\\*`"), "`INFINITY`"},
		{regexp.MustCompile("`\\*\\*NAN\\*\\*`"), "`NAN`"},
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

	// 初始化键盘监听
	if err := keyboard.Open(); err != nil {
		log.Fatalf("无法初始化键盘监听: %v", err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	files, err := findMarkdownFiles(dir)
	if err != nil {
		fmt.Printf("查找 .md 文件时出错: %v\n", err)
		return
	}
	for {
		replacedCount := 0
		for _, file := range files {
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
		//fmt.Println("是否再次执行？输入 1 再次执行，输入 2 退出：")
		//labelForJudge:
		//	reader := bufio.NewReader(os.Stdin)
		//	input, err := reader.ReadString('\n')
		//	if err != nil {
		//		fmt.Printf("读取输入时出错: %v\n", err)
		//		return
		//	}
		//	input = strings.TrimSpace(input)
		//	if input == "2" {
		//		break
		//	} else if input != "1" {
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
