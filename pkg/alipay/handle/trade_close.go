package handle

import (
	kpkg "dgo/pkg"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Jeffail/gabs"
	"strings"
)

//统一收单交易关闭接口(用于交易创建后，用户在一定时间内未进行支付，可调用该接口直接将未付款的交易进行关闭。)
//https://opendocs.alipay.com/apis/api_1/alipay.trade.close

func NewAliTradeClose() *aliTradeClose {
	return &aliTradeClose{}
}

type aliTradeClose struct {
	reqParam tradeCloseParam
}

//请求参数
type tradeCloseParam struct {
	tradeNo    string //该交易在支付宝系统中的交易流水号
	outTradeNo string //订单支付时传入的商户订单号,和支付宝交易号不能同时为空
}

//结果
type AliPayTradeCloseResp struct {
	Code       string `json:"code"`
	Msg        string `json:"msg"`
	SubCode    string `json:"sub_code"`
	SubMsg     string `json:"sub_msg"`
	TradeNo    string `json:"trade_no"`
	OutTradeNo string `json:"out_trade_no"`
}

func (t *aliTradeClose) SetTradeNo(tradeNo string) {
	t.reqParam.tradeNo = tradeNo
}

func (t *aliTradeClose) SetOutTradeNo(outTradeNo string) {
	t.reqParam.outTradeNo = outTradeNo
}

func (t *aliTradeClose) GetParam() map[string]interface{} {
	param := make(map[string]interface{})
	param["trade_no"] = t.reqParam.tradeNo
	param["out_trade_no"] = t.reqParam.outTradeNo
	return param
}

func (t *aliTradeClose) GetExtendParam() map[string]interface{} {
	return nil
}

func (t *aliTradeClose) Method() string {
	return "alipay.trade.close"
}

func (t *aliTradeClose) GetResult(jsonParsed *gabs.Container) (interface{}, error) {
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
