package telnetter

import "fmt"

var version string = "__VERSION_INFO__"

func PrintVersion() {
	fmt.Println("Telnetter version:", version)
}