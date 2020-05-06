package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rabarar/dylan/mosaic"
	"github.com/rabarar/dylan/pal"
)

var (
	colorMode = map[string]mosaic.ColorMode{
		"tile":   mosaic.ColorModeMeanTile,
		"random": mosaic.ColorModeRandom,
		"mean":   mosaic.ColorModeMean,
		"copy":   mosaic.ColorModeCopy,
	}
)

func main() {

	srcFilename := flag.String("src", "", "input jpeg filename ")
	dstFilename := flag.String("dst", "output.jpg", "output jpeg filename ")
	palFilename := flag.String("palette", "palette.json", "json for palette")
	mode := flag.String("mode", "tile", "tiling mode: tile | random | mean | copy")

	flag.Parse()

	if *srcFilename == "" {
		fmt.Printf("must specify a source filename, exiting...\n")
		os.Exit(1)
	}

	mo, err := mosaic.NewMosaic(*srcFilename, mosaic.WINDOW_SIZE)
	if err != nil {
		panic(err)
	}

	p, err := pal.LoadPalette(*palFilename)
	if err != nil {
		panic(err)
	}
	fmt.Printf("loaded Palette...\n")

	err = p.FillPalette()
	if err != nil {
		panic(err)
	}
	fmt.Printf("filled Palette...\n")

	var cmode mosaic.ColorMode = colorMode[*mode]

	err = mo.Color(p, cmode)
	if err != nil {
		panic(err)
	}

	err = mo.Save(*dstFilename, mosaic.JPEG_QUALITY)
	if err != nil {
		panic(err)
	}

}
