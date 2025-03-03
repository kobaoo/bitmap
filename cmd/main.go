package main

import (
	"fmt"
	"os"

	"platform.alem.school/git/kseipoll/bitmap/internal/utils"
)

func main() {
	args := os.Args
	file, err := os.Open(args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	res, err := utils.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(res[:100])
}
