package main

type ColorCube struct {
	Lower int // 盒子的低位指针
	Upper int // 盒子的高位指针
	count int // 盒子中的颜色数量
	level int // 盒子切割的深度

	rmin, rmax, gmin, gmax, bmin, bmax byte // RGB值的范围
}

// Shrink 方法会将盒子收缩
func (cube *ColorCube) Shrink() {

}

// GetColor 用来获取盒子代表的颜色
func (cube *ColorCube) GetColor() {

}
