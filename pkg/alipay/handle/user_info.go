package handle

import (
	kpkg "dgo/pkg"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Jeffail/gabs"
	"strings"
)

//支付宝会员授权信息
//https://opendocs.alipay.com/apis/api_2/alipay.user.info.share

func NewAliUserInfo() *aliUserInfo {
	return &aliUserInfo{}
}

type aliUserInfo struct {
	reqParam userInfoParam
}

//请求参数
type userInfoParam struct {
	authToken string //针对用户授权接口，获取用户相关数据时，用于标识用户授权关系
}

//结果
type AliPayUserInfoShareResp struct {
	Code     string `json:"code"`
	Msg      string `json:"msg"`
	SubCode  string `json:"sub_code"`
	SubMsg   string `json:"sub_msg"`
	UserId   string `json:"user_id"`
	Province string `json:"province"`
	Avatar   string `json:"avatar"`
	City     string `json:"city"`
	NickName string `json:"nick_name"`
	Gender   string `json:"gender"`
}

func (t *aliUserInfo) SetAuthToken(authToken string) {
	t.reqParam.authToken = authToken
}

func (t *aliUserInfo) GetParam() map[string]interface{} {
	return nil
}

func (t *aliUserInfo) GetExtendParam() map[string]interface{} {
	param := make(map[string]interface{})
	param["auth_token"] = t.reqParam.authToken
	return param
}

func (t *aliUserInfo) Method() string {
	return "alipay.user.info.share"
}

func (t *aliUserInfo) GetResult(jsonParsed *gabs.Container) (interface{}, error) {
	var result AliPayUserInfoShareResp
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
