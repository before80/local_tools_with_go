package main

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/tailscale/win"
)

func main() {
	// Step 1: 打开 Typora
	fmt.Println("正在启动 Typora...")
	execTypora()

	// 等待 Typora 启动（根据系统性能调整等待时间）
	time.Sleep(3 * time.Second)

	// Step 2: 模拟点击“文件”菜单
	fmt.Println("点击 Typora 的文件菜单...")
	robotgo.KeyTap("alt") // Alt 键激活菜单栏
	time.Sleep(500 * time.Millisecond)
	robotgo.TypeStr("f") // 按下 "F" 键选择“文件”菜单
	time.Sleep(500 * time.Millisecond)

	// Step 3: 模拟点击“打开”菜单项
	fmt.Println("点击打开菜单项...")
	//robotgo.TypeStr("o") // 按下 "O" 键选择“打开”菜单项
	robotgo.KeyTap("o", "ctrl")
	time.Sleep(1 * time.Second)

	// Step 4: 移动打开窗口到指定位置
	fmt.Println("移动打开窗口到 (200, 200)...")
	moveOpenWindowTo(200, 200)

	// Step 5: 输入文件路径并确认
	fmt.Println("输入文件路径并打开...")
	inputFilePathAndConfirm("D:\\Docs\\hugos\\lang\\content\\c\\types\\_index.md")
}

// 启动 Typora
func execTypora() {
	// 使用 os/exec 启动 Typora 的可执行文件
	cmd := exec.Command("D:\\tools\\typora\\Typora.exe")
	err := cmd.Start()
	if err != nil {
		fmt.Println("无法启动 Typora:", err)
		return
	}
	fmt.Println("Typora 已启动")
}

// 移动打开窗口到指定位置
func moveOpenWindowTo(x, y int) {
	// 查找打开窗口（假设标题包含“打开”）
	hwnd := robotgo.FindWindow("打开")
	if hwnd == 0 {
		fmt.Println("未找到打开窗口！")
		return
	}

	// 使用 Windows API 移动窗口
	moveWindow(hwnd, x, y)
}

// 输入文件路径并确认
func inputFilePathAndConfirm(filePath string) {
	// 输入文件路径
	robotgo.TypeStr(filePath)
	time.Sleep(500 * time.Millisecond)

	// 按下 Enter 键确认
	robotgo.KeyTap("enter")
}

// 使用 Windows API 移动窗口
func moveWindow(hwnd win.HWND, x, y int) {
	// 调用 SetWindowPos 函数移动窗口
	win.SetWindowPos(
		hwnd,           // 窗口句柄
		win.HWND_TOP,   // 置顶
		int32(x),       // 新的 X 坐标
		int32(y),       // 新的 Y 坐标
		0,              // 宽度（保持不变）
		0,              // 高度（保持不变）
		win.SWP_NOSIZE, // 不改变窗口大小
	)
}
