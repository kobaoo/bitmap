package tools

import (
	"fmt"
	"os"
)

// Helper function to read a uint32 from a byte slice
func readUint32(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

// Helper function to read an int16 from a byte slice
func readInt16(b []byte) int16 {
	return int16(b[0]) | int16(b[1])<<8
}

// Helper function to read a uint16 from a byte slice
func readUint16(b []byte) uint16 {
	return uint16(b[0]) | uint16(b[1])<<8
}

// Helper function to write a byte slice from uint32
func writeUint32(b []byte, v uint32) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
}

// Helper function to write a byte slice from uint16
func writeInt16(b []byte, v int16) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
}

// Helper function to write a byte slice from uint16
func writeUint16(b []byte, v uint16) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
}

func (bm *Bitmap) copyHeader(newfilename string) (*os.File, error) {
	file, err := os.Create(newfilename)
	if err != nil {
		return nil, fmt.Errorf("failed to create BMP file: %w", err)
	}
	if err = bm.writeHeader(file); err != nil {
		return nil, err
	}
	return file, nil
}

func (bm *Bitmap) writeHeader(file *os.File) error {
	// Write the BMP header
	header := make([]byte, 54)
	header[0] = 'B'
	header[1] = 'M'
	writeUint32(header[2:6], bm.H.FSize)
	writeUint16(header[6:10], bm.H.Reserved)
	writeUint16(header[10:14], bm.H.HeaderSize)
	writeUint16(header[14:18], bm.H.DIBHeaderSize)
	writeUint16(header[18:22], bm.H.W)
	writeInt16(header[22:26], bm.H.H)
	writeUint16(header[26:28], bm.H.ColorPlanes)
	writeUint16(header[28:30], bm.H.BitsPerPx)
	writeUint16(header[30:34], bm.H.Comp)
	writeUint32(header[34:38], bm.H.ImgSize)
	writeUint16(header[38:42], bm.H.XRes)
	writeUint16(header[42:46], bm.H.YRes)
	writeUint32(header[46:50], bm.H.ColorsInPalette)
	writeUint32(header[50:54], bm.H.ImportantColors)

	_, err := file.Write(header)
	if err != nil {
		return fmt.Errorf("failed to write header to BMP file: %w", err)
	}
	return nil
}

func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}
