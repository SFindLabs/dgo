package handle

import (
	kpkg "dgo/pkg"
	kwechat "dgo/pkg/wechat"
	"encoding/xml"
	"errors"
	"strings"
)

//查询订单
//https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=9_2&index=4

func NewWxOrderQuery() *wxOrderQuery {
	return &wxOrderQuery{}
}

type wxOrderQuery struct {
	reqParam orderQueryParam
}

//请求参数
type orderQueryParam struct {
	transactionId string //微信的订单号，优先使用
	outTradeNo    string //商户系统内部的订单号，当没提供transaction_id时需要传这个。
	nonceStr      string //随机字符串，不长于32位
}

//结果
type WeChatOrderQueryResp struct {
	ReturnCode     string `xml:"return_code" json:"return_code"`
	ReturnMsg      string `xml:"return_msg" json:"return_msg"`
	AppId          string `xml:"appid" json:"-"`
	MchId          string `xml:"mch_id" json:"-"`
	NonceStr       string `xml:"nonce_str" json:"-"`
	ResultCode     string `xml:"result_code" json:"result_code"`
	ErrCode        string `xml:"err_code" json:"err_code"`
	ErrCodeDes     string `xml:"err_code_des" json:"err_code_des"`
	OpenId         string `xml:"openid" json:"-"`
	IsSubscribe    string `xml:"is_subscribe" json:"is_subscribe"`
	TradeType      string `xml:"trade_type" json:"trade_type"`
	TradeState     string `xml:"trade_state" json:"trade_state"`
	BankType       string `xml:"bank_type" json:"bank_type"`
	TotalFee       int64  `xml:"total_fee" json:"total_fee"`
	CashFee        int64  `xml:"cash_fee" json:"cash_fee"`
	TransactionId  string `xml:"transaction_id" json:"transaction_id"`
	OutTradeNo     string `xml:"out_trade_no" json:"out_trade_no"`
	TimeEnd        string `xml:"time_end" json:"time_end"`
	TradeStateDesc string `xml:"trade_state_desc" json:"trade_state_desc"`
}

func (t *wxOrderQuery) SetTransactionId(transactionId string) {
	t.reqParam.transactionId = transactionId
}

func (t *wxOrderQuery) SetOutTradeNo(outTradeNo string) {
	t.reqParam.outTradeNo = outTradeNo
}

func (t *wxOrderQuery) SetNonceStr(nonceStr string) {
	t.reqParam.nonceStr = nonceStr
}

func (t *wxOrderQuery) GetParam(wx *kwechat.WeClient) map[string]interface{} {
	param := make(map[string]interface{})
	param["appid"] = wx.GetAppsAppId()
	param["mch_id"] = wx.GetMchId()
	param["nonce_str"] = t.reqParam.nonceStr
	param["out_trade_no"] = t.reqParam.outTradeNo
	param["transaction_id"] = t.reqParam.transactionId
	return param
}

func (t *wxOrderQuery) ApiName() string {
	return "pay/orderquery"
}

func (t *wxOrderQuery) QueryPOST() bool {
	return true
}

func (t *wxOrderQuery) GetResult(wx *kwechat.WeClient, body string) (interface{}, error) {
	var resReturn WeChatOrderQueryResp

	err := xml.Unmarshal([]byte(body), &resReturn)
	if err != nil {
		return resReturn, err
	}

	if !(strings.ToUpper(resReturn.ReturnCode) == kpkg.WeChatCodeSuccess && strings.ToUpper(resReturn.ResultCode) == kpkg.WeChatCodeSuccess) {
		payErrMsg := resReturn.ErrCodeDes
		if payErrMsg == "" {
			payErrMsg = resReturn.ReturnMsg
		}
		payErr := errors.New(payErrMsg)
		return resReturn, payErr
	}
	return resReturn, nil
}
