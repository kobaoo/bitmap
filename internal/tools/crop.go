package tools

import (
	"fmt"
	"strconv"
	"strings"
)

func (bm *Bitmap) Crop(newfilename string, params []string) error {
	if len(params) != 1 {
		return fmt.Errorf("invalid crop parameters, expected format: xOffset-yOffset-width-height")
	}

	parts := strings.Split(params[0], "-")
	if len(parts) < 2 || len(parts) > 4 {
		return fmt.Errorf("invalid crop parameters, expected format: xOffset-yOffset-width-height")
	}

	xOffset, _ := strconv.Atoi(parts[0])
	yOffset, _ := strconv.Atoi(parts[1])

	newWidth := int(bm.Px.W) - xOffset
	newHeight := int(bm.Px.H) - yOffset

	if len(parts) > 2 {
		newWidth, _ = strconv.Atoi(parts[2])
	}
	if len(parts) > 3 {
		newHeight, _ = strconv.Atoi(parts[3])
	}

	// Проверяем границы
	if xOffset < 0 || yOffset < 0 || newWidth <= 0 || newHeight <= 0 ||
		xOffset+newWidth > int(bm.Px.W) || yOffset+newHeight > int(bm.Px.H) {
		return fmt.Errorf("invalid crop dimensions")
	}

	// Calculate the new row size with padding
	newRowSize := (newWidth*int(bm.Px.BytesPerPx) + 3) &^ 3 // Align to 4 bytes
	newDataSize := newRowSize * newHeight
	newData := make([]byte, newDataSize)

	// Обрезаем изображение
	for y := 0; y < newHeight; y++ {
		srcRowStart := (yOffset + y) * (int(bm.Px.RowSize) + int(bm.Px.PadSize))
		dstRowStart := y * newRowSize

		srcStart := srcRowStart + xOffset*int(bm.Px.BytesPerPx)
		srcEnd := srcStart + newWidth*int(bm.Px.BytesPerPx)

		dstStart := dstRowStart
		dstEnd := dstStart + newWidth*int(bm.Px.BytesPerPx)

		copy(newData[dstStart:dstEnd], bm.Px.Data[srcStart:srcEnd])
	}

	// Update the bitmap headers and pixel data
	bm.Px.Data = newData
	bm.Px.W = uint16(newWidth)
	bm.Px.H = uint16(newHeight)
	bm.Px.RowSize = uint16(newWidth * int(bm.Px.BytesPerPx))
	bm.Px.PadSize = uint16(newRowSize - (newWidth * int(bm.Px.BytesPerPx)))

	bm.H.W = uint16(newWidth)
	bm.H.H = int16(newHeight)
	bm.H.ImgSize = uint32(newDataSize)
	bm.H.FSize = 54 + bm.H.ImgSize

	return bm.Save(newfilename)
}
