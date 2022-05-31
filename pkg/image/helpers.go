package imageprocessing

import (
	"errors"
	"path/filepath"
	"strings"
)

var formatExts = map[string]uint8{
	"jpg":  1,
	"jpeg": 2,
	"png":  3,
	"gif":  4,
	"tif":  5,
	"tiff": 6,
	"bmp":  7,
}

func GetExtensionFromFileName(filename string) (string, error) {
	ext := filepath.Ext(filename)

	formattedExt := strings.ToLower(strings.TrimPrefix(ext, "."))

	if formatExts[formattedExt] == 0 {
		return "", errors.New("Unsupported file")
	}

	return formattedExt, nil
}
