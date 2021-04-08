package utils

import (
	kcommon "dgo/framework/tools/common"
	khttp "dgo/framework/tools/http"
	kruntime "dgo/framework/tools/runtime"
	kcode "dgo/work/code"
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/gin-gonic/gin"
	"html"
	"html/template"
	"io"
	"io/ioutil"
	"math"
	"math/big"
	"mime/multipart"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func SaveFile(c *gin.Context, dir string, fileHeader *multipart.FileHeader) (error, int, string, string) {
	is, err := kruntime.PathExists(dir)
	if err != nil {
		return err, kcode.FILE_UPLOAD_FAIL, "", ""
	}
	if !is {
		_ = os.MkdirAll(dir, 755)
	}
	randNum := kcommon.GetRandomString(6, 0)
	fileExt := path.Ext(fileHeader.Filename)
	times := strconv.FormatInt(time.Now().UnixNano(), 10)
	fileName := fmt.Sprintf("%s%s%s", times, randNum, fileExt)
	if err := c.SaveUploadedFile(fileHeader, fmt.Sprint(dir, fileName)); err != nil {
		return err, kcode.FILE_UPLOAD_FAIL, "", ""
	}
	return nil, kcode.SUCCESS_STATUS, fileName, fileExt
}

func SaveBurstFile(c *gin.Context, dir string, index int64, fileExt string, fileHeader *multipart.FileHeader) (error, int, string, string) {
	is, err := kruntime.PathExists(dir)
	if err != nil {
		return err, kcode.FILE_UPLOAD_FAIL, "", ""
	}
	if !is {
		_ = os.MkdirAll(dir, 755)
	}
	savePath := fmt.Sprintf("%s/%d%s", dir, index, fileExt)
	if err := c.SaveUploadedFile(fileHeader, savePath); err != nil {
		return err, kcode.FILE_UPLOAD_FAIL, "", ""
	}
	return nil, kcode.SUCCESS_STATUS, savePath, fileExt
}

func SaveUrlFile(url, dir, fileName, fileExt string) (string, error) {
	_, body, err := khttp.UrlGetGetJsonObj(url, 10)
	is, err := kruntime.PathExists(dir)
	if err != nil {
		return "", err
	}
	if !is {
		_ = os.MkdirAll(dir, 755)
	}
	if fileExt == "" {
		fileExt = "jpg"
	}
	saveFile := fmt.Sprintf("%s%s.%s", dir, fileName, fileExt)
	out, err := os.Create(saveFile)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = out.Close()
	}()
	_, err = io.Copy(out, bytes.NewReader([]byte(body)))
	if err != nil {
		return "", err
	}
	return saveFile, nil
}

func AssertFloat64(i interface{}) float64 {
	if v, ok := i.(float64); ok {
		return v
	}
	return 0
}
func AssertString(i interface{}) string {
	if v, ok := i.(string); ok {
		return v
	}
	return ""
}
func AssertInt64(i interface{}) int64 {
	if v, ok := i.(int64); ok {
		return int64(v)
	}
	return 0
}
func AssertBool(i interface{}) bool {
	if v, ok := i.(bool); ok {
		return v
	}
	return false
}

func FormatTime(timeStr string, format ...string) string {
	formatTime, err := time.ParseInLocation("2006-01-02T15:04:05+08:00", timeStr, time.Local)
	if err != nil {
		return timeStr
	}
	formatStr := "2006-01-02 15:04:05"
	if len(format) > 0 {
		formatStr = format[0]
	}
	return time.Unix(formatTime.Unix(), 0).Format(formatStr)
}

func FormatUnixToTimeStr(timeValue int64, isUnix, isMillisecond, isUnixNano bool) string {
	var value int64
	if isUnixNano {
		value = timeValue / 1e9
	}
	if isMillisecond {
		value = timeValue / 1e3
	}
	if isUnix {
		value = timeValue
	}
	return time.Unix(value, 0).Format("2006-01-02 15:04:05")
}

func FormatTimeStrToSecond(timeStr string, format ...string) int64 {
	formatStr := "2006-01-02T15:04:05+08:00"
	if len(format) > 0 {
		formatStr = format[0]
	}
	formatTime, err := time.ParseInLocation(formatStr, timeStr, time.Local)
	if err != nil {
		return 0
	}
	return formatTime.Unix()
}

//获取当前目录下的一层的文件数(用于分片上传时统计文件缓存目录下文件数)
func GetOneDirFileNum(path string) (int64, []string, error) {
	var count int64
	fileName := make([]string, 0)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return count, fileName, err
	}
	for _, f := range files {
		fileName = append(fileName, f.Name())
		count++
	}
	return count, fileName, nil
}

//分片合并文件
func MergeFile(avatarDir, fileName string, path string, length int, suffix string) error {
	is, err := kruntime.PathExists(avatarDir)
	if err != nil {
		return err
	}
	if !is {
		_ = os.MkdirAll(avatarDir, 755)
	}
	outPutFile := fmt.Sprintf("%s%s", avatarDir, fileName)
	tmpFile, err := os.OpenFile(outPutFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	defer func() {
		_ = tmpFile.Close()
	}()
	if err != nil {
		return err
	}
	for i := 1; i <= length; i++ {
		f, err := os.OpenFile(fmt.Sprintf("%s/%d%s", path, i, suffix), os.O_RDONLY, os.ModePerm)
		if err != nil {
			return err
		}
		blob, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		if _, err := tmpFile.Write(blob); err != nil {
			return err
		}
		_ = f.Close()
	}
	return nil
}

//html显示空格用于菜单树形格式
func HtmlSpaceFunc(num int64) template.HTML {
	return template.HTML(strings.Repeat("&nbsp;&nbsp;", int(num)))
}

func Bytes2BitsString(data []byte) string {
	dst := ""
	for _, v := range data {
		for i := 0; i < 8; i++ {
			move := uint(7 - i)
			dst += strconv.Itoa(int((v >> move) & 1))
		}
	}
	return dst
}

// 生成区间[-m, n]的安全随机数
func RangeRand(min, max int64) int64 {
	if min > max {
		return 0
	}

	if min < 0 {
		f64Min := math.Abs(float64(min))
		i64Min := int64(f64Min)
		result, _ := rand.Int(rand.Reader, big.NewInt(max+1+i64Min))

		return result.Int64() - i64Min
	} else {
		result, _ := rand.Int(rand.Reader, big.NewInt(max-min+1))
		return min + result.Int64()
	}
}

//转义特殊字符
func XssConvertFunc(srcStr string, isEscape bool) string {
	result := ""
	if isEscape {
		result = html.EscapeString(srcStr)
	} else {
		result = html.UnescapeString(srcStr)
	}
	return result
}

//截取字符串
func Substr(str string, substr string, count int) string {
	if str == "" {
		return ""
	}
	s, c := 0, 0
	ss := str
	for {
		i := UnicodeIndex(ss, substr)
		s += i + 1
		if i > -1 {
			ss = string([]rune(ss)[i+1:])
			c++
			if c > count {
				break
			}
		} else {
			break
		}
	}
	s -= 1
	return string([]rune(str)[:s])
}

//获取中文字符串的子串字符位置
func UnicodeIndex(str, substr string) int {
	result := strings.Index(str, substr)
	if result >= 0 {
		prefix := []byte(str)[0:result]
		rs := []rune(string(prefix))
		result = len(rs)
	}
	return result
}
