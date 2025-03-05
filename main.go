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
		bm.Header.Print()
	case "copy":
		err := bm.Copy()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	default:
		fmt.Println("Unknown command:", cfg.Command)
	}
}
