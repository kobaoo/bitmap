package tools

import "fmt"

func (bm *Bitmap) Mirror(newfilename string, commands []string) error {
	for _, command := range commands {
		h := int(bm.Px.H)
		rowSize := int(bm.Px.RowSize)
		btsPerPx := int(bm.Px.BytesPerPx)
		pxdata := bm.Px.Data
		mirrored := make([]byte, len(pxdata))
		copy(mirrored, pxdata)

		switch command {
		case "h", "hor", "horizontally":
			for y := 0; y < h; y++ {
				rowStart := y * rowSize
				row := mirrored[rowStart : rowStart+rowSize]

				// Flip pixels in the row
				for i := 0; i < (rowSize-btsPerPx)/2; i += btsPerPx {
					for j := 0; j < btsPerPx; j++ {
						row[i+j], row[rowSize-btsPerPx-i+j] = row[rowSize-btsPerPx-i+j], row[i+j]
					}
				}
			}

		case "v", "ver", "vertically":
			if h > 0 { // Только если BMP хранится снизу вверх
				for y := 0; y < h/2; y++ {
					srcRowStart := y * rowSize
					dstRowStart := (h - 1 - y) * rowSize

					// Swap rows
					for i := 0; i < rowSize; i++ {
						mirrored[srcRowStart+i], mirrored[dstRowStart+i] = mirrored[dstRowStart+i], mirrored[srcRowStart+i]
					}
				}
			}

		default:
			return fmt.Errorf("invalid mirroring type: %v", command)
		}

		copy(pxdata, mirrored)
	}

	bm.H.ImgSize = uint32(len(bm.Px.Data))
	bm.H.FSize = 54 + bm.H.ImgSize
	bm.H.HeaderSize = 54
	bm.H.DIBHeaderSize = 40
	return bm.Save(newfilename)
}
