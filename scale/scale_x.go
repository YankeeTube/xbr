package scale

import "image"

func Xbr(img image.Image, ratio int) image.Image {
	switch ratio {
	case 4:
		return xbr4X(img)
	case 3:
		return xbr3X(img)
	default:
		return xbr2X(img)
	}
}
