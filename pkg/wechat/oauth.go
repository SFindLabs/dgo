package wechat

import (
	khttp "dgo/framework/tools/http"
	kpkg "dgo/pkg"
	kcommon "dgo/work/utils"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

//通过 code 获取 access_token、用户个人信息等
var WxOauthGlobal WxOauthInfo

type WeChatCode struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type WxOauthUser struct {
	OpenId  string `json:"openid"`
	UnionId string `json:"unionid"`
}

type WxOauthToken struct {
	WxOauthUser
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type WxOauthUserInfo struct {
	WxOauthUser
	NickName   string   `json:"nickname"`
	Sex        int64    `json:"sex"` //普通用户性别，1 为男性，2 为女性
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	HeadImgUrl string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"` //用户特权信息，json 数组，如微信沃卡用户为（chinaunicom）
}

type WxOauthInfo struct {
	Code  WeChatCode
	Token WxOauthToken
	Info  WxOauthUserInfo
}

//===================================数据处理方法==================================================

//通过 code 获取 access_token
func (wx WxOauthInfo) GetAccessToken(code, appId, appSecret string) (info WxOauthInfo, err error) {
	url := fmt.Sprintf("%s/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code", kpkg.WeChatUrl, appId, appSecret, code)
	jsonParsed, body, err := khttp.UrlGetGetJsonObj(url, 10)
	if err != nil {
		return
	}
	errCode := kcommon.AssertFloat64(jsonParsed.Path("errcode").Data())
	if errCode > 0 {
		info.Code.ErrCode = int64(errCode)
		info.Code.ErrMsg = kcommon.AssertString(jsonParsed.Path("errmsg").Data())
		return
	}
	err = json.Unmarshal([]byte(body), &info.Token)
	return
}

//刷新或续期 access_token 使用
func (wx WxOauthInfo) RefreshAccessToken(refreshToken, appId string) (info WxOauthInfo, err error) {
	url := fmt.Sprintf("%s/sns/oauth2/refresh_token?appid=%s&grant_type=refresh_token&refresh_token=%s", kpkg.WeChatUrl, appId, refreshToken)
	jsonParsed, body, err := khttp.UrlGetGetJsonObj(url, 10)
	if err != nil {
		return
	}

	errCode := kcommon.AssertFloat64(jsonParsed.Path("errcode").Data())
	if errCode > 0 {
		info.Code.ErrCode = int64(errCode)
		info.Code.ErrMsg = kcommon.AssertString(jsonParsed.Path("errmsg").Data())
		return
	}
	err = json.Unmarshal([]byte(body), &info.Token)

	return
}

//检验授权凭证（access_token）是否有效
func (wx WxOauthInfo) CheckAccessToken(accessToken, openid string) (bool, error) {
	url := fmt.Sprintf("%s/sns/auth?access_token=%s&openid=%s", kpkg.WeChatUrl, accessToken, openid)
	jsonParsed, _, err := khttp.UrlGetGetJsonObj(url, 10)
	if err != nil {
		return false, err
	}

	errCode := kcommon.AssertFloat64(jsonParsed.Path("errcode").Data())
	if errCode > 0 {
		return false, errors.New(kcommon.AssertString(jsonParsed.Path("errmsg").Data()))
	}
	return true, nil
}

//获取用户个人信息（UnionID 机制）
func (wx WxOauthInfo) GetUserInfo(accessToken, openid string) (info WxOauthInfo, err error) {
	url := fmt.Sprintf("%s/sns/userinfo?access_token=%s&openid=%s", kpkg.WeChatUrl, accessToken, openid)
	jsonParsed, body, err := khttp.UrlGetGetJsonObj(url, 10)
	if err != nil {
		return
	}

	errCode := kcommon.AssertFloat64(jsonParsed.Path("errcode").Data())
	if errCode > 0 {
		info.Code.ErrCode = int64(errCode)
		info.Code.ErrMsg = kcommon.AssertString(jsonParsed.Path("errmsg").Data())
		return
	}
	err = json.Unmarshal([]byte(body), &info.Info)
	return
}

//通过 code 获取 access_token并且获取用户个人信息（UnionID 机制）
func (wx WxOauthInfo) GetTokenAndInfo(code, appId, appSecret string) (info WxOauthInfo, err error) {
	token, err := wx.GetAccessToken(code, appId, appSecret)
	if err != nil {
		return
	}
	if token.Code.ErrCode != 0 {
		info.Code = token.Code
		return
	}
	info.Token = token.Token
	user, err := wx.GetUserInfo(token.Token.AccessToken, token.Token.OpenId)
	if err != nil {
		return
	}
	if user.Code.ErrCode != 0 {
		info.Code = user.Code
		return
	}
	info.Info = user.Info
	return
}

//网页通过 code 获取 access_token并且获取用户个人信息（UnionID 机制）
func (wx WxOauthInfo) GetWebTokenAndInfo(code, appId, appSecret string) (info WxOauthInfo, err error) {
	token, err := wx.GetAccessToken(code, appId, appSecret)
	if err != nil {
		return
	}
	if token.Code.ErrCode != 0 {
		info.Code = token.Code
		return
	}
	info.Token = token.Token
	user, err := wx.GetUserInfo(token.Token.AccessToken, token.Token.OpenId)
	if err != nil {
		return
	}
	if user.Code.ErrCode != 0 {
		info.Code = user.Code
		return
	}
	info.Info = user.Info
	return
}
