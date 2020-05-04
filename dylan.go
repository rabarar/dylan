package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"math/rand"
	"os"

	"github.com/rabarar/dylan/pal"
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

func main() {

	data, err := ioutil.ReadFile("dylan.jpg")
	if err != nil {
		panic(err)
	}

	img, err := jpeg.Decode(bytes.NewReader(data))

	if err != nil {
		panic(err)
	}

	minx := img.Bounds().Min.X
	miny := img.Bounds().Min.Y

	maxx := img.Bounds().Max.X
	maxy := img.Bounds().Max.Y

	windows := []*WindowBox{}

	// set the window size
	size := WINDOW_SIZE

	for x := minx; x < maxx; x += size {

		for y := miny; y < maxy; y += size {

			// create a new windowbox
			w := NewWindowBox(size)

			// get the color at the x,y coordinate and copy into the window
			w.Min = image.Point{x, y}
			w.Max = image.Point{x + size, y + size}

			for wx := 0; wx < size; wx++ {
				for wy := 0; wy < size; wy++ {
					w.Window[wx][wy] = img.At(wx+x, wy+y)
				}
			}

			windows = append(windows, w)

		}

	}

	// dump the window array
	if false {
		fmt.Printf("len of windows = %d\n", len(windows))
		for i := 0; i < len(windows); i++ {

			fmt.Printf("window[%d] min = %d, max = %d\n", i+1, windows[i].Min, windows[i].Max)
		}
	}

	// create new image from windows
	newImg := image.NewRGBA(image.Rect(minx, miny, maxx, maxy))

	p, err := pal.LoadPalette("palette.json")
	if err != nil {
		panic(err)
	}
	fmt.Printf("loaded Palette...\n")

	err = p.FillPalette()
	if err != nil {
		panic(err)
	}
	fmt.Printf("filled Palette...\n")

	// fill it with a (random) tile
	for x := minx; x < maxx; x += size {

		for y := miny; y < maxy; y += size {

			// select a palette
			ip := *(p.List[rand.Intn(len(p.List))].Image())

			for wx := 0; wx < size; wx++ {
				for wy := 0; wy < size; wy++ {
					newImg.Set(x+wx, y+wy, ip.At(wx, wy))
				}
			}
		}
	}

	// fill it with the mean
	if false {
		for i := 0; i < len(windows); i++ {
			windows[i].CalcMean()
			for x := windows[i].Min.X; x < windows[i].Max.X; x++ {
				for y := windows[i].Min.Y; y < windows[i].Max.Y; y++ {
					newImg.Set(x, y, windows[i].Mean())
				}
			}
		}
	}

	// fill it with the closest tile
	if true {
		for i := 0; i < len(windows); i++ {
			//ip := *(p.List[0].Image()) // p.Closest(windows[i].Mean())
			windows[i].CalcMean()
			var ip *image.Image = p.Closest(windows[i].Mean())
			if ip == nil {
				panic("can't be nil!")

			}
			for x := windows[i].Min.X; x < windows[i].Max.X; x++ {
				for y := windows[i].Min.Y; y < windows[i].Max.Y; y++ {
					newImg.Set(x, y, (*ip).At(x-windows[i].Min.X, y-windows[i].Min.Y))
				}
			}
		}
	}

	// basic copy
	if false {
		for x := minx; x < maxx; x++ {
			for y := miny; y < maxy; y++ {
				newImg.Set(x, y, img.At(x, y))
			}
		}
	}

	// write it to disk
	file, err := os.Create("debug.jpg")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = jpeg.Encode(file, newImg, &jpeg.Options{Quality: JPEG_QUALITY})
	if err != nil {
		panic(err)
	}

}
