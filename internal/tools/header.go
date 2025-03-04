package tools

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

type BMPHeader struct {
	Signature        string // "BM"
	FileSize         uint32 // Size of the file in bytes
	Reserved         uint32 // Unused
	PixelArrayOffset uint32 // Offset to the pixel array
	DIBHeaderSize    uint32 // Size of the DIB header
	Width            uint32 // Width of the image
	Height           uint32 // Height of the image
	ColorPlanes      uint32 // Number of color planes
	BitsPerPixel     uint32 // Bits per pixel
	Compression      uint32 // Compression method
	ImageSize        uint32 // Size of the raw pixel data
	XResolution      uint32 // Horizontal resolution (pixels per meter)
	YResolution      uint32 // Vertical resolution (pixels per meter)
	ColorsInPalette  uint32 // Number of colors in the palette
	ImportantColors  uint32 // Number of important colors
}

func ReadImageHeader(r io.Reader) (BMPHeader, error) {
	bh := BMPHeader{}

	// Use a buffered reader to read larger chunks
	reader := bufio.NewReader(r)
	buf := make([]byte, 54) // Read in 1 KB chunks for faster I/O
	builder := strings.Builder{}

	// Read the 'BM' signature for the BMP file
	for {
		n, err := reader.Read(buf) // Read one byte at a time to check for 'BM'
		if err != nil {
			if err == io.EOF {
				return bh, nil
			}
			return BMPHeader{}, err
		}
		if n > 0 {
			// Process each byte in the buffer
			for i := n - 1; i >= 0; i-- {
				// Convert the byte to hexadecimal and append to the builder
				builder.WriteString(toHex(int(buf[i])))
			}
		}
		builderLen := len(builder.String())
		signature := builder.String()[builderLen-4:]
		if signature != "4d42" {
			return BMPHeader{}, errors.New("Given image file format is not .bmp")
		}
		bh.Signature = "BM"
		builderLen -= 4

		fileSize := builder.String()[builderLen-8 : builderLen]
		bh.FileSize = toDec(fileSize)
		builderLen -= 8

		reserved := builder.String()[builderLen-8 : builderLen]
		bh.Reserved = toDec(reserved)
		builderLen -= 8

		pixelArrayOffset := builder.String()[builderLen-8 : builderLen]
		bh.PixelArrayOffset = toDec(pixelArrayOffset)
		builderLen -= 8

		dIBHeaderSize := builder.String()[builderLen-8 : builderLen]
		bh.DIBHeaderSize = toDec(dIBHeaderSize)
		builderLen -= 8

		height := builder.String()[builderLen-8 : builderLen]
		bh.Height = toDec(height)
		builderLen -= 8

		width := builder.String()[builderLen-8 : builderLen]
		bh.Width = toDec(width)
		builderLen -= 8

		colorPlanes := builder.String()[builderLen-4 : builderLen]
		bh.ColorPlanes = toDec(colorPlanes)
		builderLen -= 4

		bitsPerPixel := builder.String()[builderLen-4 : builderLen]
		bh.BitsPerPixel = toDec(bitsPerPixel)
		builderLen -= 4

		compression := builder.String()[builderLen-8 : builderLen]
		bh.Compression = toDec(compression)
		builderLen -= 8

		imageSize := builder.String()[builderLen-8 : builderLen]
		bh.ImageSize = toDec(imageSize)
		builderLen -= 8

		xResolution := builder.String()[builderLen-8 : builderLen]
		bh.XResolution = toDec(xResolution)
		builderLen -= 8

		yResolution := builder.String()[builderLen-8 : builderLen]
		bh.YResolution = toDec(yResolution)
		builderLen -= 8

		colorsInPalette := builder.String()[builderLen-8 : builderLen]
		bh.ColorsInPalette = toDec(colorsInPalette)
		builderLen -= 8

		importantColors := builder.String()[builderLen:builderLen]
		bh.ImportantColors = toDec(importantColors)

		return bh, nil
	}
}

func (bh *BMPHeader) Print() {
	fmt.Println("BMP Header:")
	fmt.Printf("- FileType %s\n", string(bh.Signature[:]))
	fmt.Printf("- FileSize %d\n", bh.FileSize)
	fmt.Printf("- Reserved %d\n", bh.Reserved)
	fmt.Printf("- PixelArrayOffset %d\n", bh.PixelArrayOffset)
	fmt.Printf("- DIBHeaderSize %d\n", bh.DIBHeaderSize)
	fmt.Printf("- Width %d\n", bh.Width)
	fmt.Printf("- Height %d\n", bh.Height)
	fmt.Printf("- ColorPlanes %d\n", bh.ColorPlanes)
	fmt.Printf("- BitsPerPixel %d\n", bh.BitsPerPixel)
	fmt.Printf("- Compression %d\n", bh.Compression)
	fmt.Printf("- ImageSize %d\n", bh.ImageSize)
	fmt.Printf("- XResolution %d\n", bh.XResolution)
	fmt.Printf("- YResolution %d\n", bh.YResolution)
	fmt.Printf("- ColorsInPalette %d\n", bh.ColorsInPalette)
	fmt.Printf("- ImportantColors %d\n", bh.ImportantColors)
}
