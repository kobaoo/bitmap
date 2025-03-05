package tools

import (
	"fmt"
	"os"
)

func (bm *Bitmap) Copy() error {
	file, err := os.Create("result.bmp")
	if err != nil {
		return fmt.Errorf("failed to create BMP file: %w", err)
	}
	defer file.Close()

	h := int(bm.Header.Height)
	if h < 0 {
		h = -h
	}
	w := int(bm.Header.Width)

	// Write the BMP header
	header := make([]byte, 54)
	header[0] = 'B'
	header[1] = 'M'
	writeUint32(header[2:6], bm.Header.FileSize)
	writeUint32(header[6:10], bm.Header.Reserved)
	writeUint32(header[10:14], bm.Header.PixelArrayOffset)
	writeUint32(header[14:18], bm.Header.DIBHeaderSize)
	writeUint32(header[18:22], bm.Header.Width)
	writeInt32(header[22:26], bm.Header.Height)
	writeUint16(header[26:28], bm.Header.ColorPlanes)
	writeUint16(header[28:30], bm.Header.BitsPerPixel)
	writeUint32(header[30:34], bm.Header.Compression)
	writeUint32(header[34:38], bm.Header.ImageSize)
	writeUint32(header[38:42], bm.Header.XResolution)
	writeUint32(header[42:46], bm.Header.YResolution)
	writeUint32(header[46:50], bm.Header.ColorsInPalette)
	writeUint32(header[50:54], bm.Header.ImportantColors)

	_, err = file.Write(header)
	if err != nil {
		return fmt.Errorf("failed to write header to BMP file: %w", err)
	}

	// Write pixel data
	bytesPerPixel := int(bm.Header.BitsPerPixel) / 8
	rowSize := w * bytesPerPixel
	paddingSize := (4 - (rowSize % 4)) % 4 // Padding to make row size a multiple of 4
	padding := make([]byte, paddingSize)

	for y := 0; y < h; y++ { // BMP stores rows bottom-to-top
		rowStart := y * rowSize
		rowEnd := rowStart + rowSize
		row := bm.PixelData[rowStart:rowEnd]

		_, err = file.Write(row)
		if err != nil {
			return fmt.Errorf("failed to write pixel data to BMP file: %w", err)
		}

		// Write padding bytes if necessary
		if paddingSize > 0 {
			_, err = file.Write(padding)
			if err != nil {
				return fmt.Errorf("failed to write padding to BMP file: %w", err)
			}
		}
	}

	return nil
}
