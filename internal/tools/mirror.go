package tools

import (
	"fmt"
)

type Mirror struct {
	V bool // Vertical mirror
	H bool // Horizontal mirror
}

func (bm *Bitmap) Mirror(newfilename string, m *Mirror) error {
	file, err := bm.copyHeader(newfilename)
	if err != nil {
		return err
	}
	h := int(bm.Px.H)
	rowSize := int(bm.Px.RowSize)
	btsPerPx := int(bm.Px.BytesPerPx)
	padSize := int(bm.Px.PadSize)
	pad := bm.Px.Pad
	pxdata := bm.Px.Data
	// Helper function to write a row with optional horizontal mirroring
	writeRow := func(row []byte) error {
		if m.H {
			// Reverse the row for horizontal mirroring
			for i := 0; i < len(row)/2; i += btsPerPx {
				for j := 0; j < btsPerPx; j++ {
					row[i+j], row[len(row)-i-btsPerPx+j] = row[len(row)-i-btsPerPx+j], row[i+j]
				}
			}
		}

		_, err := file.Write(row)
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
		return nil
	}

	// Write rows with optional vertical mirroring
	if m.V {
		for y := h - 1; y >= 0; y-- { // BMP stores rows bottom-to-top
			rowStart := y * rowSize
			rowEnd := rowStart + rowSize
			row := pxdata[rowStart:rowEnd]
			if err := writeRow(row); err != nil {
				return err
			}
		}
	} else {
		for y := 0; y < h; y++ { // BMP stores rows bottom-to-top
			rowStart := y * rowSize
			rowEnd := rowStart + rowSize
			row := pxdata[rowStart:rowEnd]
			if err := writeRow(row); err != nil {
				return err
			}
		}
	}

	return nil
}
