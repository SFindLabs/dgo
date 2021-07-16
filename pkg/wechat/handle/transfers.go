package handle

import (
	kpkg "dgo/pkg"
	kwechat "dgo/pkg/wechat"
	"encoding/xml"
	"errors"
	"strings"
)

//企业付款
//https://pay.weixin.qq.com/wiki/doc/api/tools/mch_pay.php?chapter=14_2

func NewWxEnterprisePayment() *wxEnterprisePayment {
	return &wxEnterprisePayment{}
}

type wxEnterprisePayment struct {
	reqParam enterpriseParam
}

//请求参数
type enterpriseParam struct {
	partnerTradeNo string  //商户订单号，需保持唯一性 (只能是字母或者数字，不能包含有其它字符)
	userOpenId     string  //商户appid下，某用户的openid
	desc           string  //付款备注，必填。注意：备注中的敏感词会被转成字符*
	reUserName     string  //收款用户真实姓名。 如果校验用户姓名，则必填用户真实姓名 如需电子回单，需要传入收款用户姓名
	nonceStr       string  //随机字符串，不长于32位
	amount         float64 //付款金额，单位为元
	isCheckName    bool    //是否校验用户姓名
}

//结果
type WxEnterprisePaymentResp struct {
	ReturnCode     string `xml:"return_code"`
	ReturnMsg      string `xml:"return_msg"`
	MchAppId       string `xml:"mch_appid"`
	MchId          string `xml:"mchid"`
	NonceStr       string `xml:"nonce_str"`
	ResultCode     string `xml:"result_code"`
	ErrCode        string `xml:"err_code"`
	ErrCodeDes     string `xml:"err_code_des"`
	PartnerTradeNo string `xml:"partner_trade_no"`
	PaymentNo      string `xml:"payment_no"`
	PaymentTime    string `xml:"payment_time"`
	IsNeedRetry    bool   //是否用原商户订单号重试
}

func (t *wxEnterprisePayment) SetPartnerTradeNo(partnerTradeNo string) {
	t.reqParam.partnerTradeNo = partnerTradeNo
}

func (t *wxEnterprisePayment) SetUserOpenId(userOpenId string) {
	t.reqParam.userOpenId = userOpenId
}

func (t *wxEnterprisePayment) SetDesc(desc string) {
	t.reqParam.desc = desc
}

func (t *wxEnterprisePayment) SetReUserName(reUserName string) {
	t.reqParam.reUserName = reUserName
}

func (t *wxEnterprisePayment) SetAmount(amount float64) {
	t.reqParam.amount = amount
}

func (t *wxEnterprisePayment) SetIsCheckName(isCheckName bool) {
	t.reqParam.isCheckName = isCheckName
}

func (t *wxEnterprisePayment) SetNonceStr(nonceStr string) {
	t.reqParam.nonceStr = nonceStr
}

func (t *wxEnterprisePayment) GetParam(wx *kwechat.WeClient) map[string]interface{} {
	param := make(map[string]interface{})
	param["mch_appid"] = wx.GetAppsAppId()
	param["mchid"] = wx.GetMchId()
	param["nonce_str"] = t.reqParam.nonceStr
	param["partner_trade_no"] = t.reqParam.partnerTradeNo
	param["openid"] = t.reqParam.userOpenId
	if t.reqParam.isCheckName {
		param["check_name"] = "FORCE_CHECK"
		param["re_user_name"] = t.reqParam.reUserName
	} else {
		param["check_name"] = "NO_CHECK"
	}
	param["amount"] = t.reqParam.amount * 100
	param["desc"] = t.reqParam.desc
	return param
}

func (t *wxEnterprisePayment) ApiName() string {
	return "mmpaymkttransfers/promotion/transfers"
}

func (t *wxEnterprisePayment) QueryPOST() bool {
	return true
}

func (t *wxEnterprisePayment) GetResult(wx *kwechat.WeClient, body string) (interface{}, error) {
	var resReturn WxEnterprisePaymentResp
	var isNeedRetry bool
	err := xml.Unmarshal([]byte(body), &resReturn)
	if err != nil {
		return resReturn, err
	}

	if !(strings.ToUpper(resReturn.ReturnCode) == kpkg.WeChatCodeSuccess && strings.ToUpper(resReturn.ResultCode) == kpkg.WeChatCodeSuccess) {
		switch strings.ToUpper(resReturn.ErrCode) {
		case "NOTENOUGH", "SYSTEMERROR", "NAME_MISMATCH", "SIGN_ERROR", "FREQ_LIMIT", "MONEY_LIMIT", "CA_ERROR", "V2_ACCOUNT_SIMPLE_BAN", "PARAM_IS_NOT_UTF8", "SENDNUM_LIMIT":
			isNeedRetry = true
		}
		payErrMsg := resReturn.ErrCodeDes
		if payErrMsg == "" {
			payErrMsg = resReturn.ReturnMsg
		}
		payErr := errors.New(payErrMsg)
		resReturn.IsNeedRetry = isNeedRetry
		return resReturn, payErr
	}
	return resReturn, nil
}
