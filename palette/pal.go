package main

import (
	"flag"
	"fmt"

	"github.com/rabarar/dylan/pal"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
)

const (
	PalActionCreate PalAction = iota
	PalActionVerify
	PalActionResize
)

var (
	actionCase = map[string]PalAction{
		"create": PalActionCreate,
		"verify": PalActionVerify,
		"resize": PalActionResize,
	}
)

type PalAction int

func main() {

	var srcDir string
	var resizeDir string
	var action string
	var palJson string
	var imgSize int

	flag.StringVar(&srcDir, "src", "bucket", "source palette directory")
	flag.StringVar(&resizeDir, "dst", "thumbs", "resize palette directory")
	flag.StringVar(&palJson, "palette", "palette.json", "palette filename")
	flag.StringVar(&action, "action", "create | verify, resize", "palette actions")
	flag.IntVar(&imgSize, "size", 25, "resize dimension")

	flag.Parse()

	switch actionCase[action] {
	case PalActionCreate:
		p, err := pal.NewPalette(srcDir)
		if err != nil {
			panic(err)
		}
		p.Save(fmt.Sprintf("%s.json", srcDir))

	case PalActionVerify:
		_, err := pal.LoadPalette(palJson)
		if err != nil {
			panic(err)
		}
		fmt.Printf("palette ok!\n")

	case PalActionResize:
		p, err := pal.LoadPalette(palJson)
		if err != nil {
			panic(err)
		}

		for i := 0; i < len(p.List); i++ {
			fmt.Printf("color: %d\n", p.List[i].Rank)
			img, err := imgio.Open(fmt.Sprintf("%s/%s", p.Dirname, p.List[i].File))

			if err != nil {
				panic(err)
			}
			resized := transform.Resize(img, imgSize, imgSize, transform.Linear)

			if err := imgio.Save(fmt.Sprintf("%s/%s", resizeDir, p.List[i].File), resized, imgio.JPEGEncoder(100)); err != nil {
				panic(err)
			}

		}

		np, err := pal.NewPalette(resizeDir)
		if err != nil {
			panic(err)
		}
		np.Save(fmt.Sprintf("%s.json", resizeDir))
	}

}
