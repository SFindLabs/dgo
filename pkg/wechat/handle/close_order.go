package handle

import (
	kpkg "dgo/pkg"
	kwechat "dgo/pkg/wechat"
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
)

//关闭订单
//https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=9_3&index=5

func NewWxCloseOrder() *wxCloseOrder {
	return &wxCloseOrder{}
}

type wxCloseOrder struct {
	reqParam closeOrderParam
}

//请求参数
type closeOrderParam struct {
	outTradeNo string //商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一。
	nonceStr   string //随机字符串，不长于32位
}

//结果
type WeChatCloseOrderResp struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	AppId      string `xml:"appid"`
	MchId      string `xml:"mch_id"`
	NonceStr   string `xml:"nonce_str"`
	ResultCode string `xml:"result_code"`
	ResultMsg  string `xml:"result_msg"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`
}

func (t *wxCloseOrder) SetOutTradeNo(outTradeNo string) {
	t.reqParam.outTradeNo = outTradeNo
}

func (t *wxCloseOrder) SetNonceStr(nonceStr string) {
	t.reqParam.nonceStr = nonceStr
}

func (t *wxCloseOrder) GetParam(wx *kwechat.WeClient) map[string]interface{} {
	param := make(map[string]interface{})
	param["appid"] = wx.GetAppsAppId()
	param["mch_id"] = wx.GetMchId()
	param["nonce_str"] = t.reqParam.nonceStr
	param["out_trade_no"] = t.reqParam.outTradeNo
	return param
}

func (t *wxCloseOrder) ApiName() string {
	return "pay/closeorder"
}

func (t *wxCloseOrder) QueryPOST() bool {
	return true
}

func (t *wxCloseOrder) GetResult(wx *kwechat.WeClient, body string) (interface{}, error) {
	var resReturn WeChatCloseOrderResp

	err := xml.Unmarshal([]byte(body), &resReturn)
	if err != nil {
		return resReturn, err
	}

	if !(strings.ToUpper(resReturn.ReturnCode) == kpkg.WeChatCodeSuccess && strings.ToUpper(resReturn.ResultCode) == kpkg.WeChatCodeSuccess) {
		return resReturn, errors.New(fmt.Sprintf("wechat CloseOrder error status, result body: %+v", resReturn))
	}
	return resReturn, nil
}
