package tools

import "fmt"

func (bm *Bitmap) Mirror(newfilename string, commands []string) error {
	for _, command := range commands {
		h := int(bm.Px.H)
		btsPerPx := int(bm.Px.BytesPerPx)
		rowSize := int(bm.Px.RowSize) + int(bm.Px.PadSize) // Include padding in row size
		pxdata := bm.Px.Data
		mirrored := make([]byte, len(pxdata))
		copy(mirrored, pxdata)

		switch command {
		case "h", "hor", "horizontally", "horizontal":
			// Horizontal mirroring
			for y := 0; y < h; y++ {
				rowStart := y * rowSize
				rowEnd := rowStart + rowSize
				row := mirrored[rowStart:rowEnd]

				// Flip pixels in the row (ignore padding)
				for x := 0; x < int(bm.Px.W)/2; x++ {
					leftOffset := x * btsPerPx
					rightOffset := (int(bm.Px.W) - 1 - x) * btsPerPx

					// Swap pixels
					for j := 0; j < btsPerPx; j++ {
						row[leftOffset+j], row[rightOffset+j] = row[rightOffset+j], row[leftOffset+j]
					}
				}
			}

		case "v", "ver", "vertically", "vertical":
			// Vertical mirroring
			for y := 0; y < h/2; y++ {
				srcRowStart := y * rowSize
				dstRowStart := (h - 1 - y) * rowSize

				// Swap entire rows (including padding)
				for i := 0; i < rowSize; i++ {
					mirrored[srcRowStart+i], mirrored[dstRowStart+i] = mirrored[dstRowStart+i], mirrored[srcRowStart+i]
				}
			}

		default:
			return fmt.Errorf("invalid mirroring type: %v", command)
		}

		// Update the original pixel data
		copy(pxdata, mirrored)
	}

	// Update BMP header metadata
	bm.H.ImgSize = uint32(len(bm.Px.Data))
	bm.H.FSize = 54 + bm.H.ImgSize
	bm.H.HeaderSize = 54
	bm.H.DIBHeaderSize = 40

	// Save the modified bitmap to the new file
	return bm.Save(newfilename)
}
