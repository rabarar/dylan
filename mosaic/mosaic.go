package mosaic

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"math/rand"

	"github.com/rabarar/dylan/pal"
)

const (
	WINDOW_SIZE  = 25
	JPEG_QUALITY = 100
)

type Mosaic struct {
	boxes   *[]*WindowBox
	srcFile string
	newImg  *image.RGBA
	size    int
	img     image.Image
}

type WindowBox struct {
	Min    image.Point
	Max    image.Point
	Window [][]color.Color
	mean   color.Color
	size   int
}

func newWindowBox(size int) *WindowBox {
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

func NewMosaic(srcFile string, size int) (*Mosaic, error) {

	mo := &Mosaic{}
	mo.srcFile = srcFile

	data, err := ioutil.ReadFile(srcFile)
	if err != nil {
		return nil, err
	}

	mo.img, err = jpeg.Decode(bytes.NewReader(data))

	if err != nil {
		return nil, err
	}

	minx := mo.img.Bounds().Min.X
	miny := mo.img.Bounds().Min.Y

	maxx := mo.img.Bounds().Max.X
	maxy := mo.img.Bounds().Max.Y

	windows := []*WindowBox{}

	// create new image from windows
	mo.newImg = image.NewRGBA(image.Rect(minx, miny, maxx, maxy))

	for x := minx; x < maxx; x += size {

		for y := miny; y < maxy; y += size {

			// create a new windowbox
			w := newWindowBox(size)

			// get the color at the x,y coordinate and copy into the window
			w.Min = image.Point{x, y}
			w.Max = image.Point{x + size, y + size}

			for wx := 0; wx < size; wx++ {
				for wy := 0; wy < size; wy++ {
					w.Window[wx][wy] = mo.img.At(wx+x, wy+y)
				}
			}

			windows = append(windows, w)

		}

	}

	mo.boxes = &windows

	return mo, nil

}

func (mo *Mosaic) Color(p *pal.Palette) error {

	minx := mo.img.Bounds().Min.X
	miny := mo.img.Bounds().Min.Y

	maxx := mo.img.Bounds().Max.X
	maxy := mo.img.Bounds().Max.Y

	// fill it with a (random) tile
	if false {
		for x := minx; x < maxx; x += mo.size {

			for y := miny; y < maxy; y += mo.size {

				// select a palette
				ip := *(p.List[rand.Intn(len(p.List))].Image())

				for wx := 0; wx < mo.size; wx++ {
					for wy := 0; wy < mo.size; wy++ {
						mo.newImg.Set(x+wx, y+wy, ip.At(wx, wy))
					}
				}
			}
		}
	}

	// fill it with the mean
	if false {
		for i := 0; i < len(*mo.boxes); i++ {
			(*mo.boxes)[i].CalcMean()
			for x := (*mo.boxes)[i].Min.X; x < (*mo.boxes)[i].Max.X; x++ {
				for y := (*mo.boxes)[i].Min.Y; y < (*mo.boxes)[i].Max.Y; y++ {
					mo.newImg.Set(x, y, (*mo.boxes)[i].Mean())
				}
			}
		}
	}

	// fill it with the closest tile
	if true {
		for i := 0; i < len((*mo.boxes)); i++ {
			//ip := *(p.List[0].Image()) // p.Closest(windows[i].Mean())
			(*mo.boxes)[i].CalcMean()
			var ip *image.Image = p.Closest((*mo.boxes)[i].Mean())
			if ip == nil {
				panic("can't be nil!")

			}
			for x := (*mo.boxes)[i].Min.X; x < (*mo.boxes)[i].Max.X; x++ {
				for y := (*mo.boxes)[i].Min.Y; y < (*mo.boxes)[i].Max.Y; y++ {
					mo.newImg.Set(x, y, (*ip).At(x-(*mo.boxes)[i].Min.X, y-(*mo.boxes)[i].Min.Y))
				}
			}
		}
	}

	// basic copy
	if false {
		for x := minx; x < maxx; x++ {
			for y := miny; y < maxy; y++ {
				mo.newImg.Set(x, y, mo.img.At(x, y))
			}
		}
	}

	return nil
}
