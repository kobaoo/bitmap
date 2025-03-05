package tools

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type BMPHeader struct {
	Signature        string // "BM"
	FileSize         uint32 // Size of the file in bytes
	Reserved         uint32 // Unused
	PixelArrayOffset uint32 // Offset to the pixel array
	DIBHeaderSize    uint32 // Size of the DIB header
	Width            uint32 // Width of the image
	Height           int32  // Height of the image
	ColorPlanes      uint16 // Number of color planes
	BitsPerPixel     uint16 // Bits per pixel
	Compression      uint32 // Compression method
	ImageSize        uint32 // Size of the raw pixel data
	XResolution      uint32 // Horizontal resolution (pixels per meter)
	YResolution      uint32 // Vertical resolution (pixels per meter)
	ColorsInPalette  uint32 // Number of colors in the palette
	ImportantColors  uint32 // Number of important colors
}

func ReadImageHeader(r io.Reader, fname string) (*BMPHeader, error) {
	bh := new(BMPHeader)

	// Use a buffered reader for efficient I/O
	reader := bufio.NewReader(r)

	// Read the entire 54-byte BMP header at once
	buf := make([]byte, 54)
	if _, err := io.ReadFull(reader, buf); err != nil {
		return nil, fmt.Errorf("failed to read BMP header: %w", err)
	}

	// Validate the signature (first 2 bytes)
	if string(buf[:2]) != "BM" {
		return nil, fmt.Errorf("Error: %s is not bitmap file", fname)
	}
	bh.Signature = "BM"

	// Decode fields from the buffer
	bh.FileSize = readUint32(buf[2:6])
	bh.Reserved = readUint32(buf[6:10])
	bh.PixelArrayOffset = readUint32(buf[10:14])
	bh.DIBHeaderSize = readUint32(buf[14:18])
	bh.Width = readUint32(buf[18:22])
	bh.Height = readInt32(buf[22:26])
	bh.ColorPlanes = readUint16(buf[26:28])
	bh.BitsPerPixel = readUint16(buf[28:30])
	bh.Compression = readUint32(buf[30:34])
	bh.ImageSize = readUint32(buf[34:38])
	bh.XResolution = readUint32(buf[38:42])
	bh.YResolution = readUint32(buf[42:46])
	bh.ColorsInPalette = readUint32(buf[46:50])
	bh.ImportantColors = readUint32(buf[50:54])
	return bh, nil
}

func (bh *BMPHeader) ReadImagePixels(r *os.File) ([]byte, error) {
	pixelData := make([]byte, bh.ImageSize)

	_, err := r.Seek(int64(bh.PixelArrayOffset), 0)
	if err != nil {
		return []byte{}, fmt.Errorf("Error: %w seeking to pixel data:", err)
	}
	_, err = r.Read(pixelData)
	if err != nil {
		return []byte{}, fmt.Errorf("Error: %w reading pixel data:", err)
	}
	return pixelData, nil
}

func (bh *BMPHeader) Print() {
	height := bh.Height
	order := "bottom-up"
	if height < 0 {
		order = "top-down"
		height = -height
	}
	fmt.Println("BMP Header:")
	fmt.Printf("- FileType %s\n", string(bh.Signature[:]))
	fmt.Printf("- FileSizeInBytes %d\n", bh.FileSize)
	fmt.Printf("- Reserved %d\n", bh.Reserved)
	fmt.Printf("- HeaderSize %d\n", bh.PixelArrayOffset)
	fmt.Println("DIB Header:")
	fmt.Printf("- DIBHeaderSize %d\n", bh.DIBHeaderSize)
	fmt.Printf("- WidthInPixels %d\n", bh.Width)
	fmt.Printf("- HeightInPixels %d\n", height)
	fmt.Printf("- ColorPlanes %d\n", bh.ColorPlanes)
	fmt.Printf("- PixelSizeInBits %d\n", bh.BitsPerPixel)
	fmt.Printf("- Compression %d\n", bh.Compression)
	fmt.Printf("- ImageSizeInBytes %d\n", bh.ImageSize)
	fmt.Printf("- XResolutionInPixels %d\n", bh.XResolution)
	fmt.Printf("- YResolutionInPixels %d\n", bh.YResolution)
	fmt.Printf("- ColorsInPalette %d\n", bh.ColorsInPalette)
	fmt.Printf("- ImportantColors %d\n", bh.ImportantColors)
	fmt.Printf("- ImageOrder %s\n", order)
}
