package telnetter

import "fmt"

var version string = "v0.2.1"

func PrintVersion() {
	fmt.Println("Telnetter version:", version)
	fmt.Println("github.com/rickcollette/telnetter")
}