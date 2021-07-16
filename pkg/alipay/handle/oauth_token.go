package handle

import (
	"encoding/json"
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/pkg/errors"
	"strings"
)

//换取授权访问令牌
//https://opendocs.alipay.com/apis/api_9/alipay.system.oauth.token

func NewAliOauthToken() *aliOauthToken {
	return &aliOauthToken{}
}

type aliOauthToken struct {
	reqParam oauthTokenParam
}

//请求参数
type oauthTokenParam struct {
	code           string //授权码，用户对应用授权后得到
	refreshToken   string //刷新令牌，上次换取访问令牌时得到
	isRefreshToken bool   //是否刷新获取新授权令牌
}

//结果
type AliPayOauthTokenResp struct {
	Code         string `json:"code"`
	Msg          string `json:"msg"`
	SubCode      string `json:"sub_code"`
	SubMsg       string `json:"sub_msg"`
	UserId       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	ReExpiresIn  int64  `json:"re_expires_in"`
}

func (t *aliOauthToken) SetCode(code string) {
	t.reqParam.code = code
}

func (t *aliOauthToken) SetRefreshToken(refreshToken string) {
	t.reqParam.refreshToken = refreshToken
}

func (t *aliOauthToken) SetIsRefreshToken(isRefreshToken bool) {
	t.reqParam.isRefreshToken = isRefreshToken
}

func (t *aliOauthToken) GetParam() map[string]interface{} {
	return nil
}

func (t *aliOauthToken) GetExtendParam() map[string]interface{} {
	param := make(map[string]interface{})
	param["code"] = t.reqParam.code
	param["refresh_token"] = t.reqParam.refreshToken
	if t.reqParam.isRefreshToken {
		param["grant_type"] = "refresh_token"
	} else {
		param["grant_type"] = "authorization_code"
	}
	return param
}

func (t *aliOauthToken) Method() string {
	return "alipay.system.oauth.token"
}

func (t *aliOauthToken) GetResult(jsonParsed *gabs.Container) (interface{}, error) {
	var result AliPayOauthTokenResp
	field := fmt.Sprintf("%s_response", strings.Replace(t.Method(), ".", "_", -1))
	resStr := jsonParsed.Search(field).String()
	err := json.Unmarshal([]byte(resStr), &result)
	if err != nil {
		return result, err
	}
	if result.AccessToken == "" {
		return result, errors.New("access token is empty")
	}
	return result, nil
}
