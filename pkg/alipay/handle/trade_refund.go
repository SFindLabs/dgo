package handle

import (
	kpkg "dgo/pkg"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Jeffail/gabs"
	"strings"
)

//统一收单交易退款接口
//https://opendocs.alipay.com/apis/api_1/alipay.trade.refund

func NewAliTradeRefund() *aliTradeRefund {
	return &aliTradeRefund{}
}

type aliTradeRefund struct {
	reqParam tradeRefundParam
}

//请求参数
type tradeRefundParam struct {
	tradeNo      string  //该交易在支付宝系统中的交易流水号
	outTradeNo   string  //订单支付时传入的商户订单号,和支付宝交易号不能同时为空
	refundAmount float64 //需要退款的金额，该金额不能大于订单金额，单位为元，支持两位小数
	refundReason string  //退款原因说明，商家自定义。
	outRequestNo string  //退款请求号。标识一次退款请求，需要保证在交易号下唯一，如需部分退款，则此参数必传。
}

//结果
type AliPayTradeRefundResp struct {
	Code         string  `json:"code"`
	Msg          string  `json:"msg"`
	SubCode      string  `json:"sub_code"`
	SubMsg       string  `json:"sub_msg"`
	TradeNo      string  `json:"trade_no"`
	OutTradeNo   string  `json:"out_trade_no"`
	BuyerLogonId string  `json:"buyer_logon_id"`
	FundChange   string  `json:"fund_change"`
	RefundFee    float64 `json:"refund_fee"`
	BuyerUserId  string  `json:"buyer_user_id"`
}

func (t *aliTradeRefund) SetTradeNo(tradeNo string) {
	t.reqParam.tradeNo = tradeNo
}

func (t *aliTradeRefund) SetOutTradeNo(outTradeNo string) {
	t.reqParam.outTradeNo = outTradeNo
}

func (t *aliTradeRefund) SetRefundAmount(refundAmount float64) {
	t.reqParam.refundAmount = refundAmount
}

func (t *aliTradeRefund) SetRefundReason(refundReason string) {
	t.reqParam.refundReason = refundReason
}

func (t *aliTradeRefund) SetOutRequestNo(outRequestNo string) {
	t.reqParam.outRequestNo = outRequestNo
}

func (t *aliTradeRefund) GetParam() map[string]interface{} {
	param := make(map[string]interface{})
	param["trade_no"] = t.reqParam.tradeNo
	param["out_trade_no"] = t.reqParam.outTradeNo
	param["refund_amount"] = t.reqParam.refundAmount
	param["refund_reason"] = t.reqParam.refundReason
	param["out_request_no"] = t.reqParam.outRequestNo
	return param
}

func (t *aliTradeRefund) GetExtendParam() map[string]interface{} {
	return nil
}

func (t *aliTradeRefund) Method() string {
	return "alipay.trade.refund"
}

func (t *aliTradeRefund) GetResult(jsonParsed *gabs.Container) (interface{}, error) {
	var result AliPayTradeRefundResp
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
