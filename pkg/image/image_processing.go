package imageprocessing

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"

	"github.com/disintegration/imaging"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

func CreateBlurDataForImages(buffer []byte, quality int, width int, height int, fileName string) ([]byte, error) {

	bytesReader := bytes.NewReader(buffer)

	img, imageType, err := image.Decode(bytesReader)
	if err != nil {
		return nil, err
	}

	resizedImg := imaging.Resize(img, width, height, imaging.Linear)

	buf := new(bytes.Buffer)

	switch imageType {
	case "jpg", "jpeg":
		err = jpeg.Encode(buf, resizedImg, &jpeg.Options{Quality: 1})
	case "png":
		err = png.Encode(buf, resizedImg)

	case "gif":
		err = gif.Encode(buf, resizedImg, nil)

	case "tif", "tiff":
		err = tiff.Encode(buf, resizedImg, nil)

	case "bmp":
		err = bmp.Encode(buf, resizedImg)

	default:
		err = jpeg.Encode(buf, resizedImg, &jpeg.Options{Quality: 1})
	}

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
