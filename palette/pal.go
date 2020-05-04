package main

import (
	"fmt"

	"github.com/rabarar/dylan/pal"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
)

func main() {

	resizeDir := "thumbs"

	if true {
		p, err := pal.NewPalette("bucket")
		if err != nil {
			panic(err)
		}
		p.Save("myBucket.json")
	}

	p, err := pal.LoadPalette("myBucket.json")
	if err != nil {
		panic(err)
	}

	if false {
		// resize and save thumbs
		fmt.Printf("successfully loaded palette: %d\n", len(p.List))
		for i := 0; i < len(p.List); i++ {
			fmt.Printf("color: %d\n", p.List[i].Rank)
			img, err := imgio.Open(fmt.Sprintf("%s/%s", "bucket", p.List[i].File))
			if err != nil {
				panic(err)
			}
			resized := transform.Resize(img, 25, 25, transform.Linear)

			if err := imgio.Save(fmt.Sprintf("%s/%s", resizeDir, p.List[i].File), resized, imgio.JPEGEncoder(100)); err != nil {
				panic(err)
			}

		}
	}

}
