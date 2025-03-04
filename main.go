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
		fname := os.Args[2]

		bm, err := tools.LoadBitmap(fname)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		bm.Header.Print()

	default:
		fmt.Println("Unknown command:", command)
		os.Exit(1)
	}
}
