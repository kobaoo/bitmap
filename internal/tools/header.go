package tools

type BitmapHeader struct {
	FileType   []byte
	FileSize   uint32
	HeaderSize uint32
}
type DIBHeader struct {
	HeaderSize uint32
	Width      uint32
	Height     uint32
	PixelSize  uint32
	ImageSize  uint32
}
