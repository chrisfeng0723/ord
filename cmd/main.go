/**
 * @Author: fxl
 * @Description:
 * @File:  main.go
 * @Version: 1.0.0
 * @Date: 2021/6/19 10:09
 */
package main

import (
	"fmt"
	"io/ioutil"
	"ord/internal"
	"ord/utils"
	"sort"
)

const PATH = "./data/"

func main() {
	//fmt.Println("hello ord")
	resultChan := make(chan internal.Content)
	go Worker(resultChan)
	for temp := range resultChan {
		fmt.Println(temp.Key)
	}
	fmt.Println("finish")
}

func Worker(resultChan chan internal.Content) {
	files, _ := ioutil.ReadDir(PATH)
	fileNumOriginSlice := make([]string, 0, len(files))

	for _, f := range files {
		fmt.Println(f.Name())
		fileNumber := utils.GetFileNumber(f.Name())
		fmt.Println(fileNumber)
		fileNumOriginSlice = append(fileNumOriginSlice, fileNumber)

		filePath := PATH + f.Name()
		content := internal.ReadFileContent(filePath)
		internal.GetOrdValueFromString(resultChan, content, fileNumber)

	}
	close(resultChan)

	//整理文件的序号
	fileNumberSlice := utils.RemoveDuplicate(fileNumOriginSlice)
	// 逆序
	sort.Sort(sort.Reverse(sort.StringSlice(fileNumberSlice)))
	fmt.Println(fileNumberSlice)
	fileNumberMap := utils.TransferSliceToMap(fileNumberSlice)
	fmt.Println(fileNumberMap)

	//整理数据
	//resultMap :=make(map[string][]string)

}
