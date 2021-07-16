package handle

import (
	kpkg "dgo/pkg"
	kwechat "dgo/pkg/wechat"
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
)

//查询订单
//https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=9_2&index=4

func NewWxRefundQuery() *wxRefundQuery {
	return &wxRefundQuery{}
}

type wxRefundQuery struct {
	reqParam refundQueryParam
}

//请求参数 四选一(微信订单号查询的优先级是： refund_id > out_refund_no > transaction_id > out_trade_no)
type refundQueryParam struct {
	refundId      string //微信退款单号
	outRefundNo   string //商户退款单号
	transactionId string //微信的订单号
	outTradeNo    string //商户系统内部的订单号
	nonceStr      string //随机字符串，不长于32位
}

//结果
type WeChatRefundQueryResp struct {
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
	TotalFee      int64  `xml:"total_fee"`
	CashFee       int64  `xml:"cash_fee"`
	RefundCount   int64  `xml:"refund_count"`
	OutRefundNo0  string `xml:"out_refund_no_0"`
	RefundId0     string `xml:"refund_id_0"`
	RefundFee0    int64  `xml:"refund_fee_0"`
	RefundStatus0 string `xml:"refund_status_0"`
}

func (t *wxRefundQuery) SetRefundId(refundId string) {
	t.reqParam.refundId = refundId
}

func (t *wxRefundQuery) SetOutRefundNo(outRefundNo string) {
	t.reqParam.outRefundNo = outRefundNo
}

func (t *wxRefundQuery) SetTransactionId(transactionId string) {
	t.reqParam.transactionId = transactionId
}

func (t *wxRefundQuery) SetOutTradeNo(outTradeNo string) {
	t.reqParam.outTradeNo = outTradeNo
}

func (t *wxRefundQuery) SetNonceStr(nonceStr string) {
	t.reqParam.nonceStr = nonceStr
}

func (t *wxRefundQuery) GetParam(wx *kwechat.WeClient) map[string]interface{} {
	param := make(map[string]interface{})
	param["appid"] = wx.GetAppsAppId()
	param["mch_id"] = wx.GetMchId()
	param["nonce_str"] = t.reqParam.nonceStr
	param["refund_id"] = t.reqParam.refundId
	param["out_refund_no"] = t.reqParam.outRefundNo
	param["transaction_id"] = t.reqParam.transactionId
	param["out_trade_no"] = t.reqParam.outTradeNo
	return param
}

func (t *wxRefundQuery) ApiName() string {
	return "pay/orderquery"
}

func (t *wxRefundQuery) QueryPOST() bool {
	return true
}

func (t *wxRefundQuery) GetResult(wx *kwechat.WeClient, body string) (interface{}, error) {
	var resReturn WeChatRefundQueryResp

	err := xml.Unmarshal([]byte(body), &resReturn)
	if err != nil {
		return resReturn, err
	}

	if !(strings.ToUpper(resReturn.ReturnCode) == kpkg.WeChatCodeSuccess && strings.ToUpper(resReturn.ResultCode) == kpkg.WeChatCodeSuccess) {
		return resReturn, errors.New(fmt.Sprintf("wechat RefundQuery error status, result body: %+v", resReturn))
	}
	return resReturn, nil
}
