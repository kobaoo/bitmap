package tools

import (
	"fmt"
)

func (bm *Bitmap) Rotate(newfilename string, angles []string) error {
	for _, angle := range angles {
		switch angle {
		case "right", "90":
			bm.Px = rotateRight(bm.Px)
			bm.H.W, bm.H.H = uint16(bm.H.H), int16(bm.H.W)
		case "left", "-90":
			bm.Px = rotateLeft(bm.Px)
			bm.H.W, bm.H.H = uint16(bm.H.H), int16(bm.H.W)
		case "180", "-180":
			bm.Px = rotate180(bm.Px)
		default:
			return fmt.Errorf("invalid rotation angle: %v", angle)
		}
	}

	bm.H.ImgSize = uint32(len(bm.Px.Data))
	bm.H.FSize = 54 + bm.H.ImgSize
	bm.H.HeaderSize = 54
	bm.H.DIBHeaderSize = 40

	return bm.Save(newfilename)
}

func rotateRight(px *Pixel) *Pixel {
	newData := make([]byte, len(px.Data))
	newPx := &Pixel{
		Data:       newData,
		BytesPerPx: px.BytesPerPx,
		W:          px.H,
		H:          px.W,
		RowSize:    px.H * px.BytesPerPx,
		PadSize:    px.PadSize,
		Pad:        px.Pad,
	}

	for y := 0; y < int(px.H); y++ {
		for x := 0; x < int(px.W); x++ {
			srcIdx := (y*int(px.RowSize) + x*int(px.BytesPerPx))
			dstIdx := ((int(px.W)-1-x)*int(newPx.RowSize) + y*int(px.BytesPerPx))

			copy(newData[dstIdx:dstIdx+int(px.BytesPerPx)], px.Data[srcIdx:srcIdx+int(px.BytesPerPx)])
		}
	}

	newPx.Data = newData
	return newPx
}

func rotateLeft(px *Pixel) *Pixel {
	newData := make([]byte, len(px.Data))
	newPx := &Pixel{
		Data:       newData,
		BytesPerPx: px.BytesPerPx,
		W:          px.H,
		H:          px.W,
		RowSize:    px.H * px.BytesPerPx,
		PadSize:    px.PadSize,
		Pad:        px.Pad,
	}

	for y := 0; y < int(px.H); y++ {
		for x := 0; x < int(px.W); x++ {
			srcIdx := (y*int(px.RowSize) + x*int(px.BytesPerPx))
			dstIdx := (x*int(newPx.RowSize) + (int(px.H)-1-y)*int(px.BytesPerPx))

			copy(newData[dstIdx:dstIdx+int(px.BytesPerPx)], px.Data[srcIdx:srcIdx+int(px.BytesPerPx)])
		}
	}

	newPx.Data = newData
	return newPx
}

func rotate180(px *Pixel) *Pixel {
	newData := make([]byte, len(px.Data))
	newPx := &Pixel{
		Data:       newData,
		BytesPerPx: px.BytesPerPx,
		W:          px.W,
		H:          px.H,
		RowSize:    px.RowSize,
		PadSize:    px.PadSize,
		Pad:        px.Pad,
	}

	for y := 0; y < int(px.H); y++ {
		for x := 0; x < int(px.W); x++ {
			srcIdx := (y*int(px.RowSize) + x*int(px.BytesPerPx))
			dstIdx := ((int(px.H)-1-y)*int(newPx.RowSize) + (int(px.W)-1-x)*int(px.BytesPerPx))
			copy(newData[dstIdx:dstIdx+int(px.BytesPerPx)], px.Data[srcIdx:srcIdx+int(px.BytesPerPx)])
		}
	}
	newPx.Data = newData
	return newPx
}
