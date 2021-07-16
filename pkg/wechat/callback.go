package wechat

import (
	kpkg "dgo/pkg"
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"io"
	"strings"
)

//支付回调等
var WxCallBackGlobal *wxCallBack

type wxCallBack struct {
}

//===================================数据处理方法==================================================
//  微信app支付第三方回调
//  http://testapi.qianfuniuniu.com/fapi/pay/wechat/notify
//  https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=9_7

type weChatNotifyBind struct {
	ReturnCode         string `form:"return_code" xml:"return_code" binding:"required"`            //返回状态码(SUCCESS/FAIL)
	ReturnMsg          string `form:"return_msg" xml:"return_msg" binding:"-"`                     //返回信息
	AppId              string `form:"appid" xml:"appid" binding:"-"`                               //微信分配的小程序ID
	MchId              string `form:"mch_id" xml:"mch_id" binding:"-"`                             //微信支付分配的商户号
	DeviceInfo         string `form:"device_info" xml:"device_info" binding:"-"`                   //微信支付分配的终端设备号
	NonceStr           string `form:"nonce_str" xml:"nonce_str" binding:"-"`                       //随机字符串，不长于32位
	Sign               string `form:"sign" xml:"sign" binding:"-"`                                 //签名
	ResultCode         string `form:"result_code" xml:"result_code" binding:"-"`                   //业务结果(SUCCESS/FAIL)
	ErrCode            string `form:"err_code" xml:"err_code" binding:"-"`                         //错误代码
	ErrCodeDes         string `form:"err_code_des" xml:"err_code_des" binding:"-"`                 //错误返回的信息描述
	Openid             string `form:"openid" xml:"openid" binding:"-"`                             //用户在商户appid下的唯一标识
	IsSubscribe        string `form:"is_subscribe" xml:"is_subscribe" binding:"-"`                 //用户是否关注公众账号，Y-关注，N-未关注
	TradeType          string `form:"trade_type" xml:"trade_type" binding:"-"`                     //交易类型(JSAPI、NATIVE、APP)
	BankType           string `form:"bank_type" xml:"bank_type" binding:"-"`                       //银行类型，采用字符串类型的银行标识
	TotalFee           int64  `form:"total_fee" xml:"total_fee" binding:"-"`                       //订单总金额，单位为分
	SettlementTotalFee int64  `form:"settlement_total_fee" xml:"settlement_total_fee" binding:"-"` //应结订单金额=订单金额-非充值代金券金额
	CashFee            int64  `form:"cash_fee" xml:"cash_fee" binding:"-"`                         //现金支付金额
	CashFeeType        string `form:"cash_fee_type" xml:"cash_fee_type" binding:"-"`               //现金支付货币类型
	TransactionId      string `form:"transaction_id" xml:"transaction_id" binding:"-"`             //微信支付订单号
	OutTradeNo         string `form:"out_trade_no" xml:"out_trade_no" binding:"-"`                 //商户内部订单号，32个字符内(数字、大小写字母_-|*@ )，且在同一个商户号下唯一
	Attach             string `form:"attach" xml:"attach" binding:"-"`                             //透传参数
	TimeEnd            string `form:"time_end" xml:"time_end" binding:"-"`                         //支付完成时间
}

func (ts *wxCallBack) Notify(c *gin.Context, project string) (string, map[string]interface{}, string, error) {
	/*var param weChatNotifyBind
	  if err := c.ShouldBindXML(&param); err != nil {
	  	kinit.LogError.Println(err)
	  	_, _ = fmt.Fprintf(c.Writer, fmt.Sprintf(result, tconf.WeChatCodeFAIL))
	  	return
	  }*/

	//returnMsg := kpkg.WeChatCodeFAIL

	/*defer func() {
		_, _ = fmt.Fprintf(c.Writer, fmt.Sprintf(result, returnMsg))
	}()*/

	result := `<xml><return_code><![CDATA[%s]]></return_code></xml>`
	respMap := make(map[string]interface{})

	data, err := c.GetRawData()
	rawData := string(data)

	if err != nil {
		return rawData, respMap, result, err
	}

	err = xml.Unmarshal(data, (*InterfaceMap)(&respMap))
	if err != nil {
		return rawData, respMap, result, err
	}

	sign, err := WxPayGlobal[project].GetSign(respMap)
	if err != nil {
		return rawData, respMap, result, err
	}

	if sign != respMap["sign"] {
		return rawData, respMap, result, errors.New("wechat app pay notify sign verify fail, create sign: " + sign)
	}

	if !(strings.ToUpper(fmt.Sprint(respMap["return_code"])) == kpkg.WeChatCodeSuccess && strings.ToUpper(fmt.Sprint(respMap["result_code"])) == kpkg.WeChatCodeSuccess) {
		return rawData, respMap, result, errors.New("wechat app pay notify code is fail")
	}

	return rawData, respMap, result, nil
}

//-----------------------------------------------------------------------------------

//重写xml转map
type InterfaceMap map[string]interface{}

type xmlMapEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

func (m InterfaceMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(m) == 0 {
		return nil
	}

	err := e.EncodeToken(start)
	if err != nil {
		return err
	}

	for k, v := range m {
		_ = e.Encode(xmlMapEntry{XMLName: xml.Name{Local: k}, Value: fmt.Sprint(v)})
	}

	return e.EncodeToken(start.End())
}

func (m *InterfaceMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = InterfaceMap{}
	for {
		var e xmlMapEntry

		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		(*m)[e.XMLName.Local] = e.Value
	}
	return nil
}
