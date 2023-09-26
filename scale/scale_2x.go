package scale

import (
	"image"
)

func left2X(n3, n2, pixel uint32, blendColors bool) []uint32 {
	return []uint32{alphaBlend192W(n3, pixel, blendColors), alphaBlend64W(n2, pixel, blendColors)}
}

func up2X(n3, n1, pixel uint32, blendColors bool) []uint32 {
	return []uint32{alphaBlend192W(n3, pixel, blendColors), alphaBlend64W(n1, pixel, blendColors)}
}

func dia2X(n3, pixel uint32, blendColors bool) uint32 {
	return alphaBlend128W(n3, pixel, blendColors)
}

func kernel2Xv5(pe, pi, ph, pf, pg, pc, pd, pb, f4, i4, h5, i5, n1, n2, n3 uint32, blendColors, scaleAlpha bool) (uint32, uint32, uint32) {
	ex := pe != ph && pe != pf
	if !ex {
		return n1, n2, n3
	}

	e := yuvDifference(pe, pc, scaleAlpha) + yuvDifference(pe, pg, scaleAlpha) + yuvDifference(pi, h5, scaleAlpha) + yuvDifference(pi, f4, scaleAlpha) + 4.0*yuvDifference(ph, pf, scaleAlpha)
	i := yuvDifference(ph, pd, scaleAlpha) + yuvDifference(ph, i5, scaleAlpha) + yuvDifference(pf, i4, scaleAlpha) + yuvDifference(pf, pb, scaleAlpha) + 4.0*yuvDifference(pe, pi, scaleAlpha)

	var px uint32
	if yuvDifference(pe, pf, scaleAlpha) <= yuvDifference(pe, ph, scaleAlpha) {
		px = pf
	} else {
		px = ph
	}

	if e < i && (!isEqual(pf, pb, scaleAlpha) && !isEqual(ph, pd, scaleAlpha) || isEqual(pe, pi, scaleAlpha) && (!isEqual(pf, i4, scaleAlpha) && !isEqual(ph, i5, scaleAlpha)) || isEqual(pe, pg, scaleAlpha) || isEqual(pe, pc, scaleAlpha)) {
		ke := yuvDifference(pf, pg, scaleAlpha)
		ki := yuvDifference(ph, pc, scaleAlpha)
		ex2 := pe != pc && pb != pc
		ex3 := pe != pg && pd != pg

		if (ke*2 <= ki && ex3) || (ke >= ki*2 && ex2) {
			if ke*2 <= ki && ex3 {
				leftOut := left2X(n3, n2, px, blendColors)
				n3 = leftOut[0]
				n2 = leftOut[1]
			}
			if ke >= ki*2 && ex2 {
				upOut := up2X(n3, n1, px, blendColors)
				n3 = upOut[0]
				n1 = upOut[1]
			}
		} else {
			n3 = dia2X(n3, px, blendColors)
		}
	} else if e <= i {
		n3 = alphaBlend64W(n3, px, blendColors)
	}

	return n1, n2, n3
}

func computeXbr2x(oriPixelView []uint32, oriX int, oriY int, oriW int, oriH int, dstPixelView []uint32, dstX int, dstY int, dstW int, blendColors bool, scaleAlpha bool) {
	relatedPoints := getRelatedPoints(oriPixelView, oriX, oriY, oriW, oriH)

	a1 := relatedPoints[0]
	b1 := relatedPoints[1]
	c1 := relatedPoints[2]
	a0 := relatedPoints[3]
	pa := relatedPoints[4]
	pb := relatedPoints[5]
	pc := relatedPoints[6]
	c4 := relatedPoints[7]
	d0 := relatedPoints[8]
	pd := relatedPoints[9]
	pe := relatedPoints[10]
	pf := relatedPoints[11]
	f4 := relatedPoints[12]
	g0 := relatedPoints[13]
	pg := relatedPoints[14]
	ph := relatedPoints[15]
	pi := relatedPoints[16]
	i4 := relatedPoints[17]
	g5 := relatedPoints[18]
	h5 := relatedPoints[19]
	i5 := relatedPoints[20]

	e0, e1, e2, e3 := pe, pe, pe, pe

	e1, e2, e3 = kernel2Xv5(pe, pi, ph, pf, pg, pc, pd, pb, f4, i4, h5, i5, e1, e2, e3, blendColors, scaleAlpha)
	e0, e3, e1 = kernel2Xv5(pe, pc, pf, pb, pi, pa, ph, pd, b1, c1, f4, c4, e0, e3, e1, blendColors, scaleAlpha)
	e2, e1, e0 = kernel2Xv5(pe, pa, pb, pd, pc, pg, pf, ph, d0, a0, b1, a1, e2, e1, e0, blendColors, scaleAlpha)
	e3, e0, e2 = kernel2Xv5(pe, pg, pd, ph, pa, pi, pb, pf, h5, g5, d0, g0, e3, e0, e2, blendColors, scaleAlpha)

	dstPixelView[dstX+dstY*dstW] = e0
	dstPixelView[dstX+1+dstY*dstW] = e1
	dstPixelView[dstX+(dstY+1)*dstW] = e2
	dstPixelView[dstX+1+(dstY+1)*dstW] = e3
}

func xbr2X(img image.Image, blend, alpha bool) image.Image {
	pixelArray := imageToUint32Array(img)
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y
	scaledPixelArray := make([]uint32, width*2*height*2)
	for c := 0; c < width; c++ {
		for d := 0; d < height; d++ {
			computeXbr2x(pixelArray, c, d, width, height, scaledPixelArray, c*2, d*2, width*2, blend, alpha)
		}
	}
	return intArrayToImage(scaledPixelArray, width*2, height*2)
}
