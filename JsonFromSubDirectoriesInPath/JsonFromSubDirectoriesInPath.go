package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	// targetPath = "C:/CG tachi [~2012]/Forever"
	categoryPrefix = "CG tachi/Foreve"
	// targetPath     = "/Users/inabamizuki/Desktop/iOS/src/lolimizuki/Projects"
)

func main() {
	subDirectoriesPaths := SubDirectoriesPathsFromPath(targetPath)

	for _, subPath := range subDirectoriesPaths {
		fmt.Println("{")
		fmt.Printf("\"category\": \"%s\",\n", CategoryFromPath(subPath))
		fmt.Printf("\"path\": \"%s\"\n", subPath)

		fmt.Println("},")
	}

	fmt.Println("任意鍵結束這一切 = =")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}

func SubDirectoriesPathsFromPath(path string) []string {
	subComponents, _ := ioutil.ReadDir(path)
	subDirectoryPaths := []string{}

	for _, s := range subComponents {
		if !s.IsDir() {
			continue
		}
		subDirectoryPath := path + "/" + s.Name()
		subDirectoryPaths = append(subDirectoryPaths, subDirectoryPath)
	}

	return subDirectoryPaths
}

func CategoryFromPath(path string) string {
	components := strings.Split(path, "/")
	lastIndex := len(components) - 1

	category := categoryPrefix + "/" + components[lastIndex]

	return category
}
