package main

import (
	"bytes"
	"image"
	"image/png"
	"log"
	"os"
	"xbr/scale"
)

func saveImageAsPNG(img image.Image, filename string) error {
	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()

	return png.Encode(outFile, img)
}

func resizeNearestNeighbor(img image.Image, ratio int) *image.RGBA {
	dx := img.Bounds().Dx()
	dy := img.Bounds().Dy()
	w := img.Bounds().Max.X * ratio
	h := img.Bounds().Max.Y * ratio
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	xRatio := float64(dx) / float64(w)
	yRatio := float64(dy) / float64(h)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			srcX := int(float64(x) * xRatio)
			srcY := int(float64(y) * yRatio)
			dst.Set(x, y, img.At(srcX, srcY))
		}
	}
	return dst
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

	res2 := scale.Xbr(img, 2, true, true)
	res3 := scale.Xbr(img, 3, true, true)
	res4 := scale.Xbr(img, 4, true, true)
	_ = saveImageAsPNG(res2, "output2.png")
	_ = saveImageAsPNG(res3, "output3.png")
	_ = saveImageAsPNG(res4, "output4.png")

	_ = saveImageAsPNG(resizeNearestNeighbor(img, 4), "near.png")
}
