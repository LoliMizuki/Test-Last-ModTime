package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type PathInfo struct {
	categoty string
	path     string
}

func main() {
	checkTime, pathInfos := ParseInfoFromJson("pathInfos.json")
	fmt.Println("Check updates after:", checkTime)
	fmt.Println("")

	fmt.Println("Updates")
	for _, pathInfo := range pathInfos {
		if IsDirectoryUpdateAfterTime(pathInfo.path, checkTime) {
			fmt.Printf(" > %s (%s)\n", pathInfo.categoty, pathInfo.path)
		}
	}

	fmt.Println("任意鍵結束這一切 = =")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}

func ParseInfoFromJson(jsonPath string) (checkTime time.Time, pathInfos []PathInfo) {
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
		pathInfo := PathInfo{info["categoty"].(string), info["path"].(string)}

		pathInfos = append(pathInfos, pathInfo)
	}

	return
}

func isTime1AfterTime2(time1 time.Time, time2 time.Time) bool {
	return time1.After(time2) || time1.Equal(time2)
}

func IsDirectoryUpdateAfterTime(directoryPath string, time time.Time) bool {
	if directoryPath == "" {
		return false
	}

	folderInfo, error := os.Stat(directoryPath)
	if error != nil {
		fmt.Println(error.Error())
		return false
	}

	subDirectoriesPaths := SubDirectoriesPathsFromPath(directoryPath)
	isMeUpdate := isTime1AfterTime2(folderInfo.ModTime(), time)

	if isMeUpdate && len(subDirectoriesPaths) == 0 {
		return true
	}

	for _, subDirectoryPath := range subDirectoriesPaths {
		if IsDirectoryUpdateAfterTime(subDirectoryPath, time) {
			return true
		}
	}

	return false
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