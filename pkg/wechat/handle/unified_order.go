package handle

import (
	kpkg "dgo/pkg"
	kwechat "dgo/pkg/wechat"
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
	"time"
)

//统一下单
//https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=9_1
//https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=9_12&index=2

func NewWxUnifiedOrder() *wxUnifiedOrder {
	return &wxUnifiedOrder{}
}

type wxUnifiedOrder struct {
	reqParam unifiedOrderParam
}

//请求参数
type unifiedOrderParam struct {
	body       string  //商品描述： APP——需传入应用市场上的APP名字-实际商品名称，天天爱消除-游戏充值。
	outTradeNo string  //商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|*且在同一个商户号下唯一
	ip         string  //终端IP
	notifyUrl  string  //微信支付异步通知回调地址，通知url必须为直接可访问的url，不能携带参数。
	attach     string  //透传参数
	totalFee   float64 //订单总金额，单位为元
	nonceStr   string  //随机字符串，不长于32位
}

//结果
type WeChatUnifiedOrderResp struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	AppId      string `xml:"appid"`
	MchId      string `xml:"mch_id"`
	NonceStr   string `xml:"nonce_str"`
	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`
	TradeType  string `xml:"trade_type"`
	PrepayId   string `xml:"prepay_id"`
}

func (t *wxUnifiedOrder) SetBody(body string) {
	t.reqParam.body = body
}

func (t *wxUnifiedOrder) SetOutTradeNo(outTradeNo string) {
	t.reqParam.outTradeNo = outTradeNo
}

func (t *wxUnifiedOrder) SetIp(ip string) {
	t.reqParam.ip = ip
}

func (t *wxUnifiedOrder) SetNotifyUrl(notifyUrl string) {
	t.reqParam.notifyUrl = notifyUrl
}

func (t *wxUnifiedOrder) SetAttach(attach string) {
	t.reqParam.attach = attach
}

func (t *wxUnifiedOrder) SetNonceStr(nonceStr string) {
	t.reqParam.nonceStr = nonceStr
}

func (t *wxUnifiedOrder) SetTotalFee(totalFee float64) {
	t.reqParam.totalFee = totalFee
}

func (t *wxUnifiedOrder) GetParam(wx *kwechat.WeClient) map[string]interface{} {
	param := make(map[string]interface{})
	param["appid"] = wx.GetAppsAppId()
	param["mch_id"] = wx.GetMchId()

	param["nonce_str"] = t.reqParam.nonceStr
	param["body"] = t.reqParam.body
	param["out_trade_no"] = t.reqParam.outTradeNo
	param["spbill_create_ip"] = t.reqParam.ip
	param["total_fee"] = t.reqParam.totalFee * 100
	param["notify_url"] = t.reqParam.notifyUrl
	if t.reqParam.attach != "" {
		param["attach"] = t.reqParam.attach //透传参数
	}
	param["trade_type"] = "APP"
	return param
}

func (t *wxUnifiedOrder) ApiName() string {
	return "pay/unifiedorder"
}

func (t *wxUnifiedOrder) QueryPOST() bool {
	return true
}

//返回map信息
func (t *wxUnifiedOrder) GetResult(wx *kwechat.WeClient, body string) (interface{}, error) {
	var resReturn WeChatUnifiedOrderResp
	err := xml.Unmarshal([]byte(body), &resReturn)
	if err != nil {
		return resReturn, err
	}

	result := make(map[string]interface{})
	if !(strings.ToUpper(resReturn.ReturnCode) == kpkg.WeChatCodeSuccess && strings.ToUpper(resReturn.ResultCode) == kpkg.WeChatCodeSuccess) {
		return result, errors.New(fmt.Sprintf("wechat UnifiedOrder error status, result body: %+v", resReturn))
	}

	//生成调起支付接口的参数信息
	result["appid"] = wx.GetAppsAppId()
	result["partnerid"] = wx.GetMchId()
	result["noncestr"] = t.reqParam.nonceStr

	result["prepayid"] = resReturn.PrepayId
	result["package"] = "Sign=WXPay"
	result["timestamp"] = fmt.Sprintf("%v", time.Now().Unix())

	orderInfoSign, err := wx.GetSign(result)
	if err != nil {
		return result, err
	}
	result["sign"] = orderInfoSign

	return result, nil
}
