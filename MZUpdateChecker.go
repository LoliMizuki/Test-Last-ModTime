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

const (
	pathInfosPath = "pathInfos.json"
)

func main() {
	checkTime, pathInfos := ParseInfoFromJson(pathInfosPath)
	fmt.Println("Check updates after:", checkTime)
	fmt.Println("")

	fmt.Println("Updates")
	for _, pathInfo := range pathInfos {
		if isUpdate, modTime := IsDirectoryUpdateAfterTime(pathInfo.path, checkTime); isUpdate {
			// fmt.Printf("%s (last: %s)\n", pathInfo.categoty, modTime.Format("2006-01-02"))
			fmt.Printf(" (path: %s)\n", pathInfo.path)
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

	subDirectoriesPaths := SubDirectoriesPathsFromPath(directoryPath)
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
