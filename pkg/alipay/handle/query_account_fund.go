package handle

import (
	kpkg "dgo/pkg"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Jeffail/gabs"
	"strings"
)

//支付宝资金账户资产查询接口
//https://opendocs.alipay.com/apis/api_28/alipay.fund.account.query

func NewAliQueryAccountFund() *aliQueryAccountFund {
	return &aliQueryAccountFund{}
}

type aliQueryAccountFund struct {
	reqParam queryAccountFundParam
}

//请求参数
type queryAccountFundParam struct {
	aliPayUserId string //支付宝会员 id
}

//结果
type AliPayFundAccountResp struct {
	Code            string `json:"code"`
	Msg             string `json:"msg"`
	SubCode         string `json:"sub_code"`
	SubMsg          string `json:"sub_msg"`
	AvailableAmount string `json:"available_amount"`
}

func (t *aliQueryAccountFund) SetAliPayUserId(aliPayUserId string) {
	t.reqParam.aliPayUserId = aliPayUserId
}

func (t *aliQueryAccountFund) GetParam() map[string]interface{} {
	param := make(map[string]interface{})
	param["alipay_user_id"] = t.reqParam.aliPayUserId
	param["account_type"] = "ACCTRANS_ACCOUNT"
	return param
}

func (t *aliQueryAccountFund) GetExtendParam() map[string]interface{} {
	return nil
}

func (t *aliQueryAccountFund) Method() string {
	return "alipay.fund.account.query"
}

func (t *aliQueryAccountFund) GetResult(jsonParsed *gabs.Container) (interface{}, error) {
	var result AliPayFundAccountResp
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
