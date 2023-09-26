package scale

import (
	"image"
)

func leftUp3X(n7, n5, n6, n2, n8, pixel uint32, blendColors bool) [5]uint32 {
	blendedN7 := alphaBlend192W(n7, pixel, blendColors)
	blendedN6 := alphaBlend64W(n6, pixel, blendColors)
	return [5]uint32{blendedN7, blendedN7, blendedN6, blendedN6, pixel}
}

func left3X(n7, n5, n6, n8, pixel uint32, blendColors bool) [4]uint32 {
	return [4]uint32{
		alphaBlend192W(n7, pixel, blendColors),
		alphaBlend64W(n5, pixel, blendColors),
		alphaBlend64W(n6, pixel, blendColors),
		pixel,
	}
}

func up3X(n5, n7, n2, n8, pixel uint32, blendColors bool) [4]uint32 {
	return [4]uint32{
		alphaBlend192W(n5, pixel, blendColors),
		alphaBlend64W(n7, pixel, blendColors),
		alphaBlend64W(n2, pixel, blendColors),
		pixel,
	}
}

func dia3X(n8, n5, n7, pixel uint32, blendColors bool) [3]uint32 {
	return [3]uint32{
		alphaBlend224W(n8, pixel, blendColors),
		alphaBlend32W(n5, pixel, blendColors),
		alphaBlend32W(n7, pixel, blendColors),
	}
}

func kernel3X(pe, pi, ph, pf, pg, pc, pd, pb, f4, i4, h5, i5, n2, n5, n6, n7, n8 uint32, blendColors bool, scaleAlpha bool) (uint32, uint32, uint32, uint32, uint32) {
	if pe != ph && pe != pf {
		return n2, n5, n6, n7, n8
	}

	e := yuvDifference(pe, pc, scaleAlpha) + yuvDifference(pe, pg, scaleAlpha) + yuvDifference(pi, h5, scaleAlpha) + yuvDifference(pi, f4, scaleAlpha) + 4.0*yuvDifference(ph, pf, scaleAlpha)
	i := yuvDifference(ph, pd, scaleAlpha) + yuvDifference(ph, i5, scaleAlpha) + yuvDifference(pf, i4, scaleAlpha) + yuvDifference(pf, pb, scaleAlpha) + 4.0*yuvDifference(pe, pi, scaleAlpha)

	state := (e < i) && (!isEqual(pf, pb, scaleAlpha) && !isEqual(pf, pc, scaleAlpha) || !isEqual(ph, pd, scaleAlpha) && !isEqual(ph, pg, scaleAlpha) || isEqual(pe, pi, scaleAlpha) && (!isEqual(pf, f4, scaleAlpha) && !isEqual(pf, i4, scaleAlpha) || !isEqual(ph, h5, scaleAlpha) && !isEqual(ph, i5, scaleAlpha)) || isEqual(pe, pg, scaleAlpha) || isEqual(pe, pc, scaleAlpha))
	if state {
		ke := yuvDifference(pf, pg, scaleAlpha)
		ki := yuvDifference(ph, pc, scaleAlpha)
		ex2 := pe != pc && pb != pc
		ex3 := pe != pg && pd != pg

		var px uint32
		if yuvDifference(pe, pf, scaleAlpha) <= yuvDifference(pe, ph, scaleAlpha) {
			px = pf
		} else {
			px = ph
		}

		if (ke*2 <= ki) && ex3 && (ke >= ki*2) && ex2 {
			result := leftUp3X(n7, n5, n6, n2, n8, px, blendColors)
			n7, n5, n6, n2, n8 = result[0], result[1], result[2], result[3], result[4]
		} else if (ke*2 <= ki) && ex3 {
			result := left3X(n7, n5, n6, n8, px, blendColors)
			n7, n5, n6, n8 = result[0], result[1], result[2], result[3]
		} else if ke >= (ki*2) && ex2 {
			result := up3X(n5, n7, n2, n8, px, blendColors)
			n5, n7, n2, n8 = result[0], result[1], result[2], result[3]
		} else {
			result := dia3X(n8, n5, n7, px, blendColors)
			n8, n5, n7 = result[0], result[1], result[2]
		}
	} else if e <= i {
		if yuvDifference(pe, pf, scaleAlpha) <= yuvDifference(pe, ph, scaleAlpha) {
			n8 = alphaBlend128W(n8, pf, blendColors)
		} else {
			n8 = alphaBlend128W(n8, ph, blendColors)
		}
	}
	return n2, n5, n6, n7, n8
}

func computeXbr3x(oriPixelView []uint32, oriX int, oriY int, oriW int, oriH int, dstPixelView []uint32, dstX int, dstY int, dstW int, blendColors, scaleAlpha bool) {
	relatedPoints := getRelatedPoints(oriPixelView, oriX, oriY, oriW, oriH)
	a1, b1, c1, a0, pa, pb, pc, c4, d0, pd, pe, pf, f4, g0, pg, ph, pi, i4, g5, h5, i5 := relatedPoints[0], relatedPoints[1], relatedPoints[2], relatedPoints[3], relatedPoints[4], relatedPoints[5], relatedPoints[6], relatedPoints[7], relatedPoints[8], relatedPoints[9], relatedPoints[10], relatedPoints[11], relatedPoints[12], relatedPoints[13], relatedPoints[14], relatedPoints[15], relatedPoints[16], relatedPoints[17], relatedPoints[18], relatedPoints[19], relatedPoints[20]
	e0, e1, e2, e3, e4, e5, e6, e7, e8 := pe, pe, pe, pe, pe, pe, pe, pe, pe

	e2, e5, e6, e7, e8 = kernel3X(pe, pi, ph, pf, pg, pc, pd, pb, f4, i4, h5, i5, e2, e5, e6, e7, e8, blendColors, scaleAlpha)
	e0, e1, e8, e5, e2 = kernel3X(pe, pc, pf, pb, pi, pa, ph, pd, b1, c1, f4, c4, e0, e1, e8, e5, e2, blendColors, scaleAlpha)
	e6, e3, e2, e1, e0 = kernel3X(pe, pa, pb, pd, pc, pg, pf, ph, d0, a0, b1, a1, e6, e3, e2, e1, e0, blendColors, scaleAlpha)
	e8, e7, e0, e3, e6 = kernel3X(pe, pg, pd, ph, pa, pi, pb, pf, h5, g5, d0, g0, e8, e7, e0, e3, e6, blendColors, scaleAlpha)

	dstPixelView[dstX+dstY*dstW] = e0
	dstPixelView[dstX+1+dstY*dstW] = e1
	dstPixelView[dstX+2+dstY*dstW] = e2
	dstPixelView[dstX+(dstY+1)*dstW] = e3
	dstPixelView[dstX+1+(dstY+1)*dstW] = e4
	dstPixelView[dstX+2+(dstY+1)*dstW] = e5
	dstPixelView[dstX+(dstY+2)*dstW] = e6
	dstPixelView[dstX+1+(dstY+2)*dstW] = e7
	dstPixelView[dstX+2+(dstY+2)*dstW] = e8
}

func xbr3X(img image.Image, blend, alpha bool) image.Image {
	pixelArray := imageToUint32Array(img)
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y
	scaledPixelArray := make([]uint32, width*3*height*3)
	for c := 0; c < width; c++ {
		for d := 0; d < height; d++ {
			computeXbr3x(pixelArray, c, d, width, height, scaledPixelArray, c*3, d*3, width*3, blend, alpha)
		}
	}
	return intArrayToImage(scaledPixelArray, width*3, height*3)
}
