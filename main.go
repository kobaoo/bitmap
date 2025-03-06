package main

import (
	"flag"
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

			// Validate mirror type
			mirror := new(tools.Mirror)
			for _, mtype := range cfg.MirrorType {
				if mirror.H && mirror.V {
					break
				}
				switch mtype {
				case "horizontal", "h", "horizontally", "hor":
					mirror.H = true
				case "vertical", "v", "vertically", "ver":
					mirror.V = true
				default:
					fmt.Fprintln(os.Stderr, "Error: invalid mirror type")
					flag.Usage()
					os.Exit(1)
				}
			}
			bm.Mirror(cfg.NewFileName, mirror)
		}
		if len(cfg.FilterType) != 0 {

			// Validate mirror type
			filter := new(tools.Filter)
			for _, ftype := range cfg.FilterType {
				switch ftype {
				case "blue":
					filter.IsBlue = true
				case "red":
					filter.IsRed = true
				case "green":
					filter.IsGreen = true
				case "negative":
					filter.IsNegative = true
				case "pixelate":
					filter.IsPixelate = true
				case "blur":
					filter.IsBlur = true
				default:
					fmt.Fprintln(os.Stderr, "Error: invalid filter type")
					flag.Usage()
					os.Exit(1)
				}
			}
			bm.Filter(cfg.NewFileName, filter)

		}
	default:
		fmt.Println("Unknown command:", cfg.Command)
	}
}
