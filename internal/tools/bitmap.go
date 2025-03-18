package tools

import (
	"fmt"
	"log"
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

func NewPixel(data *[]byte, bitsPerPx, w uint16, h int16, imgSize uint32) *Pixel {
	if h < 0 {
		h = -h
	}
	bytesPerPx := bitsPerPx / 8
	rowSize := w * bytesPerPx
	padSize := (4 - (rowSize % 4)) % 4
	expectedSize := uint32(rowSize+padSize) * uint32(h)
	if uint32(len(*data)) < expectedSize {
		log.Fatal("Not enough data for image")
	}

	return &Pixel{
		Data:       *data,
		BytesPerPx: bytesPerPx,
		W:          w,
		H:          uint16(h),
		RowSize:    rowSize,
		PadSize:    padSize,
	}
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
	bm.Px = NewPixel(&pixels, bh.BitsPerPx, bh.W, bh.H, bh.ImgSize)
	return bm, nil
}

func (bm *Bitmap) Save(fname string) error {
	if len(fname) < 5 {
		log.Fatal("Error: too short name of the files")
	}
	if fname[len(fname)-4:] != ".bmp" {
		log.Fatal("Error: not bmp format")
	}
	file, err := os.Create(fname)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", fname, err)
	}
	defer file.Close()
	// Обновляем заголовок перед записью
	bm.H.ImgSize = uint32(len(bm.Px.Data)) // Размер пиксельных данных
	bm.H.FSize = 54 + bm.H.ImgSize         // Общий размер файла (заголовок + пиксели)
	bm.H.HeaderSize = 54                   // Смещение до пикселей
	bm.H.DIBHeaderSize = 40                // Размер DIB-заголовка (обычно 40)

	// Записываем заголовок BMP
	err = bm.H.Write(file)
	if err != nil {
		return fmt.Errorf("failed to write BMP header: %w", err)
	}

	// Записываем пиксельные данные
	_, err = file.Write(bm.Px.Data)
	if err != nil {
		return fmt.Errorf("failed to write pixel data: %w", err)
	}

	return nil
}
