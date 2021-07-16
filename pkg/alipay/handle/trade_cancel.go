package handle

import (
	kpkg "dgo/pkg"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Jeffail/gabs"
	"strings"
)

//统一收单交易撤销接口(支付交易返回失败或支付系统超时，调用该接口撤销交易)
//https://opendocs.alipay.com/apis/api_1/alipay.trade.cancel

func NewAliTradeCancel() *aliTradeCancel {
	return &aliTradeCancel{}
}

type aliTradeCancel struct {
	reqParam tradeCancelParam
}

//请求参数
type tradeCancelParam struct {
	tradeNo    string //该交易在支付宝系统中的交易流水号
	outTradeNo string //订单支付时传入的商户订单号,和支付宝交易号不能同时为空
}

//结果
type AliPayTradeCancelResp struct {
	Code               string `json:"code"`
	Msg                string `json:"msg"`
	SubCode            string `json:"sub_code"`
	SubMsg             string `json:"sub_msg"`
	TradeNo            string `json:"trade_no"`
	OutTradeNo         string `json:"out_trade_no"`
	RetryFlag          string `json:"retry_flag"`
	Action             string `json:"action"`
	GmtRefundPay       string `json:"gmt_refund_pay"`
	RefundSettlementId string `json:"refund_settlement_id"`
}

func (t *aliTradeCancel) SetTradeNo(tradeNo string) {
	t.reqParam.tradeNo = tradeNo
}

func (t *aliTradeCancel) SetOutTradeNo(outTradeNo string) {
	t.reqParam.outTradeNo = outTradeNo
}

func (t *aliTradeCancel) GetParam() map[string]interface{} {
	param := make(map[string]interface{})
	param["trade_no"] = t.reqParam.tradeNo
	param["out_trade_no"] = t.reqParam.outTradeNo
	return param
}

func (t *aliTradeCancel) GetExtendParam() map[string]interface{} {
	return nil
}

func (t *aliTradeCancel) Method() string {
	return "alipay.trade.cancel"
}

func (t *aliTradeCancel) GetResult(jsonParsed *gabs.Container) (interface{}, error) {
	var result AliPayTradeCancelResp
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
