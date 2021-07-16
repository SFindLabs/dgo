package alipay

import (
	"crypto"
	"crypto/md5"
	"crypto/rsa"
	"crypto/x509"
	khttp "dgo/framework/tools/http"
	kpkg "dgo/pkg"
	pcrypto "dgo/pkg/crypto"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Jeffail/gabs"
	"io/ioutil"
	"net/url"
	"sort"
	"strings"
	"time"
)

type WxApiObj interface {
	Method() string
	GetParam() map[string]interface{}
	GetExtendParam() map[string]interface{}
	GetResult(*gabs.Container) (interface{}, error)
}

var AliGlobal = make(map[string]*AliClient)
var AliWithdrawGlobal = make(map[string]*AliClient)

type AliClient struct {
	appId  string
	apiUrl string

	appPrivateKey *rsa.PrivateKey // 应用私钥
	appCertSN     string          // 应用公钥证书 SN
	rootCertSN    string          // 支付宝根证书 SN
	aliPublicKey  *rsa.PublicKey  // 支付宝公钥
}

func InitAliPayGlobal(project string, isProduction bool, isPay bool) *AliClient {
	client := &AliClient{}
	if isProduction {
		client.apiUrl = kpkg.AliPayProductionURL
	} else {
		client.apiUrl = kpkg.AliPaySandboxURL
	}
	if isPay {
		AliGlobal[project] = client
	} else {
		AliWithdrawGlobal[project] = client
	}
	return client
}

func (ali *AliClient) SetAppsAppId(appId string) {
	ali.appId = appId
}

func (ali *AliClient) LoadPrivateKey(privateKey, privateKeyPath string) error {
	if privateKey == "" && privateKeyPath != "" {
		content, err := ali.LoadFromFile(privateKeyPath)
		if err != nil {
			return err
		}
		privateKey = content
	}
	priKey, err := pcrypto.ParsePKCS1PrivateKey(pcrypto.FormatPKCS1PrivateKey(privateKey))
	if err != nil {
		priKey, err = pcrypto.ParsePKCS8PrivateKey(pcrypto.FormatPKCS8PrivateKey(privateKey))
		if err != nil {
			return err
		}
	}
	ali.appPrivateKey = priKey
	return nil
}

func (ali *AliClient) LoadPublicKeyPath(publicKey, publicKeyPath string) error {
	if publicKey == "" && publicKeyPath != "" {
		content, err := ali.LoadFromFile(publicKeyPath)
		if err != nil {
			return err
		}
		publicKey = content
	}
	pubKey, err := ali.GetAppPublicCert(publicKey)
	if err != nil {
		return err
	}
	ali.appCertSN = pubKey
	return nil
}

func (ali *AliClient) LoadAliRootKeyPath(aliRootKey, aliRootKeyPath string) error {
	if aliRootKey == "" && aliRootKeyPath != "" {
		content, err := ali.LoadFromFile(aliRootKeyPath)
		if err != nil {
			return err
		}
		aliRootKey = content
	}
	rootKey, err := ali.GetAliRootCert(aliRootKey)
	if err != nil {
		return err
	}
	ali.rootCertSN = rootKey
	return nil
}

//添加支付宝公钥
func (ali *AliClient) AddAliPublicKey(publicKey, publicKeyPath string) error {
	if publicKey == "" && publicKeyPath != "" {
		content, err := ali.LoadFromFile(publicKeyPath)
		if err != nil {
			return err
		}
		publicKey = content
	}
	pubKey, err := ali.LoadAliPayPublicKey(publicKey)
	if err != nil {
		return err
	}
	ali.aliPublicKey = pubKey
	return nil
}

//==================================================方法处理====================================================

//从文件中加载证书
func (ali *AliClient) LoadFromFile(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

//获取应用公钥证书SN
func (ali *AliClient) GetAppPublicCert(s string) (string, error) {
	cert, err := pcrypto.ParseCertificate([]byte(s))
	if err != nil {
		return "", err
	}
	return ali.GetCertSN(cert), nil
}

//获取支付宝根证书SN
func (ali *AliClient) GetAliRootCert(s string) (string, error) {
	certStrList := strings.Split(s, "-----END CERTIFICATE-----")
	certSNList := make([]string, 0, len(certStrList))
	for _, certStr := range certStrList {
		certStr = certStr + "-----END CERTIFICATE-----"
		cert, _ := pcrypto.ParseCertificate([]byte(certStr))
		if cert != nil && (cert.SignatureAlgorithm == x509.SHA256WithRSA || cert.SignatureAlgorithm == x509.SHA1WithRSA) {
			certSNList = append(certSNList, ali.GetCertSN(cert))
		}
	}
	return strings.Join(certSNList, "_"), nil
}

// 加载支付宝公钥
func (ali *AliClient) LoadAliPayPublicKey(aliPublicKey string) (*rsa.PublicKey, error) {
	if len(aliPublicKey) < 0 {
		return nil, errors.New("alipay public key not found")
	}
	pub, err := pcrypto.ParseCertificate([]byte(aliPublicKey))
	if err != nil {
		return nil, err
	}
	key, ok := pub.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, err
	}
	return key, nil
}

//生成SN值
func (ali *AliClient) GetCertSN(cert *x509.Certificate) string {
	value := md5.Sum([]byte(cert.Issuer.String() + cert.SerialNumber.String()))
	return hex.EncodeToString(value[:])
}

//==================================================签名处理====================================================

//生成签名
func (ali *AliClient) GetSign(queryParam map[string]interface{}, method string, extendParam map[string]interface{}) (map[string]interface{}, error) {
	param := make(map[string]interface{})
	if len(queryParam) > 0 {
		bytes, err := json.Marshal(queryParam)
		if err != nil {
			return nil, err
		}
		param["biz_content"] = string(bytes)
	}

	param["app_id"] = ali.appId
	param["method"] = method
	param["format"] = kpkg.AliPayFormat
	param["charset"] = kpkg.AliPayCharset
	param["sign_type"] = kpkg.AliPaySignType
	param["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	param["version"] = kpkg.AliPayVersion
	param["app_cert_sn"] = ali.appCertSN
	param["alipay_root_cert_sn"] = ali.rootCertSN

	if extendParam != nil {
		for key, val := range extendParam {
			param[key] = val
		}
	}

	keyList := make([]string, 0)
	valList := make(map[string]string)
	for key, val := range param {
		tmpVal := fmt.Sprint(val)
		if key == "sign" || tmpVal == "" {
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
	rawStr = strings.TrimRight(rawStr, "&")

	sign, err := pcrypto.RSASignWithKey([]byte(rawStr), ali.appPrivateKey, crypto.SHA256)
	if err != nil {
		return nil, err
	}
	param["sign"] = base64.StdEncoding.EncodeToString(sign)
	return param, nil
}

//验签
func (ali *AliClient) VerifySign(data map[string]interface{}) (bool, error) {
	keys := make([]string, 0)
	for key := range data {
		if key == "sign" || key == "sign_type" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	pList := make([]string, 0)
	for _, key := range keys {
		tmpData, _ := url.QueryUnescape(fmt.Sprint(data[key]))
		pList = append(pList, key+"="+tmpData)
	}
	rawStr := strings.Join(pList, "&")
	signBytes, err := base64.StdEncoding.DecodeString(fmt.Sprint(data["sign"]))
	if err != nil {
		return false, err
	}

	if err = pcrypto.RSAVerifyWithKey([]byte(rawStr), signBytes, ali.aliPublicKey, crypto.SHA256); err != nil {
		return false, err
	}
	return true, nil
}

//==================================================请求处理====================================================

func (ali *AliClient) Request(api WxApiObj) (interface{}, error) {
	queryParam, err := ali.GetSign(api.GetParam(), api.Method(), api.GetExtendParam())
	if err != nil {
		return nil, err
	}

	urlParamArr := make([]string, 0)
	for k, v := range queryParam {
		urlParamArr = append(urlParamArr, k+"="+url.QueryEscape(fmt.Sprint(v)))
	}

	queryUrl := fmt.Sprintf("%s?%s", ali.apiUrl, strings.Join(urlParamArr, "&"))
	jsonParsed, _, err := khttp.UrlGetGetJsonObj(queryUrl, 10)
	if err != nil {
		return nil, err
	}

	return api.GetResult(jsonParsed)
}

//==================================================不用请求处理的方法====================================================

/**构造app授权信息
targetId 商户标识该次用户授权请求的 ID，该值在商户端应保持唯一
*/
//https://opendocs.alipay.com/open/218/105325

func (ali *AliClient) AppAuth(targetId string) (appAuthInfo string, err error) {
	param := make(map[string]string)
	param["apiname"] = "com.alipay.account.auth"
	param["app_id"] = ali.appId
	param["app_name"] = "mc"
	param["auth_type"] = "AUTHACCOUNT"
	param["biz_type"] = "openservice"
	param["pid"] = kpkg.AliPayPid
	param["product_id"] = "APP_FAST_LOGIN"
	param["scope"] = "kuaijie"
	param["method"] = "alipay.open.auth.sdk.code.get"
	param["sign_type"] = "RSA2"
	param["target_id"] = targetId

	keyList := make([]string, 0)
	valList := make(map[string]string)
	for key, val := range param {
		keyList = append(keyList, key)
		valList[key] = val
	}
	sort.Strings(keyList)
	rawStr := ""
	for _, key := range keyList {
		rawStr = fmt.Sprintf("%s%s=%s&", rawStr, key, valList[key])
	}
	rawStr = strings.TrimRight(rawStr, "&")

	sign, err := pcrypto.RSASignWithKey([]byte(rawStr), ali.appPrivateKey, crypto.SHA256)
	if err != nil {
		return "", err
	}
	param["sign"] = base64.StdEncoding.EncodeToString(sign)

	urlParamArr := make([]string, 0)
	for k, v := range param {
		urlParamArr = append(urlParamArr, k+"="+url.QueryEscape(fmt.Sprint(v)))
	}

	return strings.Join(urlParamArr, "&"), nil
}

//app支付接口2.0
//https://opendocs.alipay.com/apis/api_1/alipay.trade.app.pay

/**
 * subject         商品标题/交易标题/订单标题/订单关键字等
 * outTradeNo      商户订单号
 * publicParam     透传参数(为空时不传)
 * notifyUrl       支付宝服务器主动通知商户服务器里指定的页面http/https路径
 * totalAmount     订单总金额，单位为元，精确到小数点后两位
 * isVirtual       是否是虚拟类商品
 */
func (ali *AliClient) AppPay(subject, outTradeNo, publicParam, notifyUrl string, totalAmount float64, isVirtual bool) (orderInfo string, err error) {

	method := "alipay.trade.app.pay"
	param := make(map[string]interface{})

	param["subject"] = subject
	param["out_trade_no"] = outTradeNo
	param["total_amount"] = totalAmount
	param["product_code"] = "QUICK_MSECURITY_PAY"

	if isVirtual {
		param["goods_type"] = "0"
	} else {
		param["goods_type"] = "1"
	}

	if publicParam != "" {
		param["passback_params"] = publicParam //透传参数
	}

	extendParam := make(map[string]interface{})
	extendParam["notify_url"] = notifyUrl

	queryParam, err := ali.GetSign(param, method, extendParam)
	if err != nil {
		return
	}

	urlParamArr := make([]string, 0)
	for k, v := range queryParam {
		urlParamArr = append(urlParamArr, k+"="+url.QueryEscape(fmt.Sprint(v)))
	}

	return strings.Join(urlParamArr, "&"), nil
}

// 解析app支付接口结果
//  https://opendocs.alipay.com/open/01dcc0
//  https://opendocs.alipay.com/apis/api_1/alipay.trade.app.pay

type AliAppPayResp struct {
	Code            string `json:"code"`
	Msg             string `json:"msg"`
	SubCode         string `json:"sub_code"`
	SubMsg          string `json:"sub_msg"`
	OutTradeNo      string `json:"out_trade_no"`
	TradeNo         string `json:"trade_no"`
	TotalAmount     string `json:"total_amount"`
	SellerId        string `json:"seller_id"`
	MerchantOrderNo string `json:"merchant_order_no"`
}

func (ali *AliClient) GetAppPayResp(body string) (AliAppPayResp, error) {
	var result AliAppPayResp
	jsonParsed, err := gabs.ParseJSON([]byte(body))
	if err != nil {
		return result, err
	}

	method := "alipay.trade.app.pay"
	field := fmt.Sprintf("%s_response", strings.Replace(method, ".", "_", -1))
	resStr := jsonParsed.Search(field).String()
	err = json.Unmarshal([]byte(resStr), &result)
	if err != nil {
		return result, err
	}
	if kpkg.AliPayCodeSuccess != result.Code {
		return result, errors.New(fmt.Sprintf("alipay AppPay error status, result body: %+v", result))
	}
	return result, nil
}
