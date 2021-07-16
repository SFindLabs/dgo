package alipay

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/url"
)

//支付回调等
var AliCallBackGlobal *aliCallBack

type aliCallBack struct {
}

//===================================数据处理方法==================================================
//  支付宝app支付第三方回调
//  http://testapi.qianfuniuniu.com/fapi/pay/alipay/notify
//  https://opendocs.alipay.com/open/204/105301

type appPayNotifyBind struct {
	NotifyTime string `form:"notify_time"  binding:"required"`  //通知的发送时间,格式为yyyy-MM-dd HH:mm:ss
	NotifyType string `form:"notify_type"  binding:"required"`  //通知的类型
	NotifyId   string `form:"notify_id"  binding:"required"`    //通知校验ID
	AppId      string `form:"app_id"  binding:"required"`       //支付宝分配给开发者的应用Id
	Charset    string `form:"charset"  binding:"required"`      //编码格式。如 utf-8、gbk、gb2312 等
	Version    string `form:"version"  binding:"required"`      //调用的接口版本
	SignType   string `form:"sign_type"  binding:"required"`    //签名类型,商户生成签名字符串所使用的签名算法类型
	Sign       string `form:"sign"  binding:"required"`         //签名
	TradeNo    string `form:"trade_no"  binding:"required"`     //支付宝交易号
	OutTradeNo string `form:"out_trade_no"  binding:"required"` //商户订单号,原支付请求的商户订单号

	OutBizNo          string  `form:"out_biz_no"  binding:"-"`          //商户业务号,商户业务ID(主要是退款通知中返回退款申请的流水号)
	BuyerId           string  `form:"buyer_id"  binding:"-"`            //买家支付宝账号对应的支付宝唯一用户号
	BuyerLogonId      string  `form:"buyer_logon_id"  binding:"-"`      //买家支付宝账号
	SellerId          string  `form:"seller_id"  binding:"-"`           //卖家支付宝用户号
	SellerEmail       string  `form:"seller_email"  binding:"-"`        //卖家支付宝账号
	TradeStatus       string  `form:"trade_status"  binding:"-"`        //交易状态,交易目前所处的状态
	TotalAmount       float64 `form:"total_amount"  binding:"-"`        //订单金额,本次交易支付的订单金额,单位为人民币（元）
	ReceiptAmount     float64 `form:"receipt_amount"  binding:"-"`      //实收金额,商家在交易中实际收到的款项,单位为人民币（元）
	InvoiceAmount     float64 `form:"invoice_amount"  binding:"-"`      //开票金额,用户在交易中支付的可开发票的金额
	BuyerPayAmount    float64 `form:"buyer_pay_amount"  binding:"-"`    //付款金额,用户在交易中支付的金额
	PointAmount       float64 `form:"point_amount"  binding:"-"`        //集分宝金额,使用集分宝支付的金额
	RefundFee         float64 `form:"refund_fee"  binding:"-"`          //总退款金额,退款通知中,返回总退款金额,单位为元,支持两位小数
	Subject           string  `form:"subject"  binding:"-"`             //订单标题,是请求时对应的参数,原样通知回来
	Body              string  `form:"body"  binding:"-"`                //商品描述,对应请求时的 body参数,原样通知回来
	GmtCreate         string  `form:"gmt_create"  binding:"-"`          //该笔交易创建的时间,格式为 yyyy-MM-dd HH:mm:ss
	GmtPayment        string  `form:"gmt_payment"  binding:"-"`         //该笔交易的买家付款时间,格式为 yyyy-MM-dd HH:mm:ss
	GmtRefund         string  `form:"gmt_refund"  binding:"-"`          //该笔交易的退款时间,格式为 yyyy-MM-dd HH:mm:ss
	GmtClose          string  `form:"gmt_close"  binding:"-"`           //该笔交易结束时间,格式为 yyyy-MM-dd HH:mm:ss
	FundBillList      string  `form:"fund_bill_list"  binding:"-"`      //支付金额信息,支付成功的各个渠道金额信息
	PassBackParams    string  `form:"passback_params"  binding:"-"`     //如果请求时传递了该参数,则会在异步通知时将该参数原样返回
	VoucherDetailList string  `form:"voucher_detail_list"  binding:"-"` //优惠券信息,本交易支付时所使用的所有优惠券信息
}

func (ts *aliCallBack) Notify(c *gin.Context, project string) (string, map[string]interface{}, error) {

	/*returnMsg := "failure"
	defer func() {
		_, _ = fmt.Fprintf(c.Writer, returnMsg)
	}()*/

	data, err := c.GetRawData()
	rawData := string(data)
	param := make(map[string]interface{})

	if err != nil {
		return rawData, param, err
	}

	tmpParam, err := url.ParseQuery(rawData)
	if err != nil {
		return rawData, param, err
	}

	for key, val := range tmpParam {
		if len(val) > 0 {
			param[key] = val[0]
		}
	}

	ok, err := AliGlobal[project].VerifySign(param)
	if err != nil {
		return rawData, param, err
	}
	if !ok {
		return rawData, param, errors.New("alipay app pay notify sign verify fail")
	}

	return rawData, param, nil
}
