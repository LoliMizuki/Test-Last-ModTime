package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	path := os.Args[0]
	info, _ := os.Stat(path)

	lastModTime := info.ModTime()

	fmt.Println(lastModTime)

	AnyKeyToExit()
}

func AnyKeyToExit() {
	fmt.Println("任意鍵結束這一切 = =")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}
