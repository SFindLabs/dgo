package handle

import (
	kpkg "dgo/pkg"
	kwechat "dgo/pkg/wechat"
	"encoding/xml"
	"errors"
	"strings"
)

//查询企业付款
//https://pay.weixin.qq.com/wiki/doc/api/tools/mch_pay.php?chapter=14_3

func NewWxTransferInfo() *wxTransferInfo {
	return &wxTransferInfo{}
}

type wxTransferInfo struct {
	reqParam transferInfoParam
}

//请求参数
type transferInfoParam struct {
	partnerTradeNo string //商户订单号，需保持唯一性 (只能是字母或者数字，不能包含有其它字符)
	nonceStr       string //随机字符串，不长于32位
}

//结果
type WeChatTransferInfoResp struct {
	ReturnCode     string `xml:"return_code" json:"return_code"`
	ReturnMsg      string `xml:"return_msg" json:"return_msg"`
	ResultCode     string `xml:"result_code" json:"result_code"`
	ErrCode        string `xml:"err_code" json:"err_code"`
	ErrCodeDes     string `xml:"err_code_des" json:"err_code_des"`
	PartnerTradeNo string `xml:"partner_trade_no" json:"partner_trade_no"`
	AppId          string `xml:"appid" json:"-"`
	MchId          string `xml:"mch_id" json:"-"`
	DetailId       string `xml:"detail_id" json:"detail_id"`
	Status         string `xml:"status" json:"status"`
	Reason         string `xml:"reason" json:"reason"`
	OpenId         string `xml:"openid" json:"-"`
	TransferName   string `xml:"transfer_name" json:"transfer_name"`
	PaymentAmount  int64  `xml:"payment_amount" json:"payment_amount"`
	TransferTime   string `xml:"transfer_time" json:"transfer_time"`
	PaymentTime    string `xml:"payment_time" json:"payment_time"`
	Desc           string `xml:"desc" json:"desc"`
}

func (t *wxTransferInfo) SetPartnerTradeNo(partnerTradeNo string) {
	t.reqParam.partnerTradeNo = partnerTradeNo
}

func (t *wxTransferInfo) SetNonceStr(nonceStr string) {
	t.reqParam.nonceStr = nonceStr
}

func (t *wxTransferInfo) GetParam(wx *kwechat.WeClient) map[string]interface{} {
	param := make(map[string]interface{})
	param["appid"] = wx.GetAppsAppId()
	param["mch_id"] = wx.GetMchId()
	param["nonce_str"] = t.reqParam.nonceStr
	param["partner_trade_no"] = t.reqParam.partnerTradeNo
	return param
}

func (t *wxTransferInfo) ApiName() string {
	return "mmpaymkttransfers/gettransferinfo"
}

func (t *wxTransferInfo) QueryPOST() bool {
	return true
}

func (t *wxTransferInfo) GetResult(wx *kwechat.WeClient, body string) (interface{}, error) {
	var resReturn WeChatTransferInfoResp

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
