package main

// RGB 将r,g,b转换成16位表示(r,g,b各占5位，高位忽略)
func RGB(r byte, g byte, b byte) int16 {
	return ((int16)(b&^7) << 7) | (int16)(((g)&^7)<<2) | (int16)((r)>>3)
}

// RED 从16位的颜色表示中取出8位的红色
func RED(color int16) byte {
	return (byte)(((color) & 31) << 3)
}

// BLUE 从16位的颜色中取出8位的蓝色
func BLUE(color int16) byte {
	return (byte)((((color) >> 10) & 255) << 3)
}

// GREEN 从16位的颜色中取出8位的绿色
func GREEN(color int16) byte {
	return (byte)((((color) >> 5) & 255) << 3)
}
