package tools

func (bm *Bitmap) Filter(newfilename string, filterCommands []string) error {
	for _, command := range filterCommands {
		pxdata := bm.Px.Data
		btsPerPx := int(bm.Px.BytesPerPx)
		switch command {
		case "blue":
			for i := 0; i < len(pxdata)-2; i += btsPerPx {
				pxdata[i+1] = 0
				pxdata[i+2] = 0
			}
		case "red":
			for i := 0; i < len(pxdata)-2; i += btsPerPx {
				pxdata[i] = 0
				pxdata[i+1] = 0
			}
		case "green":
			for i := 0; i < len(pxdata)-2; i += btsPerPx {
				pxdata[i] = 0
				pxdata[i+2] = 0
			}
		case "negative":
			for i := 0; i < len(pxdata)-2; i += btsPerPx {
				pxdata[i] = 255 - pxdata[i]
				pxdata[i+1] = 255 - pxdata[i+1]
				pxdata[i+2] = 255 - pxdata[i+2]
			}
		case "pixelate":
			bm.pixelate(10)
		case "blur":
			bm.blur(15)
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
	pxdata := bm.Px.Data

	for y := 0; y < h; y += blockSize {
		for x := 0; x < w; x += blockSize {
			// Calculate the average color in the block
			var r, g, b, count int
			for dy := 0; dy < blockSize && y+dy < h; dy++ {
				for dx := 0; dx < blockSize && x+dx < w; dx++ {
					idx := ((y+dy)*w + (x + dx)) * btsPerPx
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
				// Set all pixels in the block to the average color
				for dy := 0; dy < blockSize && y+dy < h; dy++ {
					for dx := 0; dx < blockSize && x+dx < w; dx++ {
						idx := ((y+dy)*w + (x + dx)) * btsPerPx
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
	pxdata := bm.Px.Data

	// Create a copy of the pixel data to avoid modifying the original during blur
	blurred := make([]byte, len(pxdata))
	copy(blurred, pxdata)

	// Calculate the radius of the kernel
	radius := kernelSize / 2

	// Apply the blur
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := (y*w + x) * btsPerPx
			var r, g, b int
			count := 0

			// Iterate over the kernel
			for dy := -radius; dy <= radius; dy++ {
				for dx := -radius; dx <= radius; dx++ {
					// Calculate the neighbor's coordinates with clamping
					ny := y + dy
					nx := x + dx

					// Clamp the coordinates to stay within the image boundaries
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

					// Get the neighbor's pixel index
					neighborIdx := (ny*w + nx) * btsPerPx

					// Accumulate the color values
					r += int(pxdata[neighborIdx])
					g += int(pxdata[neighborIdx+1])
					b += int(pxdata[neighborIdx+2])
					count++
				}
			}

			// Calculate the average color for the kernel
			blurred[idx] = byte(r / count)
			blurred[idx+1] = byte(g / count)
			blurred[idx+2] = byte(b / count)
		}
	}

	// Replace the original pixel data with the blurred data
	copy(pxdata, blurred)
}
