/**
 * @Author: fxl
 * @Description:
 * @File:  Worker.go
 * @Version: 1.0.0
 * @Date: 2021/6/19 10:31
 */
package internal

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"time"
)

type Content struct {
	Key string
	Value string
	LocationFile string
}

/**
func Worker() {
	//resultSlice := make([]schema.Result,0)
	//hf的map值，key为所在的文件，value为hf值
	hfValueMap := make(map[float64]string, 0)
	hfValueSlice := make([]float64, 0)
	//文件编号的slice
	locationSlice := make([]int, 0)
	//所有CH的值
	allResultSlice := make([]schema.Result, 0)
	//C的数值slice
	cResultSlice := make([]int, 0)
	//H的数值slice
	hResultSlice := make([]int, 0)

	files, _ := ioutil.ReadDir(PATH)
	for _, f := range files {
		fmt.Println("正在处理" + f.Name())
		hfValue, location, resultSlice, cSlice, hSlice := GetValueByFileName(f.Name())

		//处理hf值
		r, _ := decimal.NewFromString(hfValue)
		hfFloat, _ := r.Round(5).Float64()
		//有重复的则直接舍弃
		if _, ok := hfValueMap[hfFloat]; !ok {
			locationSlice = append(locationSlice, cast.ToInt(location))
			hfValueMap[hfFloat] = location
			hfValueSlice = append(hfValueSlice, hfFloat)
			//所有结果
			allResultSlice = append(allResultSlice, resultSlice...)
			cResultSlice = append(cResultSlice, cSlice...)
			hResultSlice = append(hResultSlice, hSlice...)
		}

	}

	//排序文件编号，写入hf值
	sort.Ints(locationSlice)
	sort.Float64s(hfValueSlice)
	//排序CH顺序
	uniqueCSlice := utils.RemoveDuplicate(cResultSlice)
	uniqueHSlice := utils.RemoveDuplicate(hResultSlice)
	sort.Ints(uniqueHSlice)
	sort.Ints(uniqueCSlice)

	//fmt.Println(hfValueSlice,hfValueMap)

	WriteExcel(locationSlice, hfValueSlice, hfValueMap, uniqueCSlice, uniqueHSlice, allResultSlice)

}
func GetValueByFileName(fileName string) (hfValue, location string, resultContentSlice []schema.Result, cResultSlice, hResultSlice []int) {
	fileName = PATH + fileName
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("open %s failed:%s", fileName, err)
	}
	//获取文件location
	location = utils.GetFileNumber(fileName)

	//获取hf值
	hfValue = utils.GetFileHF(string(file))

	for _, line := range strings.Split(string(file), "\n") {
		if strings.Contains(line, "Isotropic") && strings.Contains(line, "Anisotropy") {
			resultSlice := strings.Fields(line)
			if resultSlice[1] == "H" || resultSlice[1] == "C" {
				result := schema.Result{
					Location: location,
					Sequence: resultSlice[0],
					Element:  resultSlice[1],
					Value:    resultSlice[4],
				}
				if resultSlice[1] == "H" {
					hResultSlice = append(hResultSlice, cast.ToInt(resultSlice[0]))
				}

				if resultSlice[1] == "C" {
					cResultSlice = append(cResultSlice, cast.ToInt(resultSlice[0]))
				}
				resultContentSlice = append(resultContentSlice, result)
			}

		}

	}

	return

}

*/
//Molar Mass =    620.4764 grams/mole, [Alpha] ( 6330.0 A) =      196.49 deg.
//找出包含[Alpha]的行，并解析出（）内的内容以及 = 后面的浮点数
func GetOrdValueFromString(resultChan chan Content,content,fileNumber string) {
	lines := strings.Split(content, "\n")
	//regex1 := `\((.*?)\)`
	regex1 := `\((.*?)\)\s*=\s*([1-9]\d*.\d*|0.\d*[1-9]\d*)\s*deg`
	//regex1 :=`[1-9]\d*.\d*|0.\d*[1-9]\d*`
	reg := regexp.MustCompile(regex1)
	for _, line := range lines {
		if strings.Contains(line, "[Alpha]") {
			//resultSlice := strings.Fields(line)
			fmt.Println(line)
			temp := reg.FindStringSubmatch(line)
			if len(temp) == 3{
				column := strings.Replace(temp[1]," ","",-1)
				line := strings.Replace(temp[2],"  ","",-1)
				fmt.Println("column: ",column," line: ",line)
				resultChan <- Content{
					Key: column,
					Value: line,
					LocationFile:fileNumber,
				}
			}
		}
	}

}

//读取文件内容为字符串
func ReadFileContent(filePath string) string {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("open %s failed:%s", filePath, err)
	}
	fmt.Printf("%s 正在处理%s", time.Now().Format("2006-01-02 15:04:05"), filePath)
	fmt.Println("")
	return string(file)
}


