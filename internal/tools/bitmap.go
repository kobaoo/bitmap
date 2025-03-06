package tools

import (
	"os"
)

type Bitmap struct {
	H  *BMPHeader
	Px *Pixel
}

type Pixel struct {
	Data       []byte
	BytesPerPx uint16
	W          uint16
	H          uint16
	RowSize    uint16
	PadSize    uint16
	Pad        []byte
}

func NewPixel(data *[]byte, bitsPerPx, w uint16, h int16, imgSize uint16) *Pixel {
	if h < 0 {
		h = -h
	}
	bytesPerPx := bitsPerPx / 8
	rowSize := w * bytesPerPx
	padSize := (4 - (rowSize % 4)) % 4 // Padding to make row size a multiple of 4
	pad := make([]byte, padSize)
	return &Pixel{Data: *data, BytesPerPx: bytesPerPx, W: w, H: uint16(h), RowSize: rowSize, PadSize: padSize, Pad: pad}
}

func LoadBitmap(fname string) (*Bitmap, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bm := new(Bitmap)
	bh, err := ReadImageHeader(file, fname)
	if err != nil {
		return nil, err
	}
	pixels, err := bh.ReadImagePixels(file)
	if err != nil {
		return nil, err
	}
	bm.H = bh
	bm.Px = NewPixel(&pixels, bh.BitsPerPx, bh.W, bh.H, uint16(bh.ImgSize))
	return bm, nil
}
