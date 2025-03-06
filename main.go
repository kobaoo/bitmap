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
		// Validate mirror type
		mirror := new(tools.Mirror)
		for _, mtype := range cfg.MirrorType {
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
	default:
		fmt.Println("Unknown command:", cfg.Command)
	}
}
