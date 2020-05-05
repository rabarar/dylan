package main

import (
	"fmt"

	"github.com/rabarar/dylan/mosaic"
	"github.com/rabarar/dylan/pal"
)

func main() {

	mo, err := mosaic.NewMosaic("dylan.jpg", mosaic.WINDOW_SIZE)
	if err != nil {
		panic(err)
	}

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

	err = mo.Color(p)
	if err != nil {
		panic(err)
	}

	err = mo.Save("debug.jpg", mosaic.JPEG_QUALITY)
	if err != nil {
		panic(err)
	}

}
