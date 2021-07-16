package handle

import (
	kpkg "dgo/pkg"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Jeffail/gabs"
	"strings"
)

//统一收单交易退款查询
//https://opendocs.alipay.com/apis/api_1/alipay.trade.fastpay.refund.query

func NewAliTradeRefundQuery() *aliTradeRefundQuery {
	return &aliTradeRefundQuery{}
}

type aliTradeRefundQuery struct {
	reqParam tradeRefundQueryParam
}

//请求参数
type tradeRefundQueryParam struct {
	tradeNo      string //该交易在支付宝系统中的交易流水号
	outTradeNo   string //订单支付时传入的商户订单号,和支付宝交易号不能同时为空
	outRequestNo string //请求退款接口时，传入的退款请求号，如果在退款请求时未传入，则该值为创建交易时的商户订单号
}

//结果
type refundRoyaltys struct {
	RefundAmount  float64 `json:"refund_amount"`
	ResultCode    string  `json:"result_code"`
	TransOut      string  `json:"trans_out"`
	TransOutEmail string  `json:"trans_out_email"`
	TransIn       string  `json:"trans_in"`
	TransInEmail  string  `json:"trans_in_email"`
}

type AliPayTradeRefundQueryResp struct {
	Code           string         `json:"code"`
	Msg            string         `json:"msg"`
	SubCode        string         `json:"sub_code"`
	SubMsg         string         `json:"sub_msg"`
	TradeNo        string         `json:"trade_no"`
	OutTradeNo     string         `json:"out_trade_no"`
	OutRequestNo   string         `json:"out_request_no"`
	TotalAmount    float64        `json:"total_amount"`
	RefundAmount   float64        `json:"refund_amount"`
	RefundRoyaltys refundRoyaltys `json:"refund_royaltys"`
	GmtRefundPay   string         `json:"gmt_refund_pay"`
}

func (t *aliTradeRefundQuery) SetTradeNo(tradeNo string) {
	t.reqParam.tradeNo = tradeNo
}

func (t *aliTradeRefundQuery) SetOutTradeNo(outTradeNo string) {
	t.reqParam.outTradeNo = outTradeNo
}

func (t *aliTradeRefundQuery) SetOutRequestNo(outRequestNo string) {
	t.reqParam.outRequestNo = outRequestNo
}

func (t *aliTradeRefundQuery) GetParam() map[string]interface{} {
	param := make(map[string]interface{})
	param["trade_no"] = t.reqParam.tradeNo
	param["out_trade_no"] = t.reqParam.outTradeNo
	param["out_request_no"] = t.reqParam.outRequestNo
	return param
}

func (t *aliTradeRefundQuery) GetExtendParam() map[string]interface{} {
	return nil
}

func (t *aliTradeRefundQuery) Method() string {
	return "alipay.trade.fastpay.refund.query"
}

func (t *aliTradeRefundQuery) GetResult(jsonParsed *gabs.Container) (interface{}, error) {
	var result AliPayTradeRefundQueryResp
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
