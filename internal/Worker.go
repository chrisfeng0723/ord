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
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/spf13/cast"
	"github.com/shopspring/decimal"

	"io/ioutil"
	"ord/utils"
	"regexp"
	"sort"
	"strings"
	"time"
)

type Content struct {
	Key          string
	Value        string
	LocationFile string
}

const PATH = "./data/"

// 定义缓冲区
var contentSlice = make([]Content,0)
var hfSlice = make([]float64, 0)
var fileNumOriginSlice = make([]string, 0)
var hfValueMap = make(map[float64]string, 0)


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


*/

func Worker() {
	//contentChannel := make(chan Content, 10)
	Producer()
	//处理数据
	WriteExcel()

}

// 定义生产者
func Producer() {
	files, _ := ioutil.ReadDir(PATH)
	for _, f := range files {
		fileNumber := utils.GetFileNumber(f.Name())
		fileNumOriginSlice = append(fileNumOriginSlice, fileNumber)
		filePath := PATH + f.Name()
		content := ReadFileContent(filePath)
		//hfSlice = append(hfSlice, utils.GetFileHF(content))
		GetOrdValueFromString(content, fileNumber)

		//HF值处理
		hfValue :=utils.GetFileHF(content)
		r, _ := decimal.NewFromString(hfValue)
		hfFloat, _ := r.Round(5).Float64()
		//有重复的则直接舍弃
		if _, ok := hfValueMap[hfFloat]; !ok {
			//locationSlice = append(locationSlice, cast.ToInt(location))
			hfValueMap[hfFloat] = fileNumber
			hfSlice = append(hfSlice, hfFloat)
		}
	}
	// 生产完毕之后关闭管道
	//close(contentChannel)
	fmt.Println("生产者停止生产")
}

//Molar Mass =    620.4764 grams/mole, [Alpha] ( 6330.0 A) =      196.49 deg.
//找出包含[Alpha]的行，并解析出（）内的内容以及 = 后面的浮点数
func GetOrdValueFromString(content string, fileNumber string) {
	lines := strings.Split(content, "\n")
	//regex1 := `\((.*?)\)`
	regex1 := `\((.*?)\)\s*=\s*([1-9]\d*.\d*|0.\d*[1-9]\d*)\s*deg`
	//regex1 :=`[1-9]\d*.\d*|0.\d*[1-9]\d*`
	reg := regexp.MustCompile(regex1)
	for _, line := range lines {
		if strings.Contains(line, "[Alpha]") {
			//resultSlice := strings.Fields(line)
			//fmt.Println(line)
			temp := reg.FindStringSubmatch(line)
			if len(temp) == 3 {
				column := strings.Replace(temp[1], " ", "", -1)
				line := strings.Replace(temp[2], "  ", "", -1)
				fmt.Println("column: ", column, " line: ", line)
				contentSlice = append(contentSlice,Content{
					Key:          column,
					Value:        line,
					LocationFile: fileNumber,
				})
			}
		}
	}
}

//读取文件内容为字符串
func ReadFileContent(filePath string) string {
	fmt.Println(filePath)
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("open %s failed:%s", filePath, err)
	}
	fmt.Printf("%s 正在处理%s", time.Now().Format("2006-01-02 15:04:05"), filePath)
	fmt.Println("")
	return string(file)
}

//处理数据并写入excel
func WriteExcel() {
	//文件号先去重，然后排序

	singleFileNum := utils.RemoveDuplicate(fileNumOriginSlice)
	sort.Sort(sort.Reverse(sort.StringSlice(singleFileNum)))

	fileName := time.Now().Format("20060102150405") + ".xlsx"
	fmt.Println("数据处理中....")
	f := excelize.NewFile()
	Sheet1 := "Sheet1"
	index := f.NewSheet(Sheet1)
	for key, val := range singleFileNum {
		coordinate, _ := excelize.CoordinatesToCellName(key+2, 1)
		f.SetCellValue(Sheet1, coordinate, val)
	}


	begin :=2
	fileNum :=utils.TransferSliceToMap(singleFileNum)
	existKey := make(map[string]int)
	for _,value := range contentSlice{
		//coordinate, _ := excelize.CoordinatesToCellName()
		var coordinate,firstColumn string
		if location,ok :=existKey[value.Key];ok{
			coordinate, _ = excelize.CoordinatesToCellName(fileNum[value.LocationFile]+2,location)
		}else{
			firstColumn = "A"+cast.ToString(begin)
			coordinate, _ = excelize.CoordinatesToCellName(fileNum[value.LocationFile]+2,begin)
			existKey[value.Key] = begin
			begin++
		}
		err :=f.SetCellValue(Sheet1,firstColumn,value.Key)
		err = f.SetCellValue(Sheet1,coordinate,value.Value)
		if err !=nil{
			fmt.Println(err)
		}

	}



	Sheet2 := "Sheet2"
	sort.Float64s(hfSlice)
	f.NewSheet(Sheet2)
	f.SetCellValue(Sheet2, "B1", "HF")
	countHF := len(fileNum)
	f.SetCellFormula(Sheet2, "F"+cast.ToString(countHF+2), "")
	sumFormula := fmt.Sprintf("SUM(F2:%s)", "F"+cast.ToString(countHF+1))
	f.SetCellFormula(Sheet2, "F"+cast.ToString(countHF+2), sumFormula)
	//fmt.Println(hfValueMap, hfSlice)
	for key, val := range hfSlice {
		yAxis := cast.ToString(key + 2)
		f.SetCellValue(Sheet2, "A"+yAxis, hfValueMap[val])
		f.SetCellValue(Sheet2, "B"+yAxis, val)
		formulaD := "B" + yAxis + "-B2"
		f.SetCellFormula(Sheet2, "C"+yAxis, formulaD)
		formulaE := "C" + yAxis + "*627.5"
		f.SetCellFormula(Sheet2, "D"+yAxis, formulaE)
		formulaF := "-D" + yAxis + "/(0.0019858955*298.15)"
		f.SetCellFormula(Sheet2, "E"+yAxis, formulaF)
		formulaG := "EXP(E" + yAxis + ")"
		f.SetCellFormula(Sheet2, "F"+yAxis, formulaG)
		formulaH := "F" + yAxis + "/F" + cast.ToString(countHF+2)
		f.SetCellFormula(Sheet2, "G"+yAxis, formulaH)

	}


	//计算最后一列的和
	totalFileNum := len(fileNumOriginSlice)
	for i:=2;i<begin;i++{
		//var sumFormula string
		sumFormulaSlice := make([]string,0)
		for key:=range hfSlice{
			location,_ := excelize.CoordinatesToCellName(key+2,i)
			yAxis := cast.ToString(key + 2)
			temp :=fmt.Sprintf("%s*Sheet2!%s",location,"G"+yAxis)
			sumFormulaSlice = append(sumFormulaSlice,temp)
		}
		fmt.Println(strings.Join(sumFormulaSlice,","))
		fmt.Println("----------++++++")
		resultFormula :=fmt.Sprintf("SUM(%s)",strings.Join(sumFormulaSlice,","))
		resultLocation,_ :=excelize.CoordinatesToCellName(totalFileNum+1,i)
		fmt.Println("++++++++++++++++",resultLocation)
		f.SetCellFormula("Sheet1",resultLocation,resultFormula)
	}

	f.SetActiveSheet(index)
	if err := f.SaveAs(fileName); err != nil {
		fmt.Println(err)
	}

	fmt.Println("Excel写入数据完毕...")

}
