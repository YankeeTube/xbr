package scale

import (
	"image"
)

func left4X(n15, n14, n11, n13, n12, n10 uint32, pixel uint32, blendColors bool) [6]uint32 {
	return [6]uint32{
		pixel,
		pixel,
		alphaBlend192W(n11, pixel, blendColors),
		alphaBlend192W(n13, pixel, blendColors),
		alphaBlend64W(n12, pixel, blendColors),
		alphaBlend64W(n10, pixel, blendColors),
	}
}

func up4X(n15, n14, n11, n3, n7, n10, pixel uint32, blendColors bool) [6]uint32 {
	return [6]uint32{
		pixel,
		alphaBlend192W(n14, pixel, blendColors),
		pixel,
		alphaBlend64W(n3, pixel, blendColors),
		alphaBlend192W(n7, pixel, blendColors),
		alphaBlend64W(n10, pixel, blendColors),
	}
}

func dia4X(n15, n14, n11, pixel uint32, blendColors bool) [3]uint32 {
	return [3]uint32{
		pixel,
		alphaBlend128W(n14, pixel, blendColors),
		alphaBlend128W(n11, pixel, blendColors),
	}
}

func kernel4Xv2(pe, pi, ph, pf, pg, pc, pd, pb, f4, i4, h5, i5, n15, n14, n11, n3, n7, n10, n13, n12 uint32, blendColors, scaleAlpha bool) (uint32, uint32, uint32, uint32, uint32, uint32, uint32, uint32) {
	ex := pe != ph && pe != pf
	if !ex {
		return n15, n14, n11, n3, n7, n10, n13, n12
	}

	e := yuvDifference(pe, pc, scaleAlpha) + yuvDifference(pe, pg, scaleAlpha) + yuvDifference(pi, h5, scaleAlpha) + yuvDifference(pi, f4, scaleAlpha) + 4.0*yuvDifference(ph, pf, scaleAlpha)
	i := yuvDifference(ph, pd, scaleAlpha) + yuvDifference(ph, i5, scaleAlpha) + yuvDifference(pf, i4, scaleAlpha) + yuvDifference(pf, pb, scaleAlpha) + 4.0*yuvDifference(pe, pi, scaleAlpha)

	px := ph
	if yuvDifference(pe, pf, scaleAlpha) <= yuvDifference(pe, ph, scaleAlpha) {
		px = pf
	}

	if e < i && (!isEqual(pf, pb, scaleAlpha) && !isEqual(ph, pd, scaleAlpha) || isEqual(pe, pi, scaleAlpha) && (!isEqual(pf, i4, scaleAlpha) && !isEqual(ph, i5, scaleAlpha)) || isEqual(pe, pg, scaleAlpha) || isEqual(pe, pc, scaleAlpha)) {
		ke := yuvDifference(pf, pg, scaleAlpha)
		ki := yuvDifference(ph, pc, scaleAlpha)
		ex2 := pe != pc && pb != pc
		ex3 := pe != pg && pd != pg
		if (ke*2) <= ki && ex3 || ke >= (ki*2) && ex2 {
			if (ke*2) <= ki && ex3 {
				result := left4X(n15, n14, n11, n13, n12, n10, px, blendColors)
				n15, n14, n11, n13, n12, n10 = result[0], result[1], result[2], result[3], result[4], result[5]
			}
			if ke >= (ki*2) && ex2 {
				result := up4X(n15, n14, n11, n3, n7, n10, px, blendColors)
				n15, n14, n11, n3, n7, n10 = result[0], result[1], result[2], result[3], result[4], result[5]
			}
		} else {
			result := dia4X(n15, n14, n11, px, blendColors)
			n15, n14, n11 = result[0], result[1], result[2]
		}
	} else if e <= i {
		n15 = alphaBlend128W(n15, px, blendColors)
	}

	return n15, n14, n11, n3, n7, n10, n13, n12
}

func computeXbr4x(oriPixelView []uint32, oriX, oriY, oriW, oriH int, dstPixelView []uint32, dstX, dstY, dstW int, blendColors bool, scaleAlpha bool) {
	relatedPoints := getRelatedPoints(oriPixelView, oriX, oriY, oriW, oriH)

	// Extract individual pixel values from the related points
	a1, b1, c1, a0, pa, pb, pc, c4, d0, pd, pe, pf, f4, g0, pg, ph, pi, i4, g5, h5, i5 := relatedPoints[0], relatedPoints[1], relatedPoints[2], relatedPoints[3], relatedPoints[4], relatedPoints[5], relatedPoints[6], relatedPoints[7], relatedPoints[8], relatedPoints[9], relatedPoints[10], relatedPoints[11], relatedPoints[12], relatedPoints[13], relatedPoints[14], relatedPoints[15], relatedPoints[16], relatedPoints[17], relatedPoints[18], relatedPoints[19], relatedPoints[20]
	e0, e1, e2, e3, e4, e5, e6, e7, e8, e9, ea, eb, ec, ed, ee, ef := pe, pe, pe, pe, pe, pe, pe, pe, pe, pe, pe, pe, pe, pe, pe, pe

	ef, ee, eb, e3, e7, ea, ed, ec = kernel4Xv2(pe, pi, ph, pf, pg, pc, pd, pb, f4, i4, h5, i5, ef, ee, eb, e3, e7, ea, ed, ec, blendColors, scaleAlpha)
	e3, e7, e2, e0, e1, e6, eb, ef = kernel4Xv2(pe, pc, pf, pb, pi, pa, ph, pd, b1, c1, f4, c4, e3, e7, e2, e0, e1, e6, eb, ef, blendColors, scaleAlpha)
	e0, e1, e4, ec, e8, e5, e2, e3 = kernel4Xv2(pe, pa, pb, pd, pc, pg, pf, ph, d0, a0, b1, a1, e0, e1, e4, ec, e8, e5, e2, e3, blendColors, scaleAlpha)
	ec, e8, ed, ef, ee, e9, e4, e0 = kernel4Xv2(pe, pg, pd, ph, pa, pi, pb, pf, h5, g5, d0, g0, ec, e8, ed, ef, ee, e9, e4, e0, blendColors, scaleAlpha)

	dstPixelView[dstX+dstY*dstW] = e0
	dstPixelView[dstX+1+dstY*dstW] = e1
	dstPixelView[dstX+2+dstY*dstW] = e2
	dstPixelView[dstX+3+dstY*dstW] = e3
	dstPixelView[dstX+(dstY+1)*dstW] = e4
	dstPixelView[dstX+1+(dstY+1)*dstW] = e5
	dstPixelView[dstX+2+(dstY+1)*dstW] = e6
	dstPixelView[dstX+3+(dstY+1)*dstW] = e7
	dstPixelView[dstX+(dstY+2)*dstW] = e8
	dstPixelView[dstX+1+(dstY+2)*dstW] = e9
	dstPixelView[dstX+2+(dstY+2)*dstW] = ea
	dstPixelView[dstX+3+(dstY+2)*dstW] = eb
	dstPixelView[dstX+(dstY+3)*dstW] = ec
	dstPixelView[dstX+1+(dstY+3)*dstW] = ed
	dstPixelView[dstX+2+(dstY+3)*dstW] = ee
	dstPixelView[dstX+3+(dstY+3)*dstW] = ef
}

func xbr4X(img image.Image) image.Image {
	pixelArray := imageToUint32Array(img)
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y
	scaledPixelArray := make([]uint32, width*4*height*4)
	for c := 0; c < width; c++ {
		for d := 0; d < height; d++ {
			computeXbr4x(pixelArray, c, d, width, height, scaledPixelArray, c*4, d*4, width*4, true, true)
		}
	}
	return intArrayToImage(scaledPixelArray, width*4, height*4)
}
