package tools

import (
	"os"
)

type Bitmap struct {
	Header    BMPHeader
	PixelData []byte
}

func LoadBitmap(fname string) (*Bitmap, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bm := new(Bitmap)
	bh, err := ReadImageHeader(file)
	if err != nil {
		return nil, err
	}
	bm.Header = bh
	return bm, nil
}
