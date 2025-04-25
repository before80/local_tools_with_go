package main

import (
	"github.com/go-vgo/robotgo"
)

func main() {
	// 移动鼠标到屏幕中心
	width, height := robotgo.GetScreenSize()
	robotgo.MoveMouse(width/2, height/2)

	// 鼠标点击
	robotgo.MouseClick("left", true)   // 左键单击
	robotgo.MouseClick("right", false) // 右键不双击

	// 鼠标拖拽
	robotgo.MoveMouseSmooth(100, 100)
	robotgo.MouseToggle("down", "left")
	robotgo.MoveMouseSmooth(200, 200)
	robotgo.MouseToggle("up", "left")

	// 获取鼠标位置
	x, y := robotgo.GetMousePos()
	println("Mouse position:", x, y)

	// 鼠标滚轮
	robotgo.ScrollMouse(10, "up") // 向上滚动10个单位
}
