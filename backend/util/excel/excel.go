package excel

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/XHXHXHX/medical_marketing/util/common"
	"net/http"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"golang.org/x/text/encoding/simplifiedchinese"
)

const UPLOAD_FILE_MAX_SIZE = 220200960

/**
 * 从上传文件中返回文件内容
 * @param req 		*http.Request
 * @param filename 	string		表单字段名
 * @param maxSize 	int64		文件最大字节数
 *
 * @return
 * 		[]byte  	文件内容
 * 		error		错误信息
 */
func GetFileContentFromUploadFile(req *http.Request, filename string, maxSize int64) ([]byte, error) {
	file, fileHeader, err := req.FormFile(filename)
	if err != nil {
		return nil, errors.New("文件错误")
	}

	defer file.Close()
	//headerByte, _ := json.Marshal(fileHeader.Header)
	//fmt.Printf("当前文件：Filename - >{%s}, Size -> {%v}, FileHeader -> {%s}", fileHeader.Filename, fileHeader.Size, string(headerByte))
	if maxSize == 0 {
		maxSize = UPLOAD_FILE_MAX_SIZE
	}
	// 目前设置最大210M 210 * 1024 * 1024 fileHeader.Size是byte
	if fileHeader.Size > maxSize {
		return nil, errors.New("超过文件最大大小")
	}

	buffer := make([]byte, fileHeader.Size)
	row, err := file.Read(buffer)
	if err != nil {
		return nil, errors.New("读取文件错误")
	}

	if row == 0 {
		return nil, errors.New("空文件")
	}

	return buffer, nil
}

/**
 * 导出Excel表格
 * @param  header   {[]string}  表头key，导出后显示的顺序
 * @param  headerKV {map[string]string}  表头、数据kv对照
 * @param  data     {[]map[string]interface{}} 数据集合
 *
 * @return 		 {bytes.Buffer}				缓存数据
 * @return       {error}                    异常
 */
func ExportExcel(header []string, headerKV map[string]string, data []map[string]interface{}) (*bytes.Buffer, error) {
	if len(header) > 676 {
		return nil, errors.New("导出数据列数上限超过限制")
	}

	f := excelize.NewFile()
	// Create a new sheet
	index := f.NewSheet("Sheet1")
	for i, v := range header {
		columnKey := GetColumnKey(i)
		f.SetCellStr("Sheet1", fmt.Sprintf("%s1", columnKey), headerKV[v])
	}

	for i, v := range data {
		for ii, key := range header {

			columnKey := GetColumnKey(ii)
			axis := fmt.Sprintf("%s%d", columnKey, i+2)

			switch val := v[key].(type) {
			case int:
				f.SetCellInt("Sheet1", axis, val)
			case int32:
				f.SetCellInt("Sheet1", axis, int(val))
			case int64:
				f.SetCellInt("Sheet1", axis, int(val))
			case string:
				f.SetCellStr("Sheet1", axis, val)
			case bool:
				f.SetCellBool("Sheet1", axis, val)
			default:
				f.SetCellStr("Sheet", axis, common.InterfaceToString(val))
			}
		}
	}

	f.SetActiveSheet(index)

	return f.WriteToBuffer()
}

/**
 * 导出Excel表格
 * @param  	data   		{[][]string}  	要导出的数据
 * @param	mergeCell 	[][]string		要合并的单元格数组 eg: [["A1", "B1]]
 * @param	filter 		[]string 		要添加过滤器的行 eg:["A1", "J5"]
 * @param	visibleFunc func(item []string) bool 是否隐藏方法，item 是每行数据集  返回true隐藏
 *
 * @return 		 {bytes.Buffer}				缓存数据
 * @return       {error}                    异常
 */
func ExportExcelBySlice(data, mergeCell [][]string, filter []string, visibleFunc func(item []string) bool) (*bytes.Buffer, error) {
	if len(data) == 0 {
		return nil, nil
	}
	f := excelize.NewFile()
	// Create a new sheet
	sheet := "Sheet1"
	index := f.NewSheet(sheet)

	style, err := f.NewStyle(`{"alignment": {"horizontal": "center", "vertical":"center"}}`)
	if err != nil {
		return nil, err
	}

	for i, v := range data[0] {
		columnKey := GetColumnKey(i)
		f.SetCellStr(sheet, fmt.Sprintf("%s1", columnKey), v)
	}

	for i, v := range data {
		if i == 0 {
			continue
		}
		for ii, vv := range v {
			columnKey := GetColumnKey(ii)
			axis := fmt.Sprintf("%s%d", columnKey, i+1)
			f.SetCellStr(sheet, axis, vv)
		}
		if visibleFunc != nil {
			f.SetRowVisible(sheet, i, visibleFunc(v))
		}
	}

	f.SetCellStyle(sheet, "A1", fmt.Sprintf("%s%d", CalculateColVal(len(data[0])), len(data)), style)
	if len(filter) >= 2 {
		err = f.AutoFilter(sheet, filter[0], filter[1], ``)
		if err != nil {
			return nil, err
		}
	}

	if len(mergeCell) > 0 {
		for _, v := range mergeCell {
			if len(v) < 1 {
				continue
			}

			f.MergeCell(sheet, v[0], v[1])
		}
	}

	f.SetActiveSheet(index)

	return f.WriteToBuffer()
}

func CalculateColVal(n int) string {
	var col string
	if n <= 26 {
		col = string(rune(65 + n - 1))
	} else {
		start := 65 + n/26 - 1
		end := 65 + n%26 - 1
		if n%26 == 0 { // 最后一个 Z
			end = 90
			start -= 1
		}
		col = fmt.Sprintf("%s%s", string(rune(start)), string(rune(end)))
	}

	return col
}

func GetColumnKey(index int) string {
	key := ""
	if index/26 <= 0 {
		key = fmt.Sprintf("%c", 'A'+index)
	} else {
		first := fmt.Sprintf("%c", 'A'+index/26-1)
		second := fmt.Sprintf("%c", 'A'+index%26)
		key = fmt.Sprintf("%s%s", first, second)
	}
	return key
}

/*
 * 返回所有sheet页中所有数据
 * @param buffer 	[]byte 		文件内容
 *
 * @return
 * 		map[string][][]string 	内容数组
 * 		error					错误信息
 */
func ReadExcelReturnAllData(buffer []byte) (map[string][][]string, error) {
	file, err := excelize.OpenReader(bytes.NewReader(buffer))
	if err != nil {
		return nil, errors.New("打开文件错误")
	}

	data := make(map[string][][]string, len(file.GetSheetMap()))

	for _, sheetName := range file.GetSheetMap() {

		rows := file.GetRows(sheetName)

		data[sheetName] = make([][]string, 0, len(rows))

		for _, item := range rows {
			tmp := make([]string, 0, len(item))
			for _, v := range item {
				tmp = append(tmp, v)
			}
			data[sheetName] = append(data[sheetName], tmp)
		}
	}

	return data, nil
}

/*
 * 已map形式返回指定sheet页中的数据
 * @param buffer 		[]byte 			文件内容
 * @param fieldMap 		map[string]int	字段名 - 下标
 * @param sheetIndex	int				sheet页下标
 *
 * @return
 * 		[]map[string]string 	内容数组
 * 		error					错误信息
 */
func ReadExcelForMap(buffer []byte, fieldMap map[string]int, sheetIndex int) ([]map[string]string, error) {

	file, err := excelize.OpenReader(bytes.NewReader(buffer))
	if err != nil {
		return nil, errors.New("打开文件错误")
	}

	data := make([]map[string]string, len(file.GetSheetMap()))
	for i, sheetName := range file.GetSheetMap() {
		if i != sheetIndex {
			continue
		}
		rows := file.GetRows(sheetName)
		data = make([]map[string]string, 0, len(rows))

		for j, item := range rows {
			if j == 0 {
				// 跳过表头
				continue
			}
			tmp := make(map[string]string, len(fieldMap))
			for field, index := range fieldMap {
				if index >= len(item) {
					return nil, errors.New("表头溢出")
				}
				tmp[field] = item[index]
			}
			data = append(data, tmp)
		}
	}

	return data, nil
}

/*
 * Excel时间格式存储的不是字符串而是从1900年1月1日到今天的天数
 * 首先尝试格式化字符串，格式化成功直接返回
 * 格式化失败尝试将字符串转为数字
 * 返回1900.1.1加上记录天数的日期
 *
 * PS：-2天有待确认
 */
func ExcelDate(s, format string) *time.Time {
	date, err := time.ParseInLocation(format, s, time.Local)
	if err == nil {
		return &date
	}

	n, err := strconv.ParseFloat(s, 10)
	if err != nil {
		return nil
	}
	excelEpoch := time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
	date = excelEpoch.Add(time.Duration(n * float64(24 * time.Hour)))
	return &date
}

// 解决CSV中文
func EnUnicode(s string) string {
	decodeBytes, _ := simplifiedchinese.GB18030.NewDecoder().Bytes([]byte(s))
	return string(decodeBytes)
}
