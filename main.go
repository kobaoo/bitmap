package main

import (
	"fmt"
	"log"

	"platform.alem.school/git/kseipoll/bitmap/internal/flags"
	"platform.alem.school/git/kseipoll/bitmap/internal/tools"
)

func main() {
	// Parse flags
	cfg := flags.ReadFlags()

	// Load the bitmap
	bm, err := tools.LoadBitmap(cfg.Filename)
	if err != nil {
		log.Fatal("Error:", err)
	}

	// Handle commands
	switch cfg.Command {
	case "header":
		bm.H.Print()
	case "copy":
		err := bm.Copy(cfg.NewFileName)
		if err != nil {
			log.Fatal("Error:", err)
		}
	case "apply":
		if len(cfg.MirrorType) != 0 {
			err := bm.Mirror(cfg.NewFileName, cfg.MirrorType)
			if err != nil {
				log.Fatal("Error:", err)
			}
		}
		if len(cfg.FilterType) != 0 {
			err := bm.Filter(cfg.NewFileName, cfg.FilterType)
			if err != nil {
				log.Fatal("Error:", err)
			}
		}
		if len(cfg.RotateType) != 0 {
			err := bm.Rotate(cfg.NewFileName, cfg.RotateType)
			if err != nil {
				log.Fatal("Error:", err)
			}
		}

		if len(cfg.CropParams) != 0 {
			err := bm.Crop(cfg.NewFileName, cfg.CropParams)
			if err != nil {
				log.Fatal("Error:", err)
			}
		}

	default:
		fmt.Println("Unknown command:", cfg.Command)
	}
}
