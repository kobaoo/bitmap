package main

import (
	"fmt"
	"os"

	"platform.alem.school/git/kseipoll/bitmap/internal/tools"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: bitmap-tool <command> [options]")
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "header":
		if len(os.Args) < 3 {
			fmt.Println("Usage: bitmap-tool header <filename>")
			os.Exit(1)
		}
		filename := os.Args[2]
		_, err := tools.LoadBitmap(filename)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	default:
		fmt.Println("Unknown command:", command)
		os.Exit(1)
	}
}
