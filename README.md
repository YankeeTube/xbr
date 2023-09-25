# XBR Pixel Scale on Golang
Unfortunately, `xbr` algorithm is not being shared in Golang.  
So I saw this [repository](https://github.com/joseprio/xBRjs/tree/master) written in JS and implemented it in Go language.  
I couldn't solve the `3x`. but `2x` and `4x` operate normally.

## Image Compare
1. original
2. 2x
3. 3x
4. 4x  
![original](./fixture/stand@1x.png)
![scale2x](./fixture/2x.png)
![scale3x](./fixture/3x.png)
![scale4x](./fixture/4x.png)

### Usage
```go
package main

import (
	"bytes"
	"image"
	"image/png"
	"log"
	"os"
	"xbr/scale" // change import
)

func saveImageAsPNG(img image.Image, filename string) error {
	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()
	return png.Encode(outFile, img)
}

func main() {
	file, err := os.ReadFile("fixture/stand@1x.png")
	if err != nil {
		log.Fatalln(err)
	}

	img, _, err := image.Decode(bytes.NewReader(file))
	if err != nil {
		log.Fatalln(err)
	}

	res := scale.Xbr(img, 2)
	res3 := scale.Xbr(img, 3)
	res4 := scale.Xbr(img, 4)
	_ = saveImageAsPNG(res, "output2.png")
	_ = saveImageAsPNG(res3, "output3.png")
	_ = saveImageAsPNG(res4, "output4.png")
}
```