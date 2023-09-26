package scale

import "image"

func Xbr(img image.Image, ratio int, blend, alpha bool) image.Image {
	switch ratio {
	case 4:
		return xbr4X(img, blend, alpha)
	case 3:
		return xbr3X(img, blend, alpha)
	default:
		return xbr2X(img, blend, alpha)
	}
}
