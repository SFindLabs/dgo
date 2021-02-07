package code

// 不可变参数

const (
	BURST_UPLOAD_TMP_DIR      = "/view/upload/tmp/"
	AVATAR_SAVE_PATH          = "/view/upload/avatar/"
	AVATAR_SAVE_PARENT_PATH   = "/view"
	DEFAULT_AVATAR_STATIC_DIR = "/upload/avatar/"
	DEFAULT_AVATAR            = "/upload/default.jpg"
	//分页数
	PAGE_NUMBER = 10
)

//-----------------------------------------------------

var codeToMsg map[int]string
var codeToChnMsg map[int]string

const (
	SUCCESS_STATUS = 200

	OPERATION_WRONG                 = 20001
	ACCESS_TOKEN_FAIL               = 20002
	ACCESS_TOKEN_EXPIRE             = 20003
	PARAM_WRONG                     = 20004
	WRONG_PASSWORD                  = 20005
	WRONG_PERMISSION_NO_SELECT      = 20006
	WRONG_REPEAT_FORM               = 20007
	PASSWD_TOO_SHORT                = 20008
	WRONG_PASSWD_SAMPLE             = 20009
	MOBILE_NOT_OK                   = 20010
	VERIFYCODE_EXPIRE               = 20011
	VERIFYCODE_ERROR                = 20012
	WRONG_PASSWD_CONFIRM            = 20013
	WRONG_NAME_EXIST                = 20014
	WRONG_SUPER_ROLE_OPERATION      = 20015
	WRONG_NAME_ILLEGAL_CHAR_CHINESE = 20016
	FILE_UPLOAD_FAIL                = 20017
	FILE_UPLOAD_KEEP                = 20018
	WRONG_NAME_ILLEGAL_CHAR         = 20019
	WRONG_PATH_EXIST                = 20020
	WRONG_PERMISSION_EMPTY          = 20021
	WRONG_PERMISSION_PATH_EMPTY     = 20022
	WRONG_PERMISSION_NO_HAVE        = 20023
	WRONG_LOGRECORD_NO_CHECK        = 20024
	FAIL_MERGE_FILE                 = 20025
	FAIL_SUFFIX_FILE                = 20026
	WRONG_LOGIN_PASSWORD            = 20027
	USER_NO_EXISTS                  = 20028
	WRONG_FIRST_NAME_EXIST          = 20029
	WRONG_PATH_NAME_ILLEGAL_CHAR    = 20030
	USER_IS_BAN                     = 20031
	USER_IS_LOGOUT                  = 20032
	WRONG_NAME_ILLEGAL_CHAR_KEY     = 20033
	WRONG_KEY_EXIST                 = 20034
)

func init() {

	codeToChnMsg = make(map[int]string)
	codeToChnMsg[SUCCESS_STATUS] = "操作成功"
	codeToChnMsg[OPERATION_WRONG] = "操作错误"
	codeToChnMsg[PARAM_WRONG] = "参数有误"
	codeToChnMsg[WRONG_PASSWORD] = "旧密码不正确"
	codeToChnMsg[PASSWD_TOO_SHORT] = "新密码不能小于6位或者超过20位"
	codeToChnMsg[WRONG_PASSWD_SAMPLE] = "新密码不能与旧密码相同"
	codeToChnMsg[WRONG_PASSWD_CONFIRM] = "新密码确认失败"
	codeToChnMsg[VERIFYCODE_EXPIRE] = "图片验证码过期"
	codeToChnMsg[VERIFYCODE_ERROR] = "图片验证码错误"
	codeToChnMsg[WRONG_REPEAT_FORM] = "操作不合法"
	codeToChnMsg[WRONG_NAME_EXIST] = "名称已存在"
	codeToChnMsg[WRONG_PATH_EXIST] = "路径已存在"
	codeToChnMsg[WRONG_NAME_ILLEGAL_CHAR_CHINESE] = "名称只能包含中英文"
	codeToChnMsg[WRONG_SUPER_ROLE_OPERATION] = "无权限操作此角色"
	codeToChnMsg[WRONG_NAME_ILLEGAL_CHAR] = "名称只能包含英文数字且开头不能为数字"
	codeToChnMsg[WRONG_PERMISSION_EMPTY] = "此角色权限分配不能为空"
	codeToChnMsg[WRONG_PERMISSION_NO_SELECT] = "请选择用户角色"
	codeToChnMsg[WRONG_PERMISSION_PATH_EMPTY] = "非顶级父菜单路径不能为空"
	codeToChnMsg[WRONG_PERMISSION_NO_HAVE] = "无权限访问"
	codeToChnMsg[WRONG_LOGRECORD_NO_CHECK] = "请勾选需要删除的选项"
	codeToChnMsg[FILE_UPLOAD_FAIL] = "文件上传失败"
	codeToChnMsg[FILE_UPLOAD_KEEP] = "文件上传中"
	codeToChnMsg[FAIL_MERGE_FILE] = "文件合并失败"
	codeToChnMsg[FAIL_SUFFIX_FILE] = "文件上传格式不支持"
	codeToChnMsg[WRONG_LOGIN_PASSWORD] = "登录密码错误"
	codeToChnMsg[USER_NO_EXISTS] = "此账号不存在"
	codeToChnMsg[USER_IS_BAN] = "此账号已被禁用"
	codeToChnMsg[WRONG_FIRST_NAME_EXIST] = "顶级父菜单名称已存在"
	codeToChnMsg[WRONG_PATH_NAME_ILLEGAL_CHAR] = "路径不合法,开头必须是/加小写英文,支持小写英文数字-_/"
	codeToChnMsg[USER_IS_LOGOUT] = "此账号已登出"
	codeToChnMsg[WRONG_NAME_ILLEGAL_CHAR_KEY] = "键的值只能包含英文_"
	codeToChnMsg[WRONG_KEY_EXIST] = "键已经存在"
}
func GetCodeMsg(code int) string {
	if msg, ok := codeToMsg[code]; ok {
		return msg
	}
	return ""
}
func GetCodeChnMsg(code int) string {
	if msg, ok := codeToChnMsg[code]; ok {
		return msg
	}
	return ""
}
