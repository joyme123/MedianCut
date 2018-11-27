package main

type ColorCube struct {
	Lower int16 // 盒子的低位指针
	Upper int16 // 盒子的高位指针
	Count int   // 盒子中的颜色数量
	Level int   // 盒子切割的深度

	Rmin, Rmax, Gmin, Gmax, Bmin, Bmax byte // RGB值的范围
}

// Shrink 方法会将盒子收缩
func (cube *ColorCube) Shrink() {

}

// GetColor 用来获取盒子代表的颜色
func (cube *ColorCube) GetColor() {

}

type HistList struct {
	histPtr []int16
}

func (hist *HistList) Init(histPtr []int16, longdim int) {
	hist.histPtr = histPtr
}

func (hist HistList) Len() int {
	return len(hist.histPtr)
}

func (hist HistList) Less(i, j int) bool {
	colorA := hist.histPtr[i]
	colorB := hist.histPtr[j]

	var c1 byte
	var c2 byte

	switch longdim {

	case 0:
		c1 = RED(colorA)
		c2 = RED(colorB)
		break
	case 1:
		c1 = GREEN(colorA)
		c2 = GREEN(colorB)
		break
	case 2:
		c1 = BLUE(colorA)
		c2 = BLUE(colorB)
		break
	}
	return c1-c2 < 0

}

func (hist HistList) Swap(i, j int) {
	hist.histPtr[i], hist.histPtr[j] = hist.histPtr[j], hist.histPtr[i]
}
