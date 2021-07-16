package handle

import (
	kpkg "dgo/pkg"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Jeffail/gabs"
	"strings"
)

//单笔转账接口(新)
//https://opendocs.alipay.com/apis/api_28/alipay.fund.trans.uni.transfer

func NewAliOneTransfer() *aliOneTransfer {
	return &aliOneTransfer{}
}

type aliOneTransfer struct {
	reqParam oneTransferParam
}

//请求参数
type oneTransferParam struct {
	outBizNo    string  //商家侧唯一订单号
	transAmount float64 //订单总金额，单位为元，精确到小数点后两位
	identity    string  //参与方的唯一标识
	name        string  //参与方真实姓名
	isIdType    bool    //是否是会员ID
}

//结果
type AliPayOneTransferResp struct {
	Code           string `json:"code"`
	Msg            string `json:"msg"`
	SubCode        string `json:"sub_code"`
	SubMsg         string `json:"sub_msg"`
	OutBizNo       string `json:"out_biz_no"`
	OrderId        string `json:"order_id"`
	PayFundOrderId string `json:"pay_fund_order_id"`
	Status         string `json:"status"`
	TransDate      string `json:"trans_date"`
}

func (t *aliOneTransfer) SetOutBizNo(outBizNo string) {
	t.reqParam.outBizNo = outBizNo
}

func (t *aliOneTransfer) SetTransAmount(transAmount float64) {
	t.reqParam.transAmount = transAmount
}

func (t *aliOneTransfer) SetIdentity(identity string) {
	t.reqParam.identity = identity
}

func (t *aliOneTransfer) SetIsIdType(isIdType bool) {
	t.reqParam.isIdType = isIdType
}

func (t *aliOneTransfer) SetName(name string) {
	t.reqParam.name = name
}

func (t *aliOneTransfer) GetParam() map[string]interface{} {
	param := make(map[string]interface{})
	param["product_code"] = "TRANS_ACCOUNT_NO_PWD"
	param["biz_scene"] = "DIRECT_TRANSFER"
	param["out_biz_no"] = t.reqParam.outBizNo
	param["trans_amount"] = t.reqParam.transAmount

	info := make(map[string]interface{})
	if t.reqParam.isIdType {
		info["identity_type"] = "ALIPAY_USER_ID"
	} else {
		info["identity_type"] = "ALIPAY_LOGON_ID"
	}
	info["identity"] = t.reqParam.identity
	info["name"] = t.reqParam.name
	param["payee_info"] = info
	return param
}

func (t *aliOneTransfer) GetExtendParam() map[string]interface{} {
	return nil
}

func (t *aliOneTransfer) Method() string {
	return "alipay.fund.trans.uni.transfer"
}

func (t *aliOneTransfer) GetResult(jsonParsed *gabs.Container) (interface{}, error) {
	var result AliPayOneTransferResp
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
