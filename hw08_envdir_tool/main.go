package main

import "os"

func main() {
	if err := Execute(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
