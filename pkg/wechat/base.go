package wechat

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	kpkg "dgo/pkg"
	pcrypto "dgo/pkg/crypto"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"io/ioutil"
	"sort"
	"strings"
	"time"
)

type WxApiObj interface {
	ApiName() string
	QueryPOST() bool
	GetParam(*WeClient) map[string]interface{}
	GetResult(*WeClient, string) (interface{}, error)
}

var WxPayGlobal = make(map[string]*WeClient)
var WxWithdrawGlobal = make(map[string]*WeClient)

type WeClient struct {
	appsAppId                string //移动应用开发者ID
	appsAppSecret            string //移动应用密钥
	officialAccountAppId     string //公众号开发者ID
	officialAccountAppSecret string //公众号开发者密钥
	mchId                    string //微信支付分配的商户号
	apiKey                   string //商户Key 、支付Key、 API密钥
	mchApiUrl                string
	conf                     *tls.Config //API证书配置
}

//-------------------------------------setter----------------------------------------------

func (wx *WeClient) SetOfficialAccountAppId(officialAccountAppId string) {
	wx.officialAccountAppId = officialAccountAppId
}

func (wx *WeClient) SetOfficialAccountAppSecret(officialAccountAppSecret string) {
	wx.officialAccountAppSecret = officialAccountAppSecret
}

func (wx *WeClient) SetAppsAppId(appsAppId string) {
	wx.appsAppId = appsAppId
}

func (wx *WeClient) SetAppsAppSecret(appsAppSecret string) {
	wx.appsAppSecret = appsAppSecret
}

func (wx *WeClient) SetMchId(mchId string) {
	wx.mchId = mchId
}

func (wx *WeClient) SetApiKey(apiKey string) {
	wx.apiKey = apiKey
}

//-------------------------------------getter----------------------------------------------

func (wx *WeClient) GetOfficialAccountAppId() string {
	return wx.officialAccountAppId
}

func (wx *WeClient) GetAppsAppId() string {
	return wx.appsAppId
}

func (wx *WeClient) GetMchId() string {
	return wx.mchId
}

func InitWeChatGlobal(project string, isProduction bool, isPay bool) *WeClient {
	client := &WeClient{}

	if isProduction {
		client.mchApiUrl = kpkg.WeChatMchProductionURL
	} else {
		client.mchApiUrl = kpkg.WeChatMchSandboxURL
	}

	if isPay {
		WxPayGlobal[project] = client
	} else {
		WxWithdrawGlobal[project] = client
	}
	return client
}

func (wx *WeClient) LoadCert(content, path string) error {
	var contentByte []byte
	if len(content) == 0 && path != "" {
		tmpByte, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		contentByte = tmpByte
	}

	cert, err := pcrypto.Pkcs12ToPem(contentByte, wx.mchId)
	if err != nil {
		return err
	}
	wx.conf = &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	return nil
}

//运行请求
func (wx *WeClient) RequestXML(api WxApiObj) (interface{}, error) {
	commonParam := api.GetParam(wx)
	queryUrl := fmt.Sprintf("%s/%s", wx.mchApiUrl, api.ApiName())
	sign, err := wx.GetSign(commonParam)
	if err != nil {
		return nil, err
	}
	commonParam["sign"] = sign
	queryStr := wx.URLValueToXML(commonParam)
	body, err := wx.UrlGetPostObj(api.QueryPOST(), queryUrl, "xml", queryStr, 10, wx.conf)
	if err != nil {
		return nil, err
	}
	return api.GetResult(wx, body)
}

//=====================================签名==============================================

func (wx *WeClient) GetSign(queryParam map[string]interface{}) (string, error) {
	if len(queryParam) == 0 {
		return "", errors.New("query param is empty")
	}
	keyList := make([]string, 0)
	valList := make(map[string]string)
	for key, val := range queryParam {
		tmpVal := fmt.Sprint(val)
		if tmpVal == "" || key == "sign" {
			continue
		}
		keyList = append(keyList, key)
		valList[key] = tmpVal
	}
	sort.Strings(keyList)
	rawStr := ""
	for _, key := range keyList {
		rawStr = fmt.Sprintf("%s%s=%s&", rawStr, key, valList[key])
	}
	rawStr = fmt.Sprintf("%skey=%s", rawStr, wx.apiKey)
	h := md5.New()
	h.Write([]byte(rawStr))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil))), nil
}

func (wx *WeClient) URLValueToXML(param map[string]interface{}) string {
	xmlBuffer := &bytes.Buffer{}
	xmlBuffer.WriteString("<xml>")

	for key, value := range param {
		xmlBuffer.WriteString("<" + key + ">" + fmt.Sprint(value) + "</" + key + ">")
	}
	xmlBuffer.WriteString("</xml>")
	return xmlBuffer.String()
}

//=====================================http请求==============================================
//typeStr supports
//    "text/html" uses "html"
//    "application/json" uses "json"
//    "application/xml" uses "xml"
//    "text/plain" uses "text"
//    "application/x-www-form-urlencoded" uses "urlencoded", "form" or "form-data"
func (wx *WeClient) UrlGetPostObj(isPost bool, url string, typeStr string, param interface{}, timeout int64, config *tls.Config) (string, error) {
	var body string
	var errs []error
	request := gorequest.New().Timeout(time.Duration(timeout) * time.Second)
	if config != nil {
		request = request.TLSClientConfig(config)
	}
	if isPost {
		_, body, errs = request.Post(url).Type(typeStr).Send(param).End()
	} else {
		_, body, errs = request.Get(url).End()
	}

	if errs != nil {
		return body, errs[0]
	}

	return body, nil
}
