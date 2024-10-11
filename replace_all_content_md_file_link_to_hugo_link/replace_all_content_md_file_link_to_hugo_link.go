package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func main() {
	ms := genPageLinkMap(".", 6000)

	//fmt.Println("已完成")
	//for k, v := range ms {
	//	fmt.Println(k, v)
	//}

	err := replaceLink(ms, ".")
	fmt.Println("完成没有锚的链接的替换")
	//err := findLinkWithAnchor(ms, ".")
	err = replaceLinkWithAnchor(ms, ".")
	fmt.Println(err)
	fmt.Println("完成替换带有锚的链接的替换")
	time.Sleep(999 * time.Second)
}

// 生成map,key为页面原本对应的网址（去除最后面的/），value为页面所在路径的 {{< ref "">}} 格式，其添加了可执行文件所在的文件夹的名称作为前缀)
func genPageLinkMap(rootDir string, initSize int) map[string]string {
	m := make(map[string]string, initSize)
	reLine, err := regexp.Compile("^> 原文: ")
	if err != nil {
		log.Fatal(fmt.Errorf("创建正则匹配出现错误：%v\n", err))
	}

	// 获取当前可执行文件的路径
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(fmt.Errorf("获取当前所处的文件夹名称出现错误"))
	}

	//// 获取当前文件夹（即content目录）下的所有第一级目录的名称，作为后续链接的开头
	//var l1DirNames []string
	//err = filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
	//	if err != nil {
	//		return err
	//	}
	//
	//	if d.IsDir() && !strings.Contains(path, string(os.PathSeparator)) && path != rootDir {
	//		fmt.Println("path = ", path)
	//		l1DirNames = append(l1DirNames, d.Name())
	//	}
	//
	//	return nil
	//})
	//
	//if err != nil {
	//	panic("获取第一级目录名称的切片发生错误")
	//}

	// 获取当前文件夹的路径
	currentDir := filepath.Dir(exePath)
	// 使用filepath包获取文件夹名称
	currentExecutableDirName := filepath.Base(currentDir)

	fmt.Println("当前可执行文件所处的文件夹的名称：", currentExecutableDirName)

	reLink, err := regexp.Compile(`\[[^\]]+\]\(([^\)]+)\)`)
	if err != nil {
		log.Fatal(fmt.Errorf("创建正则匹配出现错误：%v\n", err))
	}

	err = filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		fileName := d.Name()
		if strings.ToLower(filepath.Ext(fileName)) == ".md" {
			var noExtPath string
			newMdFilepath := strings.TrimPrefix(path, "./")

			if fileName == "_index.md" {
				noExtPath = strings.Replace(newMdFilepath, "_index.md", "", 1)
			} else {
				noExtPath = strings.Replace(newMdFilepath, ".md", "", 1)
			}
			noExtPath = filepath.ToSlash(noExtPath)
			noExtPath = strings.TrimSuffix(noExtPath, "/")

			v := fmt.Sprintf("{{< ref \"%s\" >}}", "/"+noExtPath)

			file, err := os.OpenFile(path, os.O_RDWR, 0755)
			if err != nil {
				return fmt.Errorf("打开%s遇到错误：%v\n", path, err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()

				if reLine.MatchString(line) {
					// 找出链接
					matches := reLink.FindStringSubmatch(line)
					if len(matches) > 1 {
						k := strings.TrimSuffix(matches[1], "/")
						m[k] = v
						break
					}
				}
			}

			// 检查Scanner是否出错或到达文件末尾
			if err := scanner.Err(); err != nil {
				log.Fatal(fmt.Errorf("Scan%s遇到出错：%v\n", path, err))
			}
		}

		return nil
	})

	if err != nil {
		panic("获取map发生错误")
	}

	return m
}

func replaceLink(m map[string]string, rootDir string) error {
	num := 0
	defer func() {
		fmt.Println("已替换链接的个数：", num)
	}()

	reLine, err := regexp.Compile("^> 原文:")
	if err != nil {
		log.Fatal(fmt.Errorf("创建正则匹配出现错误1：%v\n", err))
	}

	reLink, err := regexp.Compile(`\]\([^\)]+\)`)
	if err != nil {
		log.Fatal(fmt.Errorf("创建正则匹配出现错误2：%v\n", err))
	}

	err = filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		fileName := d.Name()
		if strings.ToLower(filepath.Ext(fileName)) == ".md" {
			file, err := os.OpenFile(path, os.O_RDWR, 0755)
			if err != nil {
				return fmt.Errorf("打开%s遇到出错：%v\n", path, err)
			}
			defer file.Close()

			var lines []string
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()

				if reLine.MatchString(line) {
					lines = append(lines, line)
					continue
				}

				if reLink.MatchString(line) {
					matches := reLink.FindAllString(line, -1)

					for _, match := range matches {
						matchStr1 := strings.TrimRight(match, "/)")
						matchStr1 = strings.TrimPrefix(matchStr1, "](")
						to, ok := m[matchStr1]
						if ok {
							toStr := fmt.Sprintf("](%s)", to)
							line = strings.ReplaceAll(line, match, toStr)
							fmt.Printf("文件%s 中的 %s被替换成%s\n", path, matchStr1, toStr)
							num++
						}

						//if !strings.HasSuffix(matchStr1, "/") {
						//	matchStr2 := matchStr1 + "/"
						//	to, ok = m[matchStr2]
						//	if ok {
						//		toStr := fmt.Sprintf("](%s)", to)
						//		line = strings.ReplaceAll(line, match, toStr)
						//		fmt.Printf("文件%s 中的 %s被替换成%s\n", path, matchStr1, toStr)
						//		num++
						//	}
						//}
					}

					lines = append(lines, line)
					continue
				}

				lines = append(lines, line)
			}

			// 检查Scanner是否出错或到达文件末尾
			if err := scanner.Err(); err != nil {
				log.Fatal(fmt.Errorf("Scan%s遇到出错：%v\n", path, err))
			}

			// 将修改后的内容写入文件
			// 清空文件内容
			if err := file.Truncate(0); err != nil {
				log.Fatal(fmt.Errorf("Truncate%s遇到出错：%v\n", path, err))
			}

			// 将文件指针移动到文件开头
			if _, err := file.Seek(0, 0); err != nil {
				log.Fatal(fmt.Errorf("Seek%s遇到出错：%v\n", path, err))
			}

			// 创建一个用于写入文件的Writer
			writer := bufio.NewWriter(file)

			// 将修改后的内容写入文件
			for _, line := range lines {
				fmt.Fprintln(writer, line)
			}

			// 刷新缓冲区，确保所有内容都写入文件
			if err := writer.Flush(); err != nil {
				log.Fatal(fmt.Errorf("Flush%s遇到出错：%v\n", path, err))
			}
			fmt.Println("run here")
		}

		return nil
	})

	if err != nil {
		panic("替换链接发生错误")
	}

	return nil
}

// 替换带有锚的链接
func replaceLinkWithAnchor(m map[string]string, rootDir string) error {
	num := 0
	defer func() {
		fmt.Println("已替换带有锚的链接的个数：", num)
	}()

	reLine, err := regexp.Compile("^> 原文:")
	if err != nil {
		log.Fatal(fmt.Errorf("创建正则匹配出现错误1：%v\n", err))
	}

	reLink, err := regexp.Compile(`\]\([^\)#]+#[^\)]+\)`)
	if err != nil {
		log.Fatal(fmt.Errorf("创建正则匹配出现错误2：%v\n", err))
	}

	err = filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		fileName := d.Name()
		if strings.ToLower(filepath.Ext(fileName)) == ".md" {
			file, err := os.OpenFile(path, os.O_RDWR, 0755)
			if err != nil {
				return fmt.Errorf("打开%s遇到出错：%v\n", path, err)
			}
			defer file.Close()

			var lines []string
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()

				if reLine.MatchString(line) {
					lines = append(lines, line)
					continue
				}

				if reLink.MatchString(line) {
					matches := reLink.FindAllString(line, -1)

					for _, match := range matches {
						matchStr1 := strings.TrimSuffix(match, ")")
						matchStr1 = strings.TrimPrefix(matchStr1, "](")
						anchorIndex := strings.Index(matchStr1, "#")
						matchStr2 := matchStr1[:anchorIndex]
						anchorStr := matchStr1[anchorIndex:]

						to, ok := m[matchStr2]
						if ok {
							toStr := strings.TrimSuffix(to, "\" >}}") + anchorStr + "\" >}}"
							toStr = fmt.Sprintf("](%s)", toStr)
							line = strings.ReplaceAll(line, match, toStr)
							fmt.Printf("文件%s 中的 %s被替换成%s\n", path, matchStr1, toStr)
							num++
						}
					}

					lines = append(lines, line)
					continue
				}

				lines = append(lines, line)
			}

			// 检查Scanner是否出错或到达文件末尾
			if err := scanner.Err(); err != nil {
				log.Fatal(fmt.Errorf("Scan%s遇到出错：%v\n", path, err))
			}

			// 将修改后的内容写入文件
			// 清空文件内容
			if err := file.Truncate(0); err != nil {
				log.Fatal(fmt.Errorf("Truncate%s遇到出错：%v\n", path, err))
			}

			// 将文件指针移动到文件开头
			if _, err := file.Seek(0, 0); err != nil {
				log.Fatal(fmt.Errorf("Seek%s遇到出错：%v\n", path, err))
			}

			// 创建一个用于写入文件的Writer
			writer := bufio.NewWriter(file)

			// 将修改后的内容写入文件
			for _, line := range lines {
				fmt.Fprintln(writer, line)
			}

			// 刷新缓冲区，确保所有内容都写入文件
			if err := writer.Flush(); err != nil {
				log.Fatal(fmt.Errorf("Flush%s遇到出错：%v\n", path, err))
			}
			fmt.Println("run here")
		}

		return nil
	})

	if err != nil {
		panic("替换带有锚的链接发生错误")
	}

	return nil
}

//
//func findLinkWithAnchor(m map[string]string, rootDir string) error {
//	num := 0
//	defer func() {
//		fmt.Println("找到带有锚的链接的个数：", num)
//	}()
//
//	reLine, err := regexp.Compile("^> 原文:")
//	if err != nil {
//		log.Fatal(fmt.Errorf("创建正则匹配出现错误1：%v\n", err))
//	}
//
//	reLink, err := regexp.Compile(`\]\([^\)#]+#[^\)]+\)`)
//	if err != nil {
//		log.Fatal(fmt.Errorf("创建正则匹配出现错误2：%v\n", err))
//	}
//
//	err = filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
//		if err != nil {
//			return err
//		}
//
//		if d.IsDir() {
//			return nil
//		}
//
//		fileName := d.Name()
//		if fileName != "_index.md" && strings.ToLower(filepath.Ext(fileName)) == ".md" {
//			file, err := os.OpenFile(path, os.O_RDWR, 0755)
//			if err != nil {
//				return fmt.Errorf("打开%s遇到出错：%v\n", path, err)
//			}
//			defer file.Close()
//
//			//var lines []string
//			scanner := bufio.NewScanner(file)
//			for scanner.Scan() {
//				line := scanner.Text()
//
//				if reLine.MatchString(line) {
//					continue
//				}
//
//				if reLink.MatchString(line) {
//					matches := reLink.FindAllString(line, -1)
//
//					for _, match := range matches {
//						matchStr1 := strings.TrimSuffix(match, ")")
//						matchStr1 = strings.TrimPrefix(matchStr1, "](")
//						anchorIndex := strings.Index(matchStr1, "#")
//						matchStr2 := matchStr1[:anchorIndex]
//						to, ok := m[matchStr2]
//						if ok {
//							toStr := fmt.Sprintf("](%s)", to)
//							line = strings.ReplaceAll(line, match, toStr)
//							fmt.Printf("文件%s 中的 %s\n", path, matchStr1)
//							num++
//						}
//					}
//
//					continue
//				}
//			}
//
//			// 检查Scanner是否出错或到达文件末尾
//			if err := scanner.Err(); err != nil {
//				log.Fatal(fmt.Errorf("Scan%s遇到出错：%v\n", path, err))
//			}
//		}
//
//		return nil
//	})
//
//	if err != nil {
//		panic("找带有锚的链接发生错误")
//	}
//
//	return nil
//}
