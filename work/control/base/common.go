package base

import (
	kinit "dgo/work/base/initialize"
	kcode "dgo/work/code"
	"encoding/json"
	"github.com/Jeffail/gabs"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type DataIStruct struct {
	Status         int         `json:"status"`
	Info           string      `json:"message"`
	Data           interface{} `json:"data"`
	InnerException string      `json:"innerException"`
}

func init() {

}

func CheckGetPostParam(c *gin.Context, key string) (string, bool) {
	if value, ok := c.GetQuery(key); ok {
		return value, ok
	}
	if value, ok := c.GetPostForm(key); ok {
		return value, ok
	}
	return "", false
}

func GetParam(c *gin.Context, key string) string {
	v := c.Query(key)
	if v == "" {
		v = c.PostForm(key)
	}
	if v == "" {
		v = c.Param(key)
	}
	if v == "" {
		v = c.GetHeader(key)
	}

	return v
}
func GetParamInt64DefaultNegative(c *gin.Context, key string) int64 {
	v := GetParam(c, key)
	if v == "" {
		return -1
	}
	vv, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return -1
	}
	return vv
}

func ReturnDataI(c *gin.Context, status int, v interface{}, callbackName string) []byte {
	object := DataIStruct{
		Status: status,
		Info:   kcode.GetCodeMsg(status),
		Data:   v,
	}
	return ReturnData(c, object, callbackName)
}
func ReturnData(c *gin.Context, v interface{}, callbackName string) []byte {
	jsonStr, err := json.Marshal(v)
	if err != nil {
		kinit.LogError.Println(err)
	}
	if callbackName == "" {
		c.Data(http.StatusOK, "text/plain", jsonStr)
	} else {
		res := []byte(callbackName)
		res = append(res, []byte("(")...)
		res = append(res, jsonStr...)
		res = append(res, []byte(");")...)
		c.Data(http.StatusOK, "application/json; charset=utf-8", res)
	}
	return jsonStr
}
func ReturnDataStr(c *gin.Context, jsonStr []byte, callbackName string) []byte {
	if callbackName == "" {
		c.Data(http.StatusOK, "text/plain", jsonStr)
	} else {
		res := []byte(callbackName)
		res = append(res, []byte("(")...)
		res = append(res, jsonStr...)
		res = append(res, []byte(");")...)
		c.Data(http.StatusOK, "application/json; charset=utf-8", res)
	}
	return jsonStr
}

func SendErrorJsonStr(c *gin.Context, code int, callbackName string) {
	jsonStr := GetErrorJsonStr(code)
	if callbackName == "" {
		c.Data(http.StatusOK, "text/plain", jsonStr)
	} else {
		res := []byte(callbackName)
		res = append(res, []byte("(")...)
		res = append(res, jsonStr...)
		res = append(res, []byte(");")...)
		c.Data(http.StatusOK, "application/json; charset=utf-8", res)
	}
}

func SendErrorJsonStrx(c *gin.Context, code int, innerException string) {
	//msg := kcode.GetCodeMsg(code)
	chnMsg := kcode.GetCodeChnMsg(code)

	jsonObj := gabs.New()
	_, _ = jsonObj.Set(code, "status")
	_, _ = jsonObj.Set(innerException, "innerException")
	_, _ = jsonObj.Set(chnMsg, "message")
	_, _ = jsonObj.Set(struct{}{}, "data")
	jsonStr := jsonObj.Bytes()
	c.Data(http.StatusOK, "text/plain", jsonStr)
}

func GetErrorJsonStr(code int) []byte {
	//msg := kcode.GetCodeMsg(code)
	chnMsg := kcode.GetCodeChnMsg(code)

	jsonObj := gabs.New()
	_, _ = jsonObj.Set(code, "status")
	_, _ = jsonObj.Set("", "innerException")
	_, _ = jsonObj.Set(chnMsg, "message")
	_, _ = jsonObj.Set(struct{}{}, "data")
	return jsonObj.Bytes()
}

func GetErrorOriginJsonStr(code int, data interface{}) []byte {
	chnMsg := kcode.GetCodeChnMsg(code)
	jsonObj := gabs.New()
	_, _ = jsonObj.Set(code, "status")
	_, _ = jsonObj.Set("", "innerException")
	_, _ = jsonObj.Set(chnMsg, "message")
	_, _ = jsonObj.Set(data, "data")
	return jsonObj.Bytes()
}

func SendErrorOriginJsonStr(c *gin.Context, code int, data interface{}, callbackName string) {
	jsonStr := GetErrorOriginJsonStr(code, data)
	if callbackName == "" {
		c.Data(http.StatusOK, "text/plain", jsonStr)
	} else {
		res := []byte(callbackName)
		res = append(res, []byte("(")...)
		res = append(res, jsonStr...)
		res = append(res, []byte(");")...)
		c.Data(http.StatusOK, "application/json; charset=utf-8", res)
	}
}

func GetErrorParamsJsonStr(code int, errors error) []byte {
	jsonObj := gabs.New()
	_, _ = jsonObj.Set(code, "status")
	_, _ = jsonObj.Set("", "innerException")
	_, _ = jsonObj.Set(errors.Error(), "message")
	_, _ = jsonObj.Set(struct{}{}, "data")
	return jsonObj.Bytes()
}

func SendErrorParamsJsonStr(c *gin.Context, code int, errors error, callbackName string) {
	jsonStr := GetErrorParamsJsonStr(code, errors)
	if callbackName == "" {
		c.Data(http.StatusOK, "text/plain", jsonStr)
	} else {
		res := []byte(callbackName)
		res = append(res, []byte("(")...)
		res = append(res, jsonStr...)
		res = append(res, []byte(");")...)
		c.Data(http.StatusOK, "application/json; charset=utf-8", res)
	}
}

func RenderTokenHtml(c *gin.Context, pageName string, data map[string]interface{}) {
	data["token"] = c.GetString("token")
	c.HTML(http.StatusOK, pageName, data)
}
