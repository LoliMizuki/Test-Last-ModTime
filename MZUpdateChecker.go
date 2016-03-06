package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

type PathInfo struct {
	categoty string
	path     string
}

const (
	pathInfosFolderPath = "PathInfos"
)

func main() {
	pathInfosPath := ShowAllPathInfoOptionsAndGetPath()
	CheckAndReport(pathInfosPath)
	AnyKeyToExit()
}

// Main Processes

func ShowAllPathInfoOptionsAndGetPath() string {
	subFilePaths := SubFilesPathsInDirectoryPath(pathInfosFolderPath)

	fmt.Println("選擇一個妳要的:")
	for index, filePath := range subFilePaths {
		fmt.Printf("  (%d) %s\n", index, filePath)
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		inputByte, _ := reader.ReadByte()

		input := string(inputByte)
		readOption, _ := strconv.Atoi(input)

		if 0 <= readOption && readOption < len(subFilePaths) {
			return subFilePaths[readOption]
		} else {
			fmt.Println("再給妳一次機會")
		}
	}

	return ""
}

func CheckAndReport(pathInfosPath string) {
	checkTime, pathInfos := CheckTimeAndParseInfoFromJsonFilePath(pathInfosPath)
	fmt.Println("Check updates after:", checkTime)
	fmt.Println("")

	fmt.Println("Updates")
	for _, pathInfo := range pathInfos {
		if isUpdate, modTime := IsDirectoryUpdateAfterTime(pathInfo.path, checkTime); isUpdate {
			fmt.Printf("%s (last: %s)\n", pathInfo.categoty, modTime.Format("2006-01-02"))
			fmt.Printf(" (path: %s)\n", pathInfo.path)
		}
	}
}

func AnyKeyToExit() {
	fmt.Println("任意鍵結束這一切 = =")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}

//  Compare Functions

func CheckTimeAndParseInfoFromJsonFilePath(jsonPath string) (checkTime time.Time, pathInfos []PathInfo) {
	jsonByte, error := ioutil.ReadFile(jsonPath)
	if error != nil {
		log.Fatalln(error.Error())
	}

	settings := make(map[string]interface{})
	json.Unmarshal(jsonByte, &settings)

	timeString := settings["time"].(string)
	checkTime, _ = time.Parse("2006-01-02", timeString)

	rawPathInfosArray, _ := settings["pathInfos"].([]interface{})

	pathInfos = make([]PathInfo, len(rawPathInfosArray))
	for _, raw := range rawPathInfosArray {
		info := raw.(map[string]interface{})
		pathInfo := PathInfo{info["category"].(string), info["path"].(string)}

		pathInfos = append(pathInfos, pathInfo)
	}

	return
}

func isTime1AfterTime2(time1 time.Time, time2 time.Time) bool {
	return time1.After(time2) || time1.Equal(time2)
}

func IsDirectoryUpdateAfterTime(directoryPath string, checkTime time.Time) (isUpdate bool, modTime time.Time) {
	isUpdate = false
	modTime = time.Now()

	if directoryPath == "" {
		return
	}

	folderInfo, error := os.Stat(directoryPath)
	if error != nil {
		fmt.Println(error.Error())
		return
	}

	subDirectoriesPaths := SubDirectoryPathsInDirectoryPath(directoryPath)
	isMeUpdate := isTime1AfterTime2(folderInfo.ModTime(), checkTime)

	if isMeUpdate && len(subDirectoriesPaths) == 0 {
		isUpdate = true
		modTime = folderInfo.ModTime()
		return
	}

	for _, subDirectoryPath := range subDirectoriesPaths {
		if isSubUpdate, subUpdateTime := IsDirectoryUpdateAfterTime(subDirectoryPath, checkTime); isSubUpdate {
			isUpdate = true
			modTime = subUpdateTime
			return
		}
	}

	return
}

// Support

func SubDirectoryPathsInDirectoryPath(path string) []string {
	subComponents, error := ioutil.ReadDir(path)

	if os.IsNotExist(error) {
		fmt.Println("Error: search path is not exist")
		return make([]string, 0)
	}

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

func SubFilesPathsInDirectoryPath(path string) []string {
	subComponents, error := ioutil.ReadDir(path)

	if os.IsNotExist(error) {
		fmt.Println("Error: search path is not exist")
		return make([]string, 0)
	}

	subFilePaths := []string{}

	for _, s := range subComponents {
		if s.IsDir() {
			continue
		}
		subFilePath := path + "/" + s.Name()
		subFilePaths = append(subFilePaths, subFilePath)
	}

	return subFilePaths
}
