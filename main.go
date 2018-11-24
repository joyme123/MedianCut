package main

import (
	"fmt"
)

const HSIZE = 32768

var histPtr [HSIZE]int16

// MedianCut 用来切割ColorCude, hist是图片中颜色的直方图
func MedianCut(hist []int16, colorMap [][3]byte, maxCubes int) int {
	var lr, lg, lb byte
	var i, median, color int16
	var count int
	var k, level, ncubes, splitpos int

	var cube, cubeA, cubeB ColorCube

	ncubes = 0
	cube.count = 0

	i = 0
	color = 0

	for ; i < HSIZE-1; i++ {
		if hist[i] != 0 {
			histPtr[color] = i
			color += 1
			cube.count = cube.count + (int)(hist[i])
		}
	}

}

func main() {
	fmt.Println("hello world")
}
