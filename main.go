package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"sort"

	"github.com/joyme123/MedianCut/util"

	_ "image/gif"
	_ "image/jpeg"
	"image/png"
)

const HSIZE = 32768
const MAXCOLORS = 256

var histPtr [HSIZE]uint16
var cubeList [MAXCOLORS]util.ColorCube
var longdim int

// MedianCut 用来切割ColorCude, hist是图片中颜色的直方图
func MedianCut(hist []uint16, colorMap [][3]byte, maxCubes int) int {
	var lr, lg, lb byte
	var i, median, color uint16
	var count int
	var k, level, ncubes, splitpos int

	var cube util.ColorCube

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

	cube.HistPtr = histPtr[:]
	cube.Shrink()
	cubeList[ncubes] = cube
	ncubes++

	for ncubes < maxCubes {
		level = 255
		splitpos = -1

		for k = 0; k <= ncubes-1; k++ {
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

		sort.Sort(util.HistList{histPtr[cube.Lower:cube.Upper], longdim})

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
	invMap(hist, colorMap, ncubes)

	return ncubes
}

func invMap(hist []uint16, colorMap [][3]byte, ncubes int) {
	var r, g, b byte
	var i, k, color uint16
	var rsum, gsum, bsum float32
	var cube util.ColorCube

	for k = 0; k <= (uint16)(ncubes)-1; k++ {
		cube = cubeList[k]
		rsum = 0.0
		gsum = 0.0
		bsum = 0.0

		// fmt.Printf("upper是%d, lower是%d\n", cube.Upper, cube.Lower)
		for i = cube.Lower; i <= cube.Upper; i++ {
			// fmt.Printf("i是%d\n", i)
			color = histPtr[i]
			r = util.RED(color)
			rsum += (float32)(r) * (float32)(hist[color])
			g = util.GREEN(color)
			gsum += (float32)(g) * (float32)(hist[color])
			b = util.BLUE(color)
			bsum += (float32)(b) * (float32)(hist[color])
		}

		colorMap[k][0] = (byte)(rsum / (float32)(cube.Count))
		colorMap[k][1] = (byte)(gsum / (float32)(cube.Count))
		colorMap[k][2] = (byte)(bsum / (float32)(cube.Count))
	}

	for k = 0; k < (uint16)(ncubes); k++ {
		cube = cubeList[k]
		for i = cube.Lower; i < cube.Upper; i++ {
			color = histPtr[i]
			hist[color] = k
		}
	}
}

func main() {

	var maxCubes int
	var img string
	var out string
	var debug bool

	flag.IntVar(&maxCubes, "num", 256, "output color number")
	flag.StringVar(&img, "img", "", "image path")
	flag.StringVar(&out, "out", "", "out put file")
	flag.BoolVar(&debug, "debug", false, "enable debug mode")

	flag.Parse()

	fmt.Println(img)

	var hist []uint16      // 图片的颜色统计直方图
	var colorMap [][3]byte // 切割后的立方体对应的颜色

	colorMap = make([][3]byte, maxCubes)

	reader, err := os.Open(img)
	if err != nil {
		log.Fatal(err)
	}

	defer reader.Close()
	m, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	bounds := m.Bounds()

	imgWidth := bounds.Size().X
	imgHeight := bounds.Size().Y

	colorCnt := imgWidth * imgHeight

	hist = make([]uint16, colorCnt)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := m.At(x, y).RGBA()

			r = r >> 8
			g = g >> 8
			b = b >> 8

			color := util.RGB((byte)(r), (byte)(g), (byte)(b))

			hist[color]++
		}
	}

	ncubes := MedianCut(hist, colorMap, maxCubes)

	// 打印出colorMap
	// fmt.Printf("colorMap: %v\n", colorMap)

	fmt.Printf("实际切割后生成的立方体数量:%d\n", ncubes)

	if out != "" {

		var palette []color.Color

		palette = make([]color.Color, ncubes)

		pIndex := 0

		for colorMapIndex, rgb := range colorMap {

			if colorMapIndex >= ncubes {
				break
			}

			rc := rgb[0]
			gc := rgb[1]
			bc := rgb[2]

			rgba := color.RGBA{rc, gc, bc, 255}
			palette[pIndex] = rgba
			pIndex++
		}

		// 输出调色板
		if debug {
			debugPaletted := image.NewPaletted(image.Rect(0, 0, 640, 640), palette)
			for y := 0; y < imgHeight; y++ {
				for x := 0; x < imgWidth; x++ {
					colorIndex := (x / 40) + 16*(y/40)

					debugPaletted.SetColorIndex(x, y, (uint8)(colorIndex))
				}
			}

			file, err := os.Create("debug.png")

			if err != nil {
				log.Fatal(err)
			}

			err2 := png.Encode(file, debugPaletted)

			if err2 != nil {
				log.Fatal(err2)
			}
		}

		// 要生成改变颜色后的图片,这里使用调色板模式
		paletted := image.NewPaletted(image.Rect(0, 0, imgWidth, imgHeight), palette)

		for y := 0; y < imgHeight; y++ {
			for x := 0; x < imgWidth; x++ {

				r, g, b, _ := m.At(x, y).RGBA()

				r = r >> 8
				g = g >> 8
				b = b >> 8

				color := util.RGB((byte)(r), (byte)(g), (byte)(b))
				paletted.SetColorIndex(x, y, (uint8)(hist[color]))
			}
		}

		file, err := os.Create(out)

		if err != nil {
			log.Fatal(err)
		}

		err2 := png.Encode(file, paletted)

		if err2 != nil {
			log.Fatal(err2)
		}
	}
}
