package mosaic

import (
	"image"
	"image/color"
)

const (
	WINDOW_SIZE  = 25
	JPEG_QUALITY = 100
)

type WindowBox struct {
	Min    image.Point
	Max    image.Point
	Window [][]color.Color
	mean   color.Color
	size   int
}

func NewWindowBox(size int) *WindowBox {
	wb := WindowBox{}
	wb.Window = make([][]color.Color, size)
	for i := 0; i < size; i++ {
		wb.Window[i] = make([]color.Color, size)
	}
	wb.size = size
	return &wb
}

func (wb *WindowBox) Mean() color.Color {
	return wb.mean
}

func (wb *WindowBox) CalcMean2() {
	wb.mean = wb.Window[0][0]

}

func (wb *WindowBox) CalcMean() {

	var r, g, b, a uint32
	var mr, mg, mb, ma uint32

	var count uint32 = 0
	for x := 0; x < wb.size; x++ {
		for y := 0; y < wb.size; y++ {
			r, g, b, a = wb.Window[x][y].RGBA()
			mr += r
			mg += g
			mb += b
			ma += a
			count++
		}
	}

	mr /= count
	mg /= count
	mb /= count
	ma /= count

	// NOTE: need to divide by 0x101 for the shiftup

	wb.mean = color.RGBA{uint8(mr / 0x101), uint8(mg / 0x101), uint8(mb / 0x101), uint8(ma / 0x101)}
}
