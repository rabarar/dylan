package pal

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"math"
	"os"

	"github.com/anthonynsimon/bild/transform"
)

type RGBA struct {
	R uint32 `json:"r"`
	G uint32 `json:"g"`
	B uint32 `json:"b"`
	A uint32 `json:"a"`
}

func (rgba *RGBA) RGBA() (uint32, uint32, uint32, uint32) {

	return rgba.R, rgba.G, rgba.B, rgba.A
}

type Brush struct {
	img  *image.Image
	Rank RGBA   `json:"rank"`
	File string `json:"filename"`
}

// Get Image
func (br *Brush) Image() *image.Image {
	return br.img
}

func (br *Brush) fill(basePath string, imgSize int) error {

	fpath := fmt.Sprintf("%s/%s", basePath, br.File)
	imgFile, err := os.Open(fpath)
	if err != nil {
		return err
	}
	defer imgFile.Close()

	img, err := jpeg.Decode(imgFile)
	if err != nil {
		return err
	}

	img = transform.Resize(img, imgSize, imgSize, transform.Linear)
	br.img = &img

	return nil
}

type Palette struct {
	Dirname string  `json:"bucket"`
	List    []Brush `json:"brush"`
}

func NewPalette(bucket string, size int) (*Palette, error) {
	p := &Palette{}
	p.List = make([]Brush, 0)

	p.Dirname = bucket
	list, err := ioutil.ReadDir(bucket)
	if err != nil {
		fmt.Printf("failed to read directory...")
		return nil, err
	}

	// open each file and create a rgb hash
	for _, fileName := range list {

		name := fileName.Name()
		img, c, err := hashPalette(bucket, fileName.Name(), size)
		if err != nil {
			fmt.Printf("failed to hashPalette for %s ...\n", fileName.Name())
			return nil, err
		}

		r, g, b, a := c.RGBA()
		rgba := RGBA{
			R: r,
			G: g,
			B: b,
			A: a,
		}

		brush := Brush{
			img:  img,
			Rank: rgba,
			File: name,
		}

		p.List = append(p.List, brush)

		fmt.Printf("adding %s: %s\n", name, brush.File)

	}

	return p, nil

}

func (p *Palette) Save(dst string) error {
	b, err := json.MarshalIndent(p, "", "\t")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(dst, b, 0644)
	if err != nil {
		return err
	}

	return nil
}

// FillPalette
func (p *Palette) FillPalette(size int) error {

	for i := 0; i < len(p.List); i++ {
		err := p.List[i].fill(p.Dirname, size)
		if err != nil {
			return err
		}
	}

	return nil

}

func norm(p1 RGBA, p2 color.Color) float64 {
	r1, g1, b1, a1 := p1.RGBA()
	r2, g2, b2, a2 := p2.RGBA()

	return math.Sqrt(
		(float64(r1)-float64(r2))*(float64(r1)-float64(r2)) +
			(float64(g1)-float64(g2))*(float64(g1)-float64(g2)) +
			(float64(b1)-float64(b2))*(float64(b1)-float64(b2)) +
			(float64(a1)-float64(a2))*(float64(a1)-float64(a2)))
}

func (p *Palette) Closest(mean color.Color) *image.Image {

	var dist float64 = math.MaxFloat64

	var image *image.Image = nil

	for _, p := range p.List {
		n1 := norm(p.Rank, mean)
		if n1 < dist {
			dist = n1
			image = p.Image()
		}

	}

	return image
}

func LoadPalette(src string) (*Palette, error) {
	p := &Palette{}

	pdata, err := ioutil.ReadFile(src)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(pdata, p)
	if err != nil {
		return nil, err
	}

	return p, nil

}

func hashPalette(bucket, fileName string, imgSize int) (*image.Image, color.Color, error) {

	fpath := fmt.Sprintf("%s/%s", bucket, fileName)
	imgFile, err := os.Open(fpath)
	if err != nil {
		return nil, color.RGBA{0, 0, 0, 0}, err
	}
	defer imgFile.Close()

	img, err := jpeg.Decode(imgFile)
	if err != nil {
		return nil, color.RGBA{0, 0, 0, 0}, err
	}

	img = transform.Resize(img, imgSize, imgSize, transform.Linear)

	var mr, mg, mb, ma uint32 = 0, 0, 0, 0
	var count uint32 = 0
	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			var r, g, b, a uint32 = img.At(x, y).RGBA()
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

	return &img, color.RGBA{uint8(mr / 0x101), uint8(mg / 0x101), uint8(mb / 0x101), uint8(ma / 0x101)}, nil
}
