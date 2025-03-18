package tools

import (
	"fmt"
)

func (bm *Bitmap) Copy(newfilename string) error {
	file, err := bm.copyHeader(newfilename)
	if err != nil {
		return err
	}
	h := int(bm.Px.H)
	rowSize := int(bm.Px.RowSize)
	padSize := int(bm.Px.PadSize)
	pad := bm.Px.Pad
	pxdata := bm.Px.Data

	for y := 0; y < h; y++ { // BMP stores rows bottom-to-top
		rowStart := y * rowSize
		rowEnd := rowStart + rowSize
		row := pxdata[rowStart:rowEnd]

		// Write the row data
		_, err = file.Write(row)
		if err != nil {
			return fmt.Errorf("failed to write pixel data to BMP file: %w", err)
		}

		// Write padding bytes if necessary
		if padSize > 0 {
			_, err = file.Write(pad)
			if err != nil {
				return fmt.Errorf("failed to write padding to BMP file: %w", err)
			}
		}
	}

	return nil
}
