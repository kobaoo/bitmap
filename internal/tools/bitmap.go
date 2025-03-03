package tools

type Bitmap struct {
	Header    BitmapHeader
	DIBHeader DIBHeader
	PixelData []byte
}

func LoadBitmap(filename string) (*Bitmap, error) {
	return nil, nil
}
