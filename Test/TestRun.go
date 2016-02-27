package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

var folderPath = "Test Folder"
var innerFolderPath = folderPath + "/" + "Inner Folder"

var file1Path = folderPath + "/" + "file1.txt"
var file2Path = folderPath + "/" + "file2.txt"
var innerFolderFile1Path = innerFolderPath + "/" + "InnerFile1.txt"
var innerFolderFile2Path = innerFolderPath + "/" + "InnerFile2.txt"

func main() {
	TestSetting()
	DoTests()
}

func TestSetting() {
	fmt.Println("Test Setting")
	CleanTestEnvironment()
	CreateTestFolder()
}

func CleanTestEnvironment() {
	if IsPathExist(folderPath) {
		fmt.Println(" > Remvoe exist", folderPath)
		os.RemoveAll(folderPath)
	}
}

func CreateTestFolder() {
	os.Mkdir(folderPath, os.ModeDir|os.ModePerm)
}

func DoTests() {
	ModTimeTest("Add File", AddFile1AndFile2)
	ModTimeTest("Remove File", RemoveFile1)
	// ModTimeTest("Modify File", ModifyFile2) // can not work
	ModTimeTest("Add Inner Folder", AddInnerFolder)
	// ModTimeTest("Add File To Inner Folder", AddFile1And2ToInnerFolder) // depth > 1, can not work
	// ModTimeTest("Remove File From Inner Folder", RemoveFile1FromInnerFolder) // depth > 1, can not work
	// // ModTimeTest("Remvoe Inner Folder", RemoveInnerFolder)
}

func ModTimeTest(title string, testAction func()) {
	sleepSecond := 1 * time.Second

	fmt.Print("Test: ", title)

	preModTime := ModTimeFromPath(folderPath)
	time.Sleep(sleepSecond)

	testAction()

	afterModTime := ModTimeFromPath(folderPath)

	result := preModTime.Before(afterModTime)

	if result {
		fmt.Println(" --> Success")
	} else {
		fmt.Println(" --> Fail")
		fmt.Println(" > pre:", preModTime)
		fmt.Println(" > after:", afterModTime)
	}
}

func AddFile1AndFile2() {
	ioutil.WriteFile(file1Path, nil, os.ModePerm)
	ioutil.WriteFile(file2Path, nil, os.ModePerm)
}

func RemoveFile1() {
	os.Remove(file1Path)
}

func ModifyFile2() {
	file, _ := os.OpenFile(file2Path, os.O_WRONLY|os.O_CREATE, 0666)

	file.WriteString("Fuck you Mizuki")
	file.Sync()
	file.Close()
}

func AddInnerFolder() {
	os.Mkdir(innerFolderPath, os.ModeDir|os.ModePerm)
}

func AddFile1And2ToInnerFolder() {
	ioutil.WriteFile(innerFolderFile1Path, nil, os.ModePerm)
	ioutil.WriteFile(innerFolderFile2Path, nil, os.ModePerm)
}

func RemoveFile1FromInnerFolder() {
	os.RemoveAll(innerFolderFile1Path)
}

func RemoveInnerFolder() {
	os.RemoveAll(innerFolderPath)
}

func IsPathExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func ModTimeFromPath(path string) time.Time {
	info, _ := os.Stat(folderPath)
	return info.ModTime()
}
