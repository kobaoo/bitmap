package tools

import "fmt"

type Filter struct {
	IsBlue     bool
	IsGreen    bool
	IsRed      bool
	IsNegative bool
	IsPixelate bool
	IsBlur     bool
}

func (bm *Bitmap) Filter(newfilename string, f *Filter) error {
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

	writeRow := func(row []byte) error {
		// Process each pixel in the row
		for i := 0; i < len(row); i += btsPerPx {
			// Apply filters
			if f.IsBlue {
				// Retain only the blue channel
				row[i+1] = 0 // Green channel
				row[i+2] = 0 // Red channel
			}
			if f.IsRed {
				// Retain only the red channel
				row[i] = 0   // Blue channel
				row[i+1] = 0 // Green channel
			}
			if f.IsGreen {
				// Retain only the green channel
				row[i+2] = 0 // Blue channel
				row[i] = 0   // Red channel
			}
			if f.IsNegative {
				// Apply negative filter
				row[i] = 255 - row[i]     // Red channel
				row[i+1] = 255 - row[i+1] // Green channel
				row[i+2] = 255 - row[i+2] // Blue channel
			}
		}
		// Write the processed row to the file
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
	// Apply pixelation if enabled
	if f.IsPixelate {
		bm.pixelate(20) // Default block size of 20
	}

	// Apply blur if enabled
	if f.IsBlur {
		bm.blur(25)
	}

	// Write all rows to the file
	for y := 0; y < h; y++ {
		rowStart := y * rowSize
		rowEnd := rowStart + rowSize
		row := pxdata[rowStart:rowEnd]
		if err := writeRow(row); err != nil {
			return err
		}
	}

	return nil
}

// Helper function to apply pixelation
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

// Helper function to apply blur with a configurable kernel size and edge handling
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
