package tools

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
				rowEnd := rowStart + rowSize
				row := mirrored[rowStart:rowEnd]

				// Reverse the row for horizontal mirroring
				for i := 0; i < len(row)/2; i += btsPerPx {
					for j := 0; j < btsPerPx; j++ {
						row[i+j], row[len(row)-i-btsPerPx+j] = row[len(row)-i-btsPerPx+j], row[i+j]
					}
				}
			}
		case "v", "ver", "vertically":
			for y := 0; y < h; y++ {
				srcRowStart := y * rowSize
				dstRowStart := (h - 1 - y) * rowSize
				copy(mirrored[dstRowStart:dstRowStart+rowSize], pxdata[srcRowStart:srcRowStart+rowSize])
			}
		}
		copy(pxdata, mirrored)
	}

	bm.H.ImgSize = uint32(len(bm.Px.Data))
	bm.H.FSize = 54 + bm.H.ImgSize
	bm.H.HeaderSize = 54
	bm.H.DIBHeaderSize = 40

	return bm.Save(newfilename)
}
