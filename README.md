# dylan
Mosaic generation tool

generate mosaic pics by "painting" with other photos...

`$ ./dylan -h
Usage of ./dylan:
  -dst string
    	output jpeg filename  (default "output.jpg")
  -palette string
    	json for palette (default "palette.json")
  -src string
    	input jpeg filename
` 

use ./palette to generate a palette for ./dylan

`
$ ./palette/./palette -h
Usage of ./palette/./palette:
  -action string
    	palette actions (default "create | verify, resize")
  -dst string
    	resize palette directory (default "thumbs")
  -palette string
    	palette filename (default "palette.json")
  -size int
    	resize dimension (default 25)
  -src string
    	source palette directory (default "bucket")
`
