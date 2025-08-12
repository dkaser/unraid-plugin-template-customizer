package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("=========================================")
	fmt.Println("   Unraid Plugin Template Setup")
	fmt.Println("=========================================")
	fmt.Println()

	if err := runSetup(); err != nil {
		fmt.Printf("Setup failed: %v\n", err)
		os.Exit(1)
	}
}
