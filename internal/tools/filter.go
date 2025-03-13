package tools

import "fmt"

func (bm *Bitmap) Filter(newfilename string, filterCommands []string) error {
	pxdata := bm.Px.Data
	btsPerPx := int(bm.Px.BytesPerPx)
	rowSize := int(bm.Px.RowSize) + int(bm.Px.PadSize) // Размер строки с учётом padding

	for _, command := range filterCommands {
		switch command {
		case "blue":
			for y := 0; y < int(bm.Px.H); y++ {
				for x := 0; x < int(bm.Px.W); x++ {
					offset := y*rowSize + x*btsPerPx
					pxdata[offset+1] = 0    
					pxdata[offset+2] = 0  
				}
			}
		case "red":
			for y := 0; y < int(bm.Px.H); y++ {
				for x := 0; x < int(bm.Px.W); x++ {
					offset := y*rowSize + x*btsPerPx
					pxdata[offset] = 0    
					pxdata[offset+1] = 0   
				}
			}
		case "green":
			for y := 0; y < int(bm.Px.H); y++ {
				for x := 0; x < int(bm.Px.W); x++ {
					offset := y*rowSize + x*btsPerPx
					pxdata[offset] = 0     
					pxdata[offset+2] = 0 
				}
			}
		case "negative":
			for y := 0; y < int(bm.Px.H); y++ {
				for x := 0; x < int(bm.Px.W); x++ {
					offset := y*rowSize + x*btsPerPx
					pxdata[offset] = 255 - pxdata[offset]    
					pxdata[offset+1] = 255 - pxdata[offset+1]  
					pxdata[offset+2] = 255 - pxdata[offset+2]  
				}
			}
		case "pixelate":
			bm.pixelate(10)
		case "blur":
			bm.blur(15)
		default:
			return fmt.Errorf("invalid filter operation: %v", command)
		}
	}
  
	bm.H.ImgSize = uint32(len(bm.Px.Data))
	bm.H.FSize = 54 + bm.H.ImgSize
	bm.H.HeaderSize = 54
	bm.H.DIBHeaderSize = 40
 
	return bm.Save(newfilename)
}
func (bm *Bitmap) pixelate(blockSize int) {
	h := int(bm.Px.H)
	w := int(bm.Px.W)
	btsPerPx := int(bm.Px.BytesPerPx)
	rowSize := int(bm.Px.RowSize) + int(bm.Px.PadSize) // Учитываем padding
	pxdata := bm.Px.Data

	for y := 0; y < h; y += blockSize {
		for x := 0; x < w; x += blockSize { 
			var r, g, b, count int
			for dy := 0; dy < blockSize && y+dy < h; dy++ {
				for dx := 0; dx < blockSize && x+dx < w; dx++ {
					idx := ((y+dy)*rowSize + (x+dx)*btsPerPx) // Учитываем padding
					r += int(pxdata[idx])
					g += int(pxdata[idx+1])
					b += int(pxdata[idx+2])
					count++
				}
			}
			if count > 0 {
				r /= count
				g /= count
				b /= count 
				for dy := 0; dy < blockSize && y+dy < h; dy++ {
					for dx := 0; dx < blockSize && x+dx < w; dx++ {
						idx := ((y+dy)*rowSize + (x+dx)*btsPerPx) // Учитываем padding
						pxdata[idx] = byte(r)
						pxdata[idx+1] = byte(g)
						pxdata[idx+2] = byte(b)
					}
				}
			}
		}
	}
}
func (bm *Bitmap) blur(kernelSize int) {
	h := int(bm.Px.H)
	w := int(bm.Px.W)
	btsPerPx := int(bm.Px.BytesPerPx)
	rowSize := int(bm.Px.RowSize) + int(bm.Px.PadSize) 
	pxdata := bm.Px.Data
 
	blurred := make([]byte, len(pxdata))
	copy(blurred, pxdata)
 
	radius := kernelSize / 2
 
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := (y*rowSize + x*btsPerPx) // Учитываем padding
			var r, g, b int
			count := 0 
			for dy := -radius; dy <= radius; dy++ {
				for dx := -radius; dx <= radius; dx++ { 
					ny := y + dy
					nx := x + dx
 
					if ny < 0 {
						ny = 0
					} else if ny >= h {
						ny = h - 1
					}
					if nx < 0 {
						nx = 0
					} else if nx >= w {
						nx = w - 1
					}
 
					neighborIdx := (ny*rowSize + nx*btsPerPx) 
 
					r += int(pxdata[neighborIdx])
					g += int(pxdata[neighborIdx+1])
					b += int(pxdata[neighborIdx+2])
					count++
				}
			}
 
			blurred[idx] = byte(r / count)
			blurred[idx+1] = byte(g / count)
			blurred[idx+2] = byte(b / count)
		}
	}
 
	copy(pxdata, blurred)
}