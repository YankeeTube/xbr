package scale

import (
	"image"
	"math"
)

func alphaBlend32W(dst, src uint32, blendColors bool) uint32 {
	if blendColors {
		return pixelInterpolate(dst, src, 7, 1)
	}
	return dst
}

func alphaBlend64W(dst, src uint32, blendColors bool) uint32 {
	if blendColors {
		return pixelInterpolate(dst, src, 3, 1)
	}
	return dst
}

func alphaBlend128W(dst, src uint32, blendColors bool) uint32 {
	if blendColors {
		return pixelInterpolate(dst, src, 1, 1)
	}
	return dst
}

func alphaBlend192W(dst, src uint32, blendColors bool) uint32 {
	if blendColors {
		return pixelInterpolate(dst, src, 1, 3)
	}
	return src
}

func alphaBlend224W(dst, src uint32, blendColors bool) uint32 {
	if blendColors {
		return pixelInterpolate(dst, src, 1, 7)
	}
	return src
}

func isEqual(A, B uint32, scaleAlpha bool) bool {
	alphaA := ((A & AlphaMask) >> 24) & 0xff
	alphaB := ((B & AlphaMask) >> 24) & 0xff

	if alphaA == 0 && alphaB == 0 {
		return true
	}

	if !scaleAlpha && (alphaA < 255 || alphaB < 255) {
		return false
	}

	if alphaA == 0 || alphaB == 0 {
		return false
	}

	yA, uA, vA := getYuv(A)
	yB, uB, vB := getYuv(B)

	if math.Abs(yA-yB) > ThreadsHoldY {
		return false
	}
	if math.Abs(uA-uB) > ThreadsHoldU {
		return false
	}
	if math.Abs(vA-vB) > ThreadsHoldV {
		return false
	}

	return true
}

func pixelInterpolate(A, B, q1, q2 uint32) uint32 {
	alphaA := (A & AlphaMask) >> 24 & 0xFF
	alphaB := (B & AlphaMask) >> 24 & 0xFF

	var r, g, b, a float64

	if alphaA == 0 {
		r = float64(B & RedMask)
		g = float64((B & GreenMask) >> 8)
		b = float64((B & BlueMask) >> 16)
	} else if alphaB == 0 {
		r = float64(A & RedMask)
		g = float64((A & GreenMask) >> 8)
		b = float64((A & BlueMask) >> 16)
	} else {
		r = float64(q2*(B&RedMask)+q1*(A&RedMask)) / float64(q1+q2)
		g = float64(q2*((B&GreenMask)>>8)+q1*((A&GreenMask)>>8)) / float64(q1+q2)
		b = float64(q2*((B&BlueMask)>>16)+q1*((A&BlueMask)>>16)) / float64(q1+q2)
	}
	a = float64(q2*alphaB+q1*alphaA) / float64(q1+q2)

	// Flooring operation to match JavaScript's behavior
	return uint32(math.Floor(r)) | (uint32(math.Floor(g)) << 8) | (uint32(math.Floor(b)) << 16) | (uint32(math.Floor(a)) << 24)
}

// Convert an ARGB byte to YUV
func getYuv(p uint32) (float64, float64, float64) {
	r := float64(p & RedMask)
	g := float64((p & GreenMask) >> 8)
	b := float64((p & BlueMask) >> 16)
	y := r*0.299000 + g*0.587000 + b*0.114000
	u := r*-0.168736 + g*-0.331264 + b*0.500000
	v := r*0.500000 + g*-0.418688 + b*-0.081312
	return y, u, v
}

func yuvDifference(A uint32, B uint32, scaleAlpha bool) float64 {
	alphaA := (A & AlphaMask) >> 24 & 0xff
	alphaB := (B & AlphaMask) >> 24 & 0xff

	if alphaA == 0 && alphaB == 0 {
		return 0
	}

	if !scaleAlpha && (alphaA < 255 || alphaB < 255) {
		// Very large value not attainable by the thresholds
		return 1000000
	}

	if alphaA == 0 || alphaB == 0 {
		// Very large value not attainable by the thresholds
		return 1000000
	}

	yA, uA, vA := getYuv(A)
	yB, uB, vB := getYuv(B)

	// Add HQx filters threshold & return
	return math.Abs(yA-yB)*ThreadsHoldY +
		math.Abs(uA-uB)*ThreadsHoldU +
		math.Abs(vA-vB)*ThreadsHoldV
}

func getRelatedPoints(oriPixelView []uint32, oriX int, oriY int, oriW int, oriH int) []uint32 {
	xm1 := oriX - 1
	if xm1 < 0 {
		xm1 = 0
	}
	xm2 := oriX - 2
	if xm2 < 0 {
		xm2 = 0
	}
	xp1 := oriX + 1
	if xp1 >= oriW {
		xp1 = oriW - 1
	}
	xp2 := oriX + 2
	if xp2 >= oriW {
		xp2 = oriW - 1
	}
	ym1 := oriY - 1
	if ym1 < 0 {
		ym1 = 0
	}
	ym2 := oriY - 2
	if ym2 < 0 {
		ym2 = 0
	}
	yp1 := oriY + 1
	if yp1 >= oriH {
		yp1 = oriH - 1
	}
	yp2 := oriY + 2
	if yp2 >= oriH {
		yp2 = oriH - 1
	}

	return []uint32{
		oriPixelView[xm1+ym2*oriW],   // a1
		oriPixelView[oriX+ym2*oriW],  // b1
		oriPixelView[xp1+ym2*oriW],   // c1
		oriPixelView[xm2+ym1*oriW],   // a0
		oriPixelView[xm1+ym1*oriW],   // pa
		oriPixelView[oriX+ym1*oriW],  // pb
		oriPixelView[xp1+ym1*oriW],   // pc
		oriPixelView[xp2+ym1*oriW],   // c4
		oriPixelView[xm2+oriY*oriW],  // d0
		oriPixelView[xm1+oriY*oriW],  // pd
		oriPixelView[oriX+oriY*oriW], // pe
		oriPixelView[xp1+oriY*oriW],  // pf
		oriPixelView[xp2+oriY*oriW],  // f4
		oriPixelView[xm2+yp1*oriW],   // g0
		oriPixelView[xm1+yp1*oriW],   // pg
		oriPixelView[oriX+yp1*oriW],  // ph
		oriPixelView[xp1+yp1*oriW],   // pi
		oriPixelView[xp2+yp1*oriW],   // i4
		oriPixelView[xm1+yp2*oriW],   // g5
		oriPixelView[oriX+yp2*oriW],  // h5
		oriPixelView[xp1+yp2*oriW],   // i5
	}
}

func imageToUint32Array(img image.Image) []uint32 {
	bounds := img.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	pixelArray := make([]uint32, width*height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixelArray[y*width+x] = (r>>8)<<24 | (g>>8)<<16 | (b>>8)<<8 | (a >> 8)
		}
	}
	return pixelArray
}

func intArrayToImage(data []uint32, width int, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			idx := y*width + x
			pixel := data[idx]
			r := uint8((pixel >> 24) & 0xFF)
			g := uint8((pixel >> 16) & 0xFF)
			b := uint8((pixel >> 8) & 0xFF)
			a := uint8(pixel & 0xFF)

			offset := img.PixOffset(x, y)
			img.Pix[offset] = r
			img.Pix[offset+1] = g
			img.Pix[offset+2] = b
			img.Pix[offset+3] = a
		}
	}
	return img
}
