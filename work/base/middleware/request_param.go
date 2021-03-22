package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

//打印请求数据
func TestPrintRequest(c *gin.Context) string {
	data, _ := c.GetRawData()
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	c.Next()
	return string(data)
}
