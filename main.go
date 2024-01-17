package main

import (
	"os"
)

func main() {
	var Info HostInfo
	err := Flag(&Info)
	if err != nil {
		os.Exit(1)
	}
}
