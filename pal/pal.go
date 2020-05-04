package pal

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"os"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
)

type RGBA struct {
	R uint32 `json:"r"`
	G uint32 `json:"g"`
	B uint32 `json:"b"`
	A uint32 `json:"a"`
}

type Brush struct {
	img  *image.Image
	Rank RGBA   `json:"rank"`
	File string `json:"filename"`
}

type Palette struct {
	Dirname string  `json:"bucket"`
	List    []Brush `json:"brush"`
}

func NewPalette(bucket string) (*Palette, error) {
	p := &Palette{}
	p.List = make([]Brush, 0)

	p.Dirname = bucket
	list, err := ioutil.ReadDir(bucket)
	if err != nil {
		return nil, err
	}

	// open each file and create a rgb hash
	for _, fileName := range list {

		name := fileName.Name()
		img, c, err := hashPalette(bucket, fileName.Name())
		if err != nil {
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

func hashPalette(bucket, fileName string) (*image.Image, color.Color, error) {

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
