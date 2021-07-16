package handle

import (
	kpkg "dgo/pkg"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Jeffail/gabs"
	"strings"
)

//转账业务单据查询接口
//https://opendocs.alipay.com/apis/api_28/alipay.fund.trans.common.query

func NewAliQueryTransCommonOrder() *aliQueryTransCommonOrder {
	return &aliQueryTransCommonOrder{}
}

type aliQueryTransCommonOrder struct {
	reqParam queryTransCommonOrderParam
}

//请求参数
type queryTransCommonOrderParam struct {
	payFundOrderId string //支付宝支付资金流水号
	orderId        string //支付宝转账单据号
	outBizNo       string //商户转账唯一订单号
	productCode    string //传递了out_biz_no则该字段为必传
	bizScene       string //传递了out_biz_no则该字段为必传
}

//结果
type AliPayTransCommonOrderResp struct {
	Code           string `json:"code"`
	Msg            string `json:"msg"`
	SubCode        string `json:"sub_code"`
	SubMsg         string `json:"sub_msg"`
	OrderId        string `json:"order_id"`
	PayFundOrderId string `json:"pay_fund_order_id"`
	OutBizNo       string `json:"out_biz_no"`
	TransAmount    string `json:"trans_amount"`
	Status         string `json:"status"`
	PayDate        string `json:"pay_date"`
	ArrivalTimeEnd string `json:"arrival_time_end"`
	OrderFee       string `json:"order_fee"`
	ErrorCode      string `json:"error_code"`
	FailReason     string `json:"fail_reason"`
}

func (t *aliQueryTransCommonOrder) SetPayFundOrderId(payFundOrderId string) {
	t.reqParam.payFundOrderId = payFundOrderId
}

func (t *aliQueryTransCommonOrder) SetOrderId(orderId string) {
	t.reqParam.orderId = orderId
}

func (t *aliQueryTransCommonOrder) SetOutBizNo(outBizNo string) {
	t.reqParam.outBizNo = outBizNo
}

func (t *aliQueryTransCommonOrder) SetProductCode(productCode string) {
	t.reqParam.productCode = productCode
}

func (t *aliQueryTransCommonOrder) SetBizScene(bizScene string) {
	t.reqParam.bizScene = bizScene
}

func (t *aliQueryTransCommonOrder) GetParam() map[string]interface{} {
	param := make(map[string]interface{})
	param["pay_fund_order_id"] = t.reqParam.payFundOrderId
	param["order_id"] = t.reqParam.orderId
	param["out_biz_no"] = t.reqParam.outBizNo
	param["product_code"] = t.reqParam.productCode
	param["biz_scene"] = t.reqParam.bizScene
	return param
}

func (t *aliQueryTransCommonOrder) GetExtendParam() map[string]interface{} {
	return nil
}

func (t *aliQueryTransCommonOrder) Method() string {
	return "alipay.fund.trans.common.query"
}

func (t *aliQueryTransCommonOrder) GetResult(jsonParsed *gabs.Container) (interface{}, error) {
	var result AliPayTransCommonOrderResp
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
