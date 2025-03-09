package main

import (
	"fmt"

	"platform.alem.school/git/kseipoll/bitmap/internal/flags"
	"platform.alem.school/git/kseipoll/bitmap/internal/tools"
)

func main() {
	// Parse flags
	cfg := flags.ReadFlags()

	// Load the bitmap
	bm, err := tools.LoadBitmap(cfg.Filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Handle commands
	switch cfg.Command {
	case "header":
		bm.H.Print()
	case "copy":
		err := bm.Copy(cfg.NewFileName)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	case "apply":
		if len(cfg.MirrorType) != 0 {
			bm.Mirror(cfg.NewFileName, cfg.MirrorType)
		}
		if len(cfg.FilterType) != 0 {
			bm.Filter(cfg.NewFileName, cfg.FilterType)
		}
		if len(cfg.RotateType) != 0 {
			err := bm.Rotate(cfg.NewFileName, cfg.RotateType)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
		}

		if len(cfg.CropParams) != 0 {
			err := bm.Crop(cfg.NewFileName, cfg.CropParams)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
		}

	default:
		fmt.Println("Unknown command:", cfg.Command)
	}
}
