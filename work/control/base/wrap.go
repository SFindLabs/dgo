package base

import (
	kinit "dgo/work/base/initialize"
	kmiddleware "dgo/work/base/middleware"
	kroute "dgo/work/base/route"
	kcode "dgo/work/code"
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	MIDDLE_TYPE_NO_CHECK_LOGIN   = 1 << 0
	MIDDLE_TYPE_CHECK_LOGIN      = 1 << 1
	MIDDLE_TYPE_CHECK_PERMISSION = 1 << 2
	MIDDLE_TYPE_CHECK_CSRF       = 1 << 3

	MIDDLE_TYPE_CHECK_LOGIN_AND_CSRF      = MIDDLE_TYPE_CHECK_LOGIN | MIDDLE_TYPE_CHECK_CSRF
	MIDDLE_TYPE_CHECK_PERMISSION_AND_CSRF = MIDDLE_TYPE_CHECK_PERMISSION | MIDDLE_TYPE_CHECK_CSRF
)

func Wrap(Method string, Path string, f func(*gin.Context), types int) kroute.RouteWrapStruct {
	wp := kroute.RouteWrapStruct{
		Method: Method,
		Path:   Path,
	}
	wp.F = func(c *gin.Context) {
		method := c.Request.Method
		if types&(MIDDLE_TYPE_CHECK_LOGIN|MIDDLE_TYPE_CHECK_PERMISSION) != 0 {
			err, code := CheckLogin(c)
			if err != nil {
				if method == "POST" || method == "PUT" || method == "DELETE" || method == "PATCH" {
					if code == 2 {
						SendErrorJsonStr(c, kcode.USER_IS_BAN, "")
					} else {
						SendErrorJsonStr(c, kcode.USER_IS_LOGOUT, "")
					}
				} else {
					if code == 2 {
						c.Redirect(http.StatusFound, "/userbanshow")
					} else {
						c.Redirect(http.StatusFound, "/login")
					}
				}
				return
			}

			//检查权限
			if types&MIDDLE_TYPE_CHECK_PERMISSION != 0 {
				if ok, code := kmiddleware.CheckPermissionIsOk(c); !ok {
					//处理删除账号
					if code == 1 {
						session := sessions.Default(c)
						session.Clear()
						_ = session.Save()
					}
					if method == "POST" || method == "PUT" || method == "DELETE" || method == "PATCH" {
						SendErrorJsonStr(c, kcode.WRONG_PERMISSION_NO_HAVE, "")
					} else {
						c.Redirect(http.StatusFound, "/nopermission")
					}
					return
				}
				//记录日志
				kmiddleware.StoreLogRecord(c)
			}

			// 检查csrf token
			if types&MIDDLE_TYPE_CHECK_CSRF != 0 {
				if ok := CheckTokenForm(c, method); !ok {
					SendErrorJsonStr(c, kcode.WRONG_REPEAT_FORM, "")
					return
				}
			}
		}
		f(c)
	}
	return wp
}

func CheckLogin(c *gin.Context) (err error, code int) {
	//设置cookie失效时间从不操作开始计算
	sessionName, _ := kinit.Conf.GetString("session.name")
	cookie, err := c.Request.Cookie(sessionName)
	if err != nil {
		return errors.New("cookie is expires"), 1
	}
	maxAge, _ := kinit.Conf.GetInt("session.max_age")
	cookie.Path = "/"
	cookie.MaxAge = maxAge
	cookie.HttpOnly = true
	c.Writer.Header().Set("Set-Cookie", cookie.String())

	//检查session
	session := sessions.Default(c)
	v := session.Get("user")
	if v == nil {
		return errors.New("session is nil"), 1
	}

	if ok := kmiddleware.CheckUserBanIsOk(v); !ok {
		session.Clear()
		if err := session.Save(); err != nil {
			return errors.New("session save fail"), 2
		}
		return errors.New("user is ban"), 2
	}

	c.Set("user", v)
	c.Set("user_name", session.Get("user_name"))
	c.Set("user_avatar", session.Get("user_avatar"))
	c.Set("token", session.Get("token"))
	return nil, 0
}

func CheckTokenForm(c *gin.Context, method string) bool {
	if method == "POST" || method == "PUT" || method == "DELETE" || method == "PATCH" {
		token := GetParam(c, "token")
		if token == "" {
			token = c.GetHeader("X-CSRF-TOKEN")
		}
		if token == "" {
			return false
		}
		return c.GetString("token") == token
	}
	return true
}
