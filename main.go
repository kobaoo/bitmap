package main

import (
	"fmt"
	"os"

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
		os.Exit(1)
	}

	// Handle commands
	switch cfg.Command {
	case "header":
		bm.H.Print()
	case "copy":
		err := bm.Copy(cfg.NewFileName)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	case "apply":
		if len(cfg.MirrorType) != 0 {
			err := bm.Mirror(cfg.NewFileName, cfg.MirrorType)
			if err != nil {
				fmt.Println("Error: ", err)
				os.Exit(1)
			}
		}
		if len(cfg.FilterType) != 0 {
			err := bm.Filter(cfg.NewFileName, cfg.FilterType)
			if err != nil {
				fmt.Println("Error: ", err)
				os.Exit(1)
			}
		}
		if len(cfg.RotateType) != 0 {
			err := bm.Rotate(cfg.NewFileName, cfg.RotateType)
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		}

		if len(cfg.CropParams) != 0 {
			err := bm.Crop(cfg.NewFileName, cfg.CropParams)
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		}

	default:
		fmt.Println("Unknown command:", cfg.Command)
	}
}
