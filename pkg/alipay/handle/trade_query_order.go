package handle

import (
	kpkg "dgo/pkg"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Jeffail/gabs"
	"strings"
)

//统一收单线下交易查询
//https://opendocs.alipay.com/apis/api_1/alipay.trade.query

func NewAliTradeQueryOrder() *aliTradeQueryOrder {
	return &aliTradeQueryOrder{}
}

type aliTradeQueryOrder struct {
	reqParam tradeQueryOrderParam
}

//请求参数
type tradeQueryOrderParam struct {
	tradeNo    string //该交易在支付宝系统中的交易流水号
	outTradeNo string //订单支付时传入的商户订单号,和支付宝交易号不能同时为空
}

//结果
type AliPayTradeQueryOrderResp struct {
	Code           string `json:"code"`
	Msg            string `json:"msg"`
	SubCode        string `json:"sub_code"`
	SubMsg         string `json:"sub_msg"`
	TradeNo        string `json:"trade_no"`
	OutTradeNo     string `json:"out_trade_no"`
	BuyerLogonId   string `json:"buyer_logon_id"`
	TradeStatus    string `json:"trade_status"`
	TotalAmount    string `json:"total_amount"`
	BuyerPayAmount string `json:"buyer_pay_amount"`
	SendPayDate    string `json:"send_pay_date"`
	ReceiptAmount  string `json:"receipt_amount"`
	BuyerUserId    string `json:"buyer_user_id"`
}

func (t *aliTradeQueryOrder) SetTradeNo(tradeNo string) {
	t.reqParam.tradeNo = tradeNo
}

func (t *aliTradeQueryOrder) SetOutTradeNo(outTradeNo string) {
	t.reqParam.outTradeNo = outTradeNo
}

func (t *aliTradeQueryOrder) GetParam() map[string]interface{} {
	param := make(map[string]interface{})
	param["trade_no"] = t.reqParam.tradeNo
	param["out_trade_no"] = t.reqParam.outTradeNo
	return param
}

func (t *aliTradeQueryOrder) GetExtendParam() map[string]interface{} {
	return nil
}

func (t *aliTradeQueryOrder) Method() string {
	return "alipay.trade.query"
}

func (t *aliTradeQueryOrder) GetResult(jsonParsed *gabs.Container) (interface{}, error) {
	var result AliPayTradeCloseResp
	field := fmt.Sprintf("%s_response", strings.Replace(t.Method(), ".", "_", -1))
	resStr := jsonParsed.Search(field).String()
	err := json.Unmarshal([]byte(resStr), &result)
	if err != nil {
		return result, err
	}
	if kpkg.AliPayCodeSuccess != result.Code {
		return result, errors.New(result.SubMsg)
	}

	return result, nil
}
