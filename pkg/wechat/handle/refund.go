package handle

import (
	kpkg "dgo/pkg"
	kwechat "dgo/pkg/wechat"
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
)

//申请退款
//https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=9_4&index=6

func NewWxRefund() *wxRefund {
	return &wxRefund{}
}

type wxRefund struct {
	reqParam refundParam
}

//请求参数
type refundParam struct {
	transactionId string  //微信的订单号，优先使用
	outTradeNo    string  //商户系统内部的订单号，当没提供transaction_id时需要传这个。
	nonceStr      string  //随机字符串，不长于32位
	totalFee      float64 //订单金额(单位为元)
	refundFee     float64 //退款金额(单位为元)
	outRefundNo   string  //商户退款单号
	notifyUrl     string  //异步接收微信支付退款结果通知的回调地址
}

//结果
type WeChatPayRefundResp struct {
	ReturnCode    string `xml:"return_code"`
	ReturnMsg     string `xml:"return_msg"`
	ResultCode    string `xml:"result_code"`
	ErrCode       string `xml:"err_code"`
	ErrCodeDes    string `xml:"err_code_des"`
	AppId         string `xml:"appid"`
	MchId         string `xml:"mch_id"`
	NonceStr      string `xml:"nonce_str"`
	TransactionId string `xml:"transaction_id"`
	OutTradeNo    string `xml:"out_trade_no"`
	OutRefundNo   string `xml:"out_refund_no"`
	RefundId      string `xml:"refund_id"`
	RefundFee     int64  `xml:"refund_fee"`
	TotalFee      int64  `xml:"total_fee"`
	CashFee       int64  `xml:"cash_fee"`
}

func (t *wxRefund) SetTransactionId(transactionId string) {
	t.reqParam.transactionId = transactionId
}

func (t *wxRefund) SetOutTradeNo(outTradeNo string) {
	t.reqParam.outTradeNo = outTradeNo
}

func (t *wxRefund) SetNonceStr(nonceStr string) {
	t.reqParam.nonceStr = nonceStr
}

func (t *wxRefund) SetTotalFee(totalFee float64) {
	t.reqParam.totalFee = totalFee
}

func (t *wxRefund) SetRefundFee(refundFee float64) {
	t.reqParam.refundFee = refundFee
}

func (t *wxRefund) SetOutRefundNo(outRefundNo string) {
	t.reqParam.outRefundNo = outRefundNo
}

func (t *wxRefund) SetNotifyUrl(notifyUrl string) {
	t.reqParam.notifyUrl = notifyUrl
}

func (t *wxRefund) GetParam(wx *kwechat.WeClient) map[string]interface{} {
	param := make(map[string]interface{})
	param["appid"] = wx.GetAppsAppId()
	param["mch_id"] = wx.GetMchId()
	param["nonce_str"] = t.reqParam.nonceStr
	param["out_trade_no"] = t.reqParam.outTradeNo
	param["transaction_id"] = t.reqParam.transactionId
	param["total_fee"] = t.reqParam.totalFee * 100
	param["refund_fee"] = t.reqParam.refundFee * 100
	param["out_refund_no"] = t.reqParam.outRefundNo
	param["notify_url"] = t.reqParam.notifyUrl
	return param
}

func (t *wxRefund) ApiName() string {
	return "secapi/pay/refund"
}

func (t *wxRefund) QueryPOST() bool {
	return true
}

func (t *wxRefund) GetResult(wx *kwechat.WeClient, body string) (interface{}, error) {
	var resReturn WeChatPayRefundResp

	err := xml.Unmarshal([]byte(body), &resReturn)
	if err != nil {
		return resReturn, err
	}

	if !(strings.ToUpper(resReturn.ReturnCode) == kpkg.WeChatCodeSuccess && strings.ToUpper(resReturn.ResultCode) == kpkg.WeChatCodeSuccess) {
		return resReturn, errors.New(fmt.Sprintf("wechat PayRefund error status, result body: %+v", resReturn))
	}
	return resReturn, nil
}
