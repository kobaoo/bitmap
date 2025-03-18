package tools

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

type BMPHeader struct {
	Ftype           string // "BM"
	FSize           uint32 // Size of the file in bytes
	Reserved        uint16 // Unused
	HeaderSize      uint16 // Offset to the pixel array
	DIBHeaderSize   uint16 // Size of the DIB header
	W               uint16 // Width of the image
	H               int16  // Height of the image
	ColorPlanes     uint16 // Number of color planes
	BitsPerPx       uint16 // Bits per pixel
	Comp            uint16 // Compression method
	ImgSize         uint32 // Size of the raw pixel data
	XRes            uint16 // Horizontal resolution (pixels per meter)
	YRes            uint16 // Vertical resolution (pixels per meter)
	ColorsInPalette uint32 // Number of colors in the palette
	ImportantColors uint32 // Number of important colors
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

	bh.Ftype = "BM"

	// Decode fields from the buffer
	bh.FSize = readUint32(buf[2:6])
	bh.Reserved = readUint16(buf[6:10])
	bh.HeaderSize = readUint16(buf[10:14])
	if bh.HeaderSize < 54 {
		log.Fatal("Error: BMP image has unsupported header size")
	}
	bh.DIBHeaderSize = readUint16(buf[14:18])
	bh.W = readUint16(buf[18:22])
	bh.H = readInt16(buf[22:26])
	bh.ColorPlanes = readUint16(buf[26:28])
	bh.BitsPerPx = readUint16(buf[28:30])
	if bh.BitsPerPx != 24 {
		log.Fatal("Error: Not 24-bit color")
	}
	bh.Comp = readUint16(buf[30:34])
	if bh.Comp != 0 {
		log.Fatal("Error: Compressed file given")
	}
	bh.ImgSize = readUint32(buf[34:38])
	if bh.ImgSize == 0 {
		bytesPerRow := ((int(bh.W)*int(bh.BitsPerPx) + 31) / 32) * 4
		bh.ImgSize = uint32(bytesPerRow) * uint32(abs(int(bh.H)))
	}
	bh.XRes = readUint16(buf[38:42])
	bh.YRes = readUint16(buf[42:46])
	bh.ColorsInPalette = readUint32(buf[46:50])
	bh.ImportantColors = readUint32(buf[50:54])
	return bh, nil
}

func (bh *BMPHeader) ReadImagePixels(r *os.File) ([]byte, error) {
	pixelData := make([]byte, bh.ImgSize)

	_, err := r.Seek(int64(bh.HeaderSize), 0)
	if err != nil {
		return []byte{}, fmt.Errorf("Error: %w seeking to pixel data:", err)
	}
	_, err = r.Read(pixelData)
	if err != nil {
		return []byte{}, fmt.Errorf("Error: %w reading pixel data:", err)
	}
	return pixelData, nil
}

func (h *BMPHeader) Write(file *os.File) error {
	buf := make([]byte, 54)

	copy(buf[0:2], "BM") // Magic bytes

	writeUint32(buf[2:6], h.FSize)
	writeUint16(buf[6:8], h.Reserved)
	writeUint16(buf[8:10], h.Reserved)
	writeUint16(buf[10:14], h.HeaderSize)
	writeUint16(buf[14:18], h.DIBHeaderSize)
	writeUint16(buf[18:22], h.W)
	writeInt16(buf[22:26], h.H)
	writeUint16(buf[26:28], h.ColorPlanes)
	writeUint16(buf[28:30], h.BitsPerPx)
	writeUint16(buf[30:34], h.Comp)
	writeUint32(buf[34:38], h.ImgSize)
	writeUint16(buf[38:42], h.XRes)
	writeUint16(buf[42:46], h.YRes)
	writeUint32(buf[46:50], h.ColorsInPalette)
	writeUint32(buf[50:54], h.ImportantColors)

	_, err := file.Write(buf)
	if err != nil {
		return fmt.Errorf("failed to write BMP header: %w", err)
	}

	return nil
}

func (bh *BMPHeader) Print() {
	height := bh.H
	order := "bottom-up"
	if height < 0 {
		order = "top-down"
		height = -height
	}
	fmt.Println("BMP Header:")
	fmt.Printf("- FileType %s\n", string(bh.Ftype[:]))
	fmt.Printf("- FileSizeInBytes %d\n", bh.FSize)
	fmt.Printf("- Reserved %d\n", bh.Reserved)
	fmt.Printf("- HeaderSize %d\n", bh.HeaderSize)
	fmt.Println("DIB Header:")
	fmt.Printf("- DIBHeaderSize %d\n", bh.DIBHeaderSize)
	fmt.Printf("- WidthInPixels %d\n", bh.W)
	fmt.Printf("- HeightInPixels %d\n", height)
	fmt.Printf("- ColorPlanes %d\n", bh.ColorPlanes)
	fmt.Printf("- PixelSizeInBits %d\n", bh.BitsPerPx)
	fmt.Printf("- Compression %d\n", bh.Comp)
	fmt.Printf("- ImageSizeInBytes %d\n", bh.ImgSize)
	fmt.Printf("- XResolutionInPixels %d\n", bh.XRes)
	fmt.Printf("- YResolutionInPixels %d\n", bh.YRes)
	fmt.Printf("- ColorsInPalette %d\n", bh.ColorsInPalette)
	fmt.Printf("- ImportantColors %d\n", bh.ImportantColors)
	fmt.Printf("- ImageOrder %s\n", order)
}
