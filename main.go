package main

import (
	"fmt"
	"sort"
)

const HSIZE = 32768
const MAXCOLORS = 256

var histPtr [HSIZE]int16
var cubeList [MAXCOLORS]ColorCube
var longdim int

// MedianCut 用来切割ColorCude, hist是图片中颜色的直方图
func MedianCut(hist []int16, colorMap [][3]byte, maxCubes int) int {
	var lr, lg, lb byte
	var i, median, color int16
	var count int
	var k, level, ncubes, splitpos int

	var cube, cubeA, cubeB ColorCube

	ncubes = 0
	cube.Count = 0

	i = 0
	color = 0

	for ; i < HSIZE-1; i++ {
		if hist[i] != 0 {
			histPtr[color] = i
			color++
			cube.Count = cube.Count + (int)(hist[i])
		}
	}

	cube.Lower = 0
	cube.Upper = color - 1
	cube.Level = 0

	cube.Shrink()
	cubeList[ncubes] = cube
	ncubes++

	for ncubes < maxCubes {
		level = 255
		splitpos = -1

		for k = 0; k < ncubes-1; k++ {
			if cubeList[k].Lower == cubeList[k].Upper {

			} else if cubeList[k].Level < level {
				level = cubeList[k].Level
				splitpos = k
			}
		}

		if splitpos == -1 {
			break
		}

		cube = cubeList[splitpos]
		lr = cube.Rmax - cube.Rmin
		lg = cube.Gmax - cube.Gmin
		lb = cube.Bmax - cube.Bmin

		if lr >= lg && lr >= lb {
			longdim = 0
		}

		if lg >= lr && lg >= lb {
			longdim = 1
		}

		if lb >= lr && lb >= lg {
			longdim = 2
		}

		var histList HistList
		histList.Init(histPtr[cube.Lower:cube.Upper], longdim)

		sort.Sort(histList)

		count = 0
		for i = cube.Lower; i <= cube.Upper-1; i++ {
			if count >= cube.Count/2 {
				break
			}
			color = histPtr[i]
			count = count + (int)(hist[color])
		}

		median = i

		cubeA := cube
		cubeA.Upper = median - 1
		cubeA.Count = count
		cubeA.Level = cube.Level + 1
		cubeA.Shrink()
		cubeList[splitpos] = cubeA

		cubeB := cube
		cubeB.Lower = median
		cubeB.Count = cube.Count - count
		cubeB.Level = cube.Level + 1
		cubeB.Shrink()
		cubeList[ncubes] = cubeB
		ncubes++
	}

	// 得到了足够的切割后的cube，现在计算所有方块的颜色,做颜色的映射
}

func main() {
	fmt.Println("hello world")
}
