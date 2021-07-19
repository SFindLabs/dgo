package pkg

//================================微信或支付宝==============================================

const (
	//=======================微信=======================

	WeChatUrl              = "https://api.weixin.qq.com"
	WeChatMchSandboxURL    = "https://api.mch.weixin.qq.com/sandboxnew"
	WeChatMchProductionURL = "https://api.mch.weixin.qq.com"

	WeChatProductionOAuthURL = "https://open.weixin.qq.com/connect/oauth2/authorize"

	WeChatCodeSuccess = "SUCCESS"
	WeChatCodeFAIL    = "FAIL"

	//---------------------------应用开发信息---------------------------------------
	//提现
	WeChatAppsProject   = "default"
	WeChatAppsAppId     = "" //移动应用开发者ID
	WeChatAppsAppSecret = "" //移动应用密钥

	//支付
	WeChatPayAppsProject   = "default"
	WeChatPayAppsAppId     = "" //移动应用开发者ID
	WeChatPayAppsAppSecret = "" //移动应用密钥

	// -----------------------------公众号------------------------------------------

	WeChatOfficialAccountAppId     = "" //公众号开发者ID
	WeChatOfficialAccountAppSecret = "" //公众号开发者密钥

	// -----------------------------微信支付------------------------------------------

	WeChatMchId  = "" //微信支付分配的商户号
	WeChatApiKey = "" //商户Key 、支付Key、 API密钥

)

const (
	//=======================支付宝============================

	AliPayProductionURL = "https://openapi.alipay.com/gateway.do"
	AliPaySandboxURL    = "https://openapi.alipaydev.com/gateway.do"

	AliPayProductionOAuthURL = "https://openauth.alipay.com"
	AliPaySandboxOAuthURL    = "https://openauth.alipaydev.com"

	//------------------------应用开发信息-------------------------------
	//提现
	AliPayAppsProject         = "default"
	AliPayProductionAppsAppId = "" //即创建应用后生成

	AliPaySandboxAppsAppId = ""

	// -----------------------------支付------------------------------------------
	AliPayProductionPayAppsAppId = "" //即创建应用后生成

	AliPaySandboxPayAppsAppId = ""

	//签约的支付宝账号对应的支付宝唯一用户号，以 2088 开头的 16 位纯数字组成。
	AliPayPid = ""

	AliPayFormat   = "JSON"
	AliPayCharset  = "utf-8"
	AliPayVersion  = "1.0"
	AliPaySignType = "RSA2"

	//应用公私钥
	AliPayAppPrivateKey = ""
	AliPayAppPublicKey  = ""

	AliPayCodeSuccess          = "10000" // 接口调用成功
	AliPayCodeUnKnowError      = "20000" // 服务不可用
	AliPayCodeInvalidAuthToken = "20001" // 授权权限不足
	AliPayCodeMissingParam     = "40001" // 缺少必选参数
	AliPayCodeInvalidParam     = "40002" // 非法的参数
	AliPayCodeBusinessFailed   = "40004" // 业务处理失败
	AliPayCodePermissionDenied = "40006" // 权限不足
)

//=========================上传================================

const (
	//--------------------------------obs-----------------------------------

	ObsAK       = ""
	ObsSK       = ""
	ObsEndpoint = ""
	ObsBucket   = ""
	ObsUrl      = ""

	//-------------------------------oss------------------------------------

	//endpoint带-internal的是内网，用于阿里云服务器快速上传(外网访问需要设置endpoint为外网[默认为http，需要https的要加上https://])

	OssEndpoint = ""
	OssAk       = ""
	OssSK       = ""
	OssBucket   = ""
)
