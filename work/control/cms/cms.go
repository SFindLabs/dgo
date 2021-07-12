package cms

import (
	kcommon "dgo/framework/tools/common"
	kgorm "dgo/framework/tools/db/gorm"
	kruntime "dgo/framework/tools/runtime"
	kinit "dgo/work/base/initialize"
	kroute "dgo/work/base/route"
	kcode "dgo/work/code"
	kbase "dgo/work/control/base"
	kdao "dgo/work/dao"
	kdaocms "dgo/work/dao/cms"
	kutils "dgo/work/utils"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
	"time"

	"regexp"
	"strconv"
	"strings"
)

type cms struct {
}

func NewCms() *cms {
	return &cms{}
}

func (ts *cms) Load() []kroute.RouteWrapStruct {
	m := make([]kroute.RouteWrapStruct, 0)
	//登录相关
	m = append(m, kbase.Wrap("GET", "/login", ts.loginpage, kbase.MIDDLE_TYPE_NO_CHECK_LOGIN))
	m = append(m, kbase.Wrap("GET", "/captcha", ts.captcha, kbase.MIDDLE_TYPE_NO_CHECK_LOGIN))
	m = append(m, kbase.Wrap("POST", "/login", ts.login, kbase.MIDDLE_TYPE_NO_CHECK_LOGIN))
	//无权限提示页面
	m = append(m, kbase.Wrap("GET", "/nopermission", ts.nopermissionpage, kbase.MIDDLE_TYPE_NO_CHECK_LOGIN))
	m = append(m, kbase.Wrap("GET", "/userbanshow", ts.userbanshowpage, kbase.MIDDLE_TYPE_NO_CHECK_LOGIN))

	//主页
	m = append(m, kbase.Wrap("GET", "/logout", ts.logout, kbase.MIDDLE_TYPE_CHECK_LOGIN))
	m = append(m, kbase.Wrap("GET", "/home", ts.home, kbase.MIDDLE_TYPE_CHECK_LOGIN))
	m = append(m, kbase.Wrap("GET", "/", ts.index, kbase.MIDDLE_TYPE_CHECK_LOGIN))

	//更改个人信息
	m = append(m, kbase.Wrap("GET", "/changpwd", ts.changpwdpage, kbase.MIDDLE_TYPE_CHECK_LOGIN))
	m = append(m, kbase.Wrap("POST", "/pwdedit", ts.pwdedit, kbase.MIDDLE_TYPE_CHECK_LOGIN_AND_CSRF))
	m = append(m, kbase.Wrap("GET", "/changpic", ts.changpicpage, kbase.MIDDLE_TYPE_CHECK_LOGIN))
	m = append(m, kbase.Wrap("POST", "/checkpic", ts.checkpic, kbase.MIDDLE_TYPE_CHECK_LOGIN_AND_CSRF))
	m = append(m, kbase.Wrap("POST", "/picburstupload", ts.picburstupload, kbase.MIDDLE_TYPE_CHECK_LOGIN_AND_CSRF))
	m = append(m, kbase.Wrap("POST", "/picburstmerge", ts.picburstmerge, kbase.MIDDLE_TYPE_CHECK_LOGIN_AND_CSRF))
	//m = append(m, kbase.Wrap("POST", "/picedit", ts.picedit, kbase.MIDDLE_TYPE_CHECK_LOGIN_AND_CSRF))

	//角色管理
	m = append(m, kbase.Wrap("GET", "/role", ts.rolepage, kbase.MIDDLE_TYPE_CHECK_PERMISSION))
	m = append(m, kbase.Wrap("GET", "/roleaddpage", ts.roleaddpage, kbase.MIDDLE_TYPE_CHECK_PERMISSION))
	m = append(m, kbase.Wrap("POST", "/roleadd", ts.roleadd, kbase.MIDDLE_TYPE_CHECK_PERMISSION_AND_CSRF))
	m = append(m, kbase.Wrap("GET", "/roleeditpage", ts.roleeditpage, kbase.MIDDLE_TYPE_CHECK_PERMISSION))
	m = append(m, kbase.Wrap("POST", "/roleedit", ts.roleedit, kbase.MIDDLE_TYPE_CHECK_PERMISSION_AND_CSRF))
	m = append(m, kbase.Wrap("POST", "/roledel", ts.roledel, kbase.MIDDLE_TYPE_CHECK_PERMISSION_AND_CSRF))

	//用户管理
	m = append(m, kbase.Wrap("GET", "/user", ts.userpage, kbase.MIDDLE_TYPE_CHECK_PERMISSION))
	m = append(m, kbase.Wrap("GET", "/useraddpage", ts.useraddpage, kbase.MIDDLE_TYPE_CHECK_PERMISSION))
	m = append(m, kbase.Wrap("POST", "/useradd", ts.useradd, kbase.MIDDLE_TYPE_CHECK_PERMISSION_AND_CSRF))
	m = append(m, kbase.Wrap("GET", "/usereditpage", ts.usereditpage, kbase.MIDDLE_TYPE_CHECK_PERMISSION))
	m = append(m, kbase.Wrap("POST", "/useredit", ts.useredit, kbase.MIDDLE_TYPE_CHECK_PERMISSION_AND_CSRF))
	m = append(m, kbase.Wrap("POST", "/userban", ts.userban, kbase.MIDDLE_TYPE_CHECK_PERMISSION_AND_CSRF))
	m = append(m, kbase.Wrap("POST", "/userdel", ts.userdel, kbase.MIDDLE_TYPE_CHECK_PERMISSION_AND_CSRF))

	//菜单管理
	m = append(m, kbase.Wrap("GET", "/permission", ts.permissionpage, kbase.MIDDLE_TYPE_CHECK_PERMISSION))
	m = append(m, kbase.Wrap("GET", "/permissionaddpage", ts.permissionaddpage, kbase.MIDDLE_TYPE_CHECK_PERMISSION))
	m = append(m, kbase.Wrap("POST", "/permissionadd", ts.permissionadd, kbase.MIDDLE_TYPE_CHECK_PERMISSION_AND_CSRF))
	m = append(m, kbase.Wrap("GET", "/permissioneditpage", ts.permissioneditpage, kbase.MIDDLE_TYPE_CHECK_PERMISSION))
	m = append(m, kbase.Wrap("POST", "/permissionedit", ts.permissionedit, kbase.MIDDLE_TYPE_CHECK_PERMISSION_AND_CSRF))
	m = append(m, kbase.Wrap("POST", "/permissiondel", ts.permissiondel, kbase.MIDDLE_TYPE_CHECK_PERMISSION_AND_CSRF))

	//权限分配相关
	m = append(m, kbase.Wrap("POST", "/getpermissionsofrole", ts.getpermissionsofrole, kbase.MIDDLE_TYPE_CHECK_PERMISSION_AND_CSRF))
	m = append(m, kbase.Wrap("GET", "/permissionsofrole", ts.permissionsofrolepage, kbase.MIDDLE_TYPE_CHECK_PERMISSION))
	m = append(m, kbase.Wrap("POST", "/permissionsofrolesave", ts.permissionsofrolesave, kbase.MIDDLE_TYPE_CHECK_PERMISSION_AND_CSRF))

	//日志管理
	m = append(m, kbase.Wrap("GET", "/logrecord", ts.logrecordpage, kbase.MIDDLE_TYPE_CHECK_PERMISSION))
	m = append(m, kbase.Wrap("POST", "/logrecorddel", ts.logrecorddel, kbase.MIDDLE_TYPE_CHECK_PERMISSION_AND_CSRF))

	//数据库生成器
	m = append(m, kbase.Wrap("GET", "/tables", ts.tables, kbase.MIDDLE_TYPE_CHECK_PERMISSION))
	m = append(m, kbase.Wrap("POST", "/tables/optimize", ts.optimize, kbase.MIDDLE_TYPE_CHECK_PERMISSION_AND_CSRF))
	m = append(m, kbase.Wrap("GET", "/generatepage", ts.generatepage, kbase.MIDDLE_TYPE_CHECK_PERMISSION))
	m = append(m, kbase.Wrap("POST", "/tables/generate", ts.generate, kbase.MIDDLE_TYPE_CHECK_PERMISSION_AND_CSRF))

	return m
}

//-----------------------------------------------------------------------------------

func (ts *cms) loginpage(c *gin.Context) {
	session := sessions.Default(c)
	v := session.Get("user")
	if v != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}

	token := session.Get("token")
	if token == nil {
		token = kcommon.GetRandomString(40, 0)
		session.Set("token", token)
		_ = session.Save()
	}

	c.HTML(http.StatusOK, "cms/login.html", gin.H{
		"token": token,
	})
}

//-----------------------------------------------------------------------------------

func (ts *cms) captcha(c *gin.Context) {
	session := sessions.Default(c)
	root, _ := kruntime.GetCurrentPath()
	fontsDir := fmt.Sprintf("%s/view/assets/fonts", root)
	_ = kutils.ReadFonts(fontsDir, ".ttf")
	captchaImage := kutils.NewCaptchaImage(160, 40, kutils.RandLightColor())
	idKeyC := kutils.RandText(4)
	_ = captchaImage.DrawNoise(kutils.CaptchaComplexMedium).
		DrawText(idKeyC).
		DrawTextNoise(kutils.CaptchaComplexMedium).Error
	base64Png, _ := captchaImage.SaveImageToBase64String(kutils.ImageFormatPng)
	session.Set("captchaId", strings.ToLower(idKeyC))
	_ = session.Save()
	c.String(http.StatusOK, base64Png)
}

//-----------------------------------------------------------------------------------
// http://127.0.0.1:8082/login?name=admin&passwd=123456

/*type loginBind struct {
	Name    string `form:"name"  binding:"required"`
	Passwd  string `form:"passwd"  binding:"required"`
	Captcha string `form:"captcha"  binding:"required"`
}*/

type loginBind struct {
	Name    string `form:"name"  validate:"required,max=50" label:"用户名"`
	Passwd  string `form:"passwd"  validate:"required,max=255" label:"密码"`
	Captcha string `form:"captcha"  validate:"required" label:"验证码"`
}

func (ts *cms) login(c *gin.Context) {
	/*var param loginBind
	callbackName := kbase.GetParam(c, "callback")
	if err := c.ShouldBind(&param); err != nil {
		kbase.SendErrorJsonStr(c, kcode.PARAM_WRONG, callbackName)
		return
	}*/

	var param loginBind
	callbackName := kbase.GetParam(c, "callback")
	if err := c.Bind(&param); err != nil {
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}
	if err := kutils.ValidateTranslate(param); err != nil {
		kbase.SendErrorParamsJsonStr(c, kcode.OPERATION_WRONG, err, callbackName)
		return
	}

	session := sessions.Default(c)
	idKey := session.Get("captchaId")
	key := fmt.Sprintf("%s", idKey)
	verify := key == strings.ToLower(param.Captcha)
	if !verify {
		kbase.SendErrorJsonStr(c, kcode.VERIFYCODE_ERROR, callbackName)
		return
	}
	obj := kdaocms.CmsAdminUsersObj.GetByName(nil, param.Name)

	if obj.ID <= 0 {
		kbase.SendErrorJsonStr(c, kcode.USER_NO_EXISTS, callbackName)
		return
	}

	if obj.Status != 1 {
		kbase.SendErrorJsonStr(c, kcode.USER_IS_BAN, callbackName)
		return
	}

	if ok := kutils.PasswordVerify(param.Passwd, obj.Password); !ok {
		kbase.SendErrorJsonStr(c, kcode.WRONG_LOGIN_PASSWORD, callbackName)
		return
	}

	//进行session设置
	session.Set("user", obj.ID)
	session.Set("user_name", obj.Name)
	session.Set("user_avatar", obj.Avatar)
	session.Delete("captchaId")
	if err := session.Save(); err != nil {
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}
	_ = kdaocms.CmsAdminUsersObj.UpdateLoginInfoById(nil, obj.ID, c.ClientIP())
	kbase.SendErrorJsonStr(c, kcode.SUCCESS_STATUS, "")

}

func (ts *cms) logout(c *gin.Context) {
	callbackName := kbase.GetParam(c, "callback")
	userId := c.GetInt64("user")
	if userId == 0 {
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}
	kbase.SendErrorJsonStr(c, kcode.SUCCESS_STATUS, "")
}

//-----------------------------------------------------------------------------------

func (ts *cms) index(c *gin.Context) {
	pic := c.GetString("user_avatar")
	userId := c.GetInt64("user")
	name := c.GetString("user_name")
	roleId := kdaocms.CmsAdminUsersObj.GetIdUserHasRole(nil, userId)
	objs := kdaocms.CmsAdminPermissionsObj.GetRoleIdHasPermissions(nil, roleId.RoleId)
	kbase.RenderTokenHtml(c, "cms/index.html", gin.H{
		"lists": kutils.TreeMenu.MenuList(objs),
		"name":  name,
		"pic":   pic,
	})
}

func (ts *cms) home(c *gin.Context) {
	kbase.RenderTokenHtml(c, "cms/home.html", gin.H{})
}

func (ts *cms) nopermissionpage(c *gin.Context) {
	kbase.RenderTokenHtml(c, "cms/common_page.html", gin.H{
		"msg":  "无权限访问此页面",
		"code": 401,
		"wait": 3,
		"url":  "",
	})
}

func (ts *cms) userbanshowpage(c *gin.Context) {
	kbase.RenderTokenHtml(c, "cms/common_page.html", gin.H{
		"msg":  "此账号已被禁用",
		"code": 404,
		"wait": 3,
		"url":  "/login",
	})
}

//-----------------------------------------------------------------------------------

func (ts *cms) rolepage(c *gin.Context) {
	count := kdaocms.CmsAdminRolesObj.CountByAll(nil)
	paginate, toUrl, toPage, pageSize := kutils.Paginate(c, count, map[string]interface{}{})
	objs := kdaocms.CmsAdminRolesObj.GetByAll(nil, count, int64(toPage), int64(pageSize))
	countNum := len(objs)
	for i := 0; i < countNum; i++ {
		objs[i].CreatedAt = kutils.FormatTime(objs[i].CreatedAt)
	}
	kbase.RenderTokenHtml(c, "cms/role_list.html", gin.H{
		"lists":    objs,
		"paginate": paginate,
		"toUrl":    toUrl,
		"count":    countNum,
	})
}

func (ts *cms) roleaddpage(c *gin.Context) {
	kbase.RenderTokenHtml(c, "cms/role_add.html", gin.H{})
}

type roleaddBind struct {
	Name string `form:"name"  validate:"required,max=50" label:"角色名称"`
}

func (ts *cms) roleadd(c *gin.Context) {
	var param roleaddBind
	callbackName := kbase.GetParam(c, "callback")
	if err := c.Bind(&param); err != nil {
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}
	if err := kutils.ValidateTranslate(param); err != nil {
		kbase.SendErrorParamsJsonStr(c, kcode.OPERATION_WRONG, err, callbackName)
		return
	}

	param.Name = strings.TrimSpace(param.Name)
	re := regexp.MustCompile("^[A-Za-z\\x{4e00}-\\x{9fa5}]*$")
	ok := re.Match([]byte(param.Name))
	if !ok {
		kbase.SendErrorJsonStr(c, kcode.WRONG_NAME_ILLEGAL_CHAR_CHINESE, callbackName)
		return
	}

	obj := kdaocms.CmsAdminRolesObj.GetByName(nil, param.Name)
	if obj.ID > 0 {
		kbase.SendErrorJsonStr(c, kcode.WRONG_NAME_EXIST, callbackName)
		return
	}

	if _, err := kdaocms.CmsAdminRolesObj.Insert(nil, param.Name); err != nil {
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}
	kbase.SendErrorJsonStr(c, kcode.SUCCESS_STATUS, "")
}

func (ts *cms) roleeditpage(c *gin.Context) {
	idStr := kbase.GetParam(c, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id == 1 {
		c.Redirect(http.StatusFound, "/role")
		return
	}
	objs := kdaocms.CmsAdminRolesObj.GetById(nil, id)
	kbase.RenderTokenHtml(c, "cms/role_edit.html", gin.H{"obj": objs})
}

type roleeditBind struct {
	ID   int64  `form:"id"  validate:"required,gt=0" label:"角色编号"`
	Name string `form:"name"  validate:"required,max=50" label:"角色名称"`
}

func (ts *cms) roleedit(c *gin.Context) {
	var param roleeditBind
	callbackName := kbase.GetParam(c, "callback")
	if err := c.Bind(&param); err != nil {
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}
	if err := kutils.ValidateTranslate(param); err != nil {
		kbase.SendErrorParamsJsonStr(c, kcode.OPERATION_WRONG, err, callbackName)
		return
	}

	param.Name = strings.TrimSpace(param.Name)
	re := regexp.MustCompile("^[A-Za-z\\x{4e00}-\\x{9fa5}]*$")
	ok := re.Match([]byte(param.Name))
	if !ok {
		kbase.SendErrorJsonStr(c, kcode.WRONG_NAME_ILLEGAL_CHAR_CHINESE, callbackName)
		return
	}

	if param.ID == 1 {
		kbase.SendErrorJsonStr(c, kcode.WRONG_SUPER_ROLE_OPERATION, callbackName)
		return
	}

	obj := kdaocms.CmsAdminRolesObj.GetByName(nil, param.Name)
	if obj.ID > 0 && obj.ID != param.ID {
		kbase.SendErrorJsonStr(c, kcode.WRONG_NAME_EXIST, callbackName)
		return
	}

	if err := kdaocms.CmsAdminRolesObj.UpdateById(nil, param.ID, param.Name); err != nil {
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}
	kbase.SendErrorJsonStr(c, kcode.SUCCESS_STATUS, "")

}

type roledelBind struct {
	ID int64 `form:"id"  validate:"required,gt=0" label:"角色编号"`
}

func (ts *cms) roledel(c *gin.Context) {
	var param roledelBind
	callbackName := kbase.GetParam(c, "callback")
	if err := c.Bind(&param); err != nil {
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}
	if err := kutils.ValidateTranslate(param); err != nil {
		kbase.SendErrorParamsJsonStr(c, kcode.OPERATION_WRONG, err, callbackName)
		return
	}

	if param.ID == 1 {
		kbase.SendErrorJsonStr(c, kcode.WRONG_SUPER_ROLE_OPERATION, callbackName)
		return
	}

	transaction := kgorm.NewTransaction()
	tmp, _ := kinit.GetMysqlConnect("")
	tx := transaction.Begin(tmp)
	defer func() {
		_ = transaction.Defer()
	}()

	if err := kdaocms.CmsAdminUserHasRolesObj.DeleteByRoleId(tx, param.ID); err != nil {
		_ = transaction.Rollback()
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}

	if err := kdaocms.CmsAdminRoleHasPermissionsObj.DeleteByRoleId(tx, param.ID); err != nil {
		_ = transaction.Rollback()
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}

	if err := kdaocms.CmsAdminRolesObj.DeleteById(tx, param.ID); err != nil {
		_ = transaction.Rollback()
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}

	kbase.SendErrorJsonStr(c, kcode.SUCCESS_STATUS, "")

}

//----------------------------------------------------------------------------------------

func (ts *cms) changpwdpage(c *gin.Context) {
	kbase.RenderTokenHtml(c, "cms/chang_pwd.html", gin.H{})
}

type pwdBind struct {
	OldPasswd     string `form:"old_passwd"  validate:"required" label:"旧密码"`
	Passwd        string `form:"passwd"  validate:"required,min=6,max=20,eqfield=ConfirmPasswd" label:"新密码"`
	ConfirmPasswd string `form:"confirm_passwd"  validate:"required,min=6,max=20" label:"确认新密码"`
}

func (ts *cms) pwdedit(c *gin.Context) {
	var param pwdBind
	callbackName := kbase.GetParam(c, "callback")
	if err := c.Bind(&param); err != nil {
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}
	if err := kutils.ValidateTranslate(param); err != nil {
		kbase.SendErrorParamsJsonStr(c, kcode.OPERATION_WRONG, err, callbackName)
		return
	}

	/*lenPwd := len(param.Passwd)
	if lenPwd < 6 || lenPwd > 20 {
		kbase.SendErrorJsonStr(c, kcode.PASSWD_TOO_SHORT, callbackName)
		return
	}

	if param.Passwd != param.ConfirmPasswd {
		kbase.SendErrorJsonStr(c, kcode.WRONG_PASSWD_CONFIRM, callbackName)
		return
	}*/

	if param.Passwd == param.OldPasswd {
		kbase.SendErrorJsonStr(c, kcode.WRONG_PASSWD_SAMPLE, callbackName)
		return
	}

	userId := c.GetInt64("user")
	user := kdaocms.CmsAdminUsersObj.GetById(nil, userId)
	if ok := kutils.PasswordVerify(param.OldPasswd, user.Password); !ok {
		kbase.SendErrorJsonStr(c, kcode.WRONG_PASSWORD, callbackName)
		return
	}

	password, err := kutils.PasswordHash(param.Passwd)
	if err != nil {
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}

	if err := kdaocms.CmsAdminUsersObj.UpdatePasswordById(nil, userId, password); err != nil {
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}
	kbase.SendErrorJsonStr(c, kcode.SUCCESS_STATUS, "")
}

func (ts *cms) changpicpage(c *gin.Context) {
	userId := c.GetInt64("user")
	pic := c.GetString("user_avatar")
	kbase.RenderTokenHtml(c, "cms/chang_pic.html", gin.H{
		"pic":    pic,
		"userId": userId,
	})
}

type checkPicBind struct {
	Filename string `form:"filename"  binding:"required"`
	Dir      string `form:"dir"  binding:"required"`
	Total    int64  `form:"total"  binding:"required"`
	Size     string `form:"size"  binding:"required"`
}

func (ts *cms) checkpic(c *gin.Context) {
	var param checkPicBind
	callbackName := kbase.GetParam(c, "callback")
	if err := c.ShouldBind(&param); err != nil {
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.PARAM_WRONG, callbackName)
		return
	}

	fileExt := strings.ToLower(path.Ext(param.Filename))
	isExt := strings.TrimLeft(fileExt, ".")
	if !(isExt == "jpg" || isExt == "jpeg" || isExt == "png" || isExt == "bmp") {
		kbase.SendErrorJsonStr(c, kcode.FAIL_SUFFIX_FILE, callbackName)
		return
	}

	root, err := kruntime.GetCurrentPath()
	if err != nil {
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}

	userId := c.GetInt64("user")
	obj := kdaocms.CmsBurstRecordObj.GetByUidAndFilename(nil, userId, param.Filename)
	if obj.ID > 0 {
		if param.Size == obj.FileTotalSize {
			kbase.SendErrorOriginJsonStr(c, kcode.SUCCESS_STATUS, gin.H{
				"count": obj.BurstCount,
				"id":    obj.ID,
				"dir":   obj.TempFolderName,
			}, callbackName)
			return
		} else {
			dir := fmt.Sprintf("%s%s%s", root, kcode.BURST_UPLOAD_TMP_DIR, obj.TempFolderName)
			if err := os.RemoveAll(dir); err != nil {
				kinit.LogError.Println(err)
				kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
				return
			}
			if err := kdaocms.CmsBurstRecordObj.DeleteById(nil, obj.ID); err != nil {
				kinit.LogError.Println(err)
				kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
				return
			}
		}
	}
	record, err := kdaocms.CmsBurstRecordObj.Insert(nil, userId, param.Dir, param.Filename, param.Size, param.Total)
	if err != nil {
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}
	kbase.SendErrorOriginJsonStr(c, kcode.SUCCESS_STATUS, gin.H{
		"count": record.BurstCount,
		"id":    record.ID,
	}, callbackName)

}

type picburstuploadBind struct {
	Filename string `form:"filename"  binding:"required"`
	Dir      string `form:"dir"  binding:"required"`
	Index    int64  `form:"index"  binding:"required"`
	RecordId int64  `form:"record_id"  binding:"required"`
}

func (ts *cms) picburstupload(c *gin.Context) {
	var param picburstuploadBind
	callbackName := kbase.GetParam(c, "callback")
	if err := c.ShouldBind(&param); err != nil {
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.PARAM_WRONG, callbackName)
		return
	}

	fileExt := strings.ToLower(path.Ext(param.Filename))
	isExt := strings.TrimLeft(fileExt, ".")
	if !(isExt == "jpg" || isExt == "jpeg" || isExt == "png" || isExt == "bmp") {
		kbase.SendErrorJsonStr(c, kcode.FAIL_SUFFIX_FILE, callbackName)
		return
	}

	_, fileHeader, err := c.Request.FormFile("file")
	if err == nil {
		root, err := kruntime.GetCurrentPath()
		if err != nil {
			kinit.LogError.Println(err)
			kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
			return
		}
		dir := fmt.Sprintf("%s%s%s", root, kcode.BURST_UPLOAD_TMP_DIR, param.Dir)
		err, code, _, _ := kutils.SaveBurstFile(c, dir, param.Index, fileExt, fileHeader)
		if err != nil {
			kinit.LogError.Println(err)
			kbase.SendErrorJsonStr(c, code, callbackName)
			return
		}

		if err := kdaocms.CmsBurstRecordObj.UpdateById(nil, param.RecordId, param.Index); err != nil {
			kinit.LogError.Println(err)
			kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
			return
		}
		kbase.SendErrorJsonStr(c, kcode.FILE_UPLOAD_KEEP, callbackName)
	} else {
		kbase.SendErrorJsonStr(c, kcode.FILE_UPLOAD_FAIL, callbackName)
	}
}

type picburstmergeBind struct {
	Filename string `form:"filename"  binding:"required"`
	Dir      string `form:"dir"  binding:"required"`
	AllCount int64  `form:"all_count"  binding:"required"`
	RecordId int64  `form:"record_id"  binding:"required"`
}

func (ts *cms) picburstmerge(c *gin.Context) {
	var param picburstmergeBind
	callbackName := kbase.GetParam(c, "callback")
	if err := c.ShouldBind(&param); err != nil {
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.PARAM_WRONG, callbackName)
		return
	}

	root, err := kruntime.GetCurrentPath()
	if err != nil {
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}
	dir := fmt.Sprintf("%s%s%s", root, kcode.BURST_UPLOAD_TMP_DIR, param.Dir)
	fileExt := path.Ext(param.Filename)

	fileNum, _, err := kutils.GetOneDirFileNum(dir)
	if err != nil {
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}

	if fileNum != param.AllCount {
		if err := os.RemoveAll(dir); err != nil {
			kinit.LogError.Println(err)
			kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
			return
		}
		if err := kdaocms.CmsBurstRecordObj.DeleteById(nil, param.RecordId); err != nil {
			kinit.LogError.Println(err)
			kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
			return
		}
		kbase.SendErrorJsonStr(c, kcode.FILE_UPLOAD_FAIL, callbackName)
		return
	}

	randNum := kcommon.GetRandomString(6, 0)
	times := strconv.FormatInt(time.Now().UnixNano(), 10)
	saveFile := fmt.Sprintf("%s%s%s", times, randNum, fileExt)
	avatarDir := fmt.Sprintf("%s%s", root, kcode.AVATAR_SAVE_PATH)

	if err := kutils.MergeFile(avatarDir, saveFile, dir, int(fileNum), fileExt); err != nil {
		kinit.LogError.Println(err)
		_ = os.Remove(fmt.Sprintf("%s%s", avatarDir, saveFile))
		kbase.SendErrorJsonStr(c, kcode.FAIL_MERGE_FILE, callbackName)
		return
	}

	if err := os.RemoveAll(dir); err != nil {
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}

	transaction := kgorm.NewTransaction()
	tmp, _ := kinit.GetMysqlConnect("")
	tx := transaction.Begin(tmp)
	defer func() {
		_ = transaction.Defer()
	}()

	if err := kdaocms.CmsBurstRecordObj.DeleteById(tx, param.RecordId); err != nil {
		_ = transaction.Rollback()
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}

	userId := c.GetInt64("user")
	objs := kdaocms.CmsAdminUsersObj.GetById(tx, userId)
	pic := objs.Avatar

	url := fmt.Sprintf("%s%s", kcode.DEFAULT_AVATAR_STATIC_DIR, saveFile)
	if err := kdaocms.CmsAdminUsersObj.UpdatePicById(tx, userId, url); err != nil {
		_ = transaction.Rollback()
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}
	session := sessions.Default(c)
	session.Set("user_avatar", url)
	_ = session.Save()

	if pic != "" && pic != kcode.DEFAULT_AVATAR {
		_ = os.Remove(fmt.Sprintf("%s%s%s", root, kcode.AVATAR_SAVE_PARENT_PATH, pic))
	}

	kbase.SendErrorJsonStr(c, kcode.SUCCESS_STATUS, "")
}

/*func (ts *cms) picedit(c *gin.Context) {
	userId := c.GetInt64("user")
	callbackName := kbase.GetParam(c, "callback")
	_, fileHeader, err := c.Request.FormFile("file")
	if err == nil {
		root, err := kruntime.GetCurrentPath()
		if err != nil {
			kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
			return
		}
		dir := fmt.Sprintf("%s%s", root, kcode.AVATAR_SAVE_PATH)
		err, code, savePath, _ := kutils.SaveFile(c, dir, fileHeader)
		if err != nil {
			kinit.LogError.Println(err)
			kbase.SendErrorJsonStr(c, code, callbackName)
			return
		}

		objs := kdaocms.CmsAdminUsersObj.GetById(nil, userId)
		pic := objs.Avatar

		picPath := fmt.Sprintf("%s%s", kcode.DEFAULT_AVATAR_STATIC_DIR, savePath)
		if err := kdaocms.CmsAdminUsersObj.UpdatePicById(nil, userId, picPath); err != nil {
			kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
			return
		}

		session := sessions.Default(c)
		session.Set("user_avatar", picPath)
		_ = session.Save()

		if pic != "" && pic != kcode.DEFAULT_AVATAR {
			_ = os.Remove(fmt.Sprintf("%s%s%s", root, kcode.AVATAR_SAVE_PARENT_PATH, pic))
		}

	} else {
		kbase.SendErrorJsonStr(c, kcode.FILE_UPLOAD_FAIL, callbackName)
		return
	}

	kbase.SendErrorJsonStr(c, kcode.SUCCESS_STATUS, "")
}*/

//-----------------------------------------------------------------------------------

func (ts *cms) userpage(c *gin.Context) {
	searchName := kbase.GetParam(c, "search_name")
	count := kdaocms.CmsAdminUsersObj.CountByName(nil, searchName)
	params := map[string]interface{}{
		"search_name": searchName,
	}
	paginate, toUrl, toPage, pageSize := kutils.Paginate(c, count, params)
	objs := kdaocms.CmsAdminUsersObj.GetByAllName(nil, count, searchName, int64(toPage), int64(pageSize))
	countNum := len(objs)
	for i := 0; i < countNum; i++ {
		objs[i].CreatedAt = kutils.FormatTime(objs[i].CreatedAt)
		objs[i].LoginAt = kutils.FormatTime(objs[i].LoginAt)
	}
	kbase.RenderTokenHtml(c, "cms/user_list.html", gin.H{
		"search_name": searchName,
		"lists":       objs,
		"paginate":    paginate,
		"toUrl":       toUrl,
		"count":       countNum,
	})
}

func (ts *cms) useraddpage(c *gin.Context) {
	var objs []kdaocms.CmsAdminRoles
	userId := c.GetInt64("user")
	if userId == 1 {
		objs = kdaocms.CmsAdminRolesObj.GetAll(nil)
	} else {
		objs = kdaocms.CmsAdminRolesObj.GetIgnoreAll(nil, 1)
	}
	kbase.RenderTokenHtml(c, "cms/user_add.html", gin.H{"lists": objs})
}

type useraddBind struct {
	Name   string `form:"name"  validate:"required,max=50" label:"用户名"`
	Passwd string `form:"passwd"  validate:"required,min=6,max=20" label:"密码"`
	RoleId int64  `form:"role_id"  validate:"gte=0" label:"角色编号"`
}

func (ts *cms) useradd(c *gin.Context) {
	var param useraddBind
	callbackName := kbase.GetParam(c, "callback")
	if err := c.Bind(&param); err != nil {
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}
	if err := kutils.ValidateTranslate(param); err != nil {
		kbase.SendErrorParamsJsonStr(c, kcode.OPERATION_WRONG, err, callbackName)
		return
	}

	if param.RoleId == 0 {
		kbase.SendErrorJsonStr(c, kcode.WRONG_PERMISSION_NO_SELECT, callbackName)
		return
	}

	/*lenPwd := len(param.Passwd)
	if lenPwd < 6 || lenPwd > 20 {
		kbase.SendErrorJsonStr(c, kcode.PASSWD_TOO_SHORT, callbackName)
		return
	}*/

	obj := kdaocms.CmsAdminUsersObj.GetByName(nil, param.Name)
	if obj.ID > 0 {
		kbase.SendErrorJsonStr(c, kcode.WRONG_NAME_EXIST, callbackName)
		return
	}

	param.Name = strings.TrimSpace(param.Name)
	re := regexp.MustCompile("^[a-zA-Z][A-Za-z0-9]*$")
	ok := re.Match([]byte(param.Name))
	if !ok {
		kbase.SendErrorJsonStr(c, kcode.WRONG_NAME_ILLEGAL_CHAR, callbackName)
		return
	}

	password, err := kutils.PasswordHash(param.Passwd)
	if err != nil {
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}

	transaction := kgorm.NewTransaction()
	tmp, _ := kinit.GetMysqlConnect("")
	tx := transaction.Begin(tmp)
	defer func() {
		_ = transaction.Defer()
	}()

	objs, err := kdaocms.CmsAdminUsersObj.Insert(tx, param.Name, password, kcode.DEFAULT_AVATAR, c.ClientIP())
	if err != nil {
		_ = transaction.Rollback()
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}

	if _, err := kdaocms.CmsAdminUserHasRolesObj.Insert(tx, objs.ID, param.RoleId); err != nil {
		_ = transaction.Rollback()
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}

	kbase.SendErrorJsonStr(c, kcode.SUCCESS_STATUS, "")
}

func (ts *cms) usereditpage(c *gin.Context) {
	var obj []kdaocms.CmsAdminRoles

	idStr := kbase.GetParam(c, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id == 1 {
		c.Redirect(http.StatusFound, "/user")
		return
	}
	objs := kdaocms.CmsAdminUsersObj.GetById(nil, id)

	userId := c.GetInt64("user")
	if userId == 1 {
		obj = kdaocms.CmsAdminRolesObj.GetAll(nil)
	} else {
		obj = kdaocms.CmsAdminRolesObj.GetIgnoreAll(nil, 1)
	}
	userRole := kdaocms.CmsAdminUserHasRolesObj.GetByAdminId(nil, id)
	kbase.RenderTokenHtml(c, "cms/user_edit.html", gin.H{
		"obj":    objs,
		"roleId": userRole.RoleId,
		"lists":  obj,
	})
}

type usereditBind struct {
	ID     int64  `form:"id"  validate:"required,gt=0" label:"用户编号"`
	Passwd string `form:"passwd"  validate:"-" label:"密码"`
	RoleId int64  `form:"role_id"  validate:"gte=0" label:"角色编号"`
}

func (ts *cms) useredit(c *gin.Context) {
	var param usereditBind
	callbackName := kbase.GetParam(c, "callback")
	if err := c.Bind(&param); err != nil {
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}
	if err := kutils.ValidateTranslate(param); err != nil {
		kbase.SendErrorParamsJsonStr(c, kcode.OPERATION_WRONG, err, callbackName)
		return
	}

	if param.RoleId == 0 {
		kbase.SendErrorJsonStr(c, kcode.WRONG_PERMISSION_NO_SELECT, callbackName)
		return
	}

	if param.ID == 1 || param.ID == c.GetInt64("user") {
		kbase.SendErrorJsonStr(c, kcode.WRONG_SUPER_ROLE_OPERATION, callbackName)
		return
	}

	password := ""
	if param.Passwd != "" {
		lenPwd := len(param.Passwd)
		if lenPwd < 6 || lenPwd > 20 {
			kbase.SendErrorJsonStr(c, kcode.PASSWD_TOO_SHORT, callbackName)
			return
		}
		pwd, err := kutils.PasswordHash(param.Passwd)
		if err != nil {
			kinit.LogError.Println(err)
			kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
			return
		}
		password = pwd
	}

	transaction := kgorm.NewTransaction()
	tmp, _ := kinit.GetMysqlConnect("")
	tx := transaction.Begin(tmp)
	defer func() {
		_ = transaction.Defer()
	}()

	if password != "" {
		if err := kdaocms.CmsAdminUsersObj.UpdatePasswordById(tx, param.ID, password); err != nil {
			_ = transaction.Rollback()
			kinit.LogError.Println(err)
			kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
			return
		}
	}

	objs := kdaocms.CmsAdminUserHasRolesObj.GetByAdminId(tx, param.ID)
	if objs.ID > 0 {
		if err := kdaocms.CmsAdminUserHasRolesObj.UpdateByAdminId(tx, param.ID, param.RoleId); err != nil {
			_ = transaction.Rollback()
			kinit.LogError.Println(err)
			kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
			return
		}
	} else {
		if _, err := kdaocms.CmsAdminUserHasRolesObj.Insert(tx, param.ID, param.RoleId); err != nil {
			_ = transaction.Rollback()
			kinit.LogError.Println(err)
			kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
			return
		}
	}

	kbase.SendErrorJsonStr(c, kcode.SUCCESS_STATUS, "")

}

type userbanBind struct {
	ID int64 `form:"id"  validate:"required,gt=0" label:"用户编号"`
}

func (ts *cms) userban(c *gin.Context) {
	var param userbanBind
	var status int64

	callbackName := kbase.GetParam(c, "callback")
	if err := c.Bind(&param); err != nil {
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}
	if err := kutils.ValidateTranslate(param); err != nil {
		kbase.SendErrorParamsJsonStr(c, kcode.OPERATION_WRONG, err, callbackName)
		return
	}

	if param.ID == 1 || param.ID == c.GetInt64("user") {
		kbase.SendErrorJsonStr(c, kcode.WRONG_SUPER_ROLE_OPERATION, callbackName)
		return
	}

	obj := kdaocms.CmsAdminUsersObj.GetById(nil, param.ID)

	if obj.Status == 1 {
		status = 0
	} else if obj.Status == 0 {
		status = 1
	}

	if err := kdaocms.CmsAdminUsersObj.UpdateStatusById(nil, param.ID, status); err != nil {
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}

	kbase.SendErrorJsonStr(c, kcode.SUCCESS_STATUS, "")

}

type userdelBind struct {
	ID int64 `form:"id"  validate:"required,gt=0" label:"用户编号"`
}

func (ts *cms) userdel(c *gin.Context) {
	var param userdelBind
	callbackName := kbase.GetParam(c, "callback")
	if err := c.Bind(&param); err != nil {
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}
	if err := kutils.ValidateTranslate(param); err != nil {
		kbase.SendErrorParamsJsonStr(c, kcode.OPERATION_WRONG, err, callbackName)
		return
	}

	if param.ID == 1 || param.ID == c.GetInt64("user") {
		kbase.SendErrorJsonStr(c, kcode.WRONG_SUPER_ROLE_OPERATION, callbackName)
		return
	}

	transaction := kgorm.NewTransaction()
	tmp, _ := kinit.GetMysqlConnect("")
	tx := transaction.Begin(tmp)
	defer func() {
		_ = transaction.Defer()
	}()

	objs := kdaocms.CmsAdminUsersObj.GetById(nil, param.ID)
	pic := objs.Avatar
	if pic != "" && pic != kcode.DEFAULT_AVATAR {
		root, err := kruntime.GetCurrentPath()
		if err != nil {
			kinit.LogError.Println(err)
			kbase.SendErrorJsonStr(c, kcode.PARAM_WRONG, callbackName)
			return
		}
		url := fmt.Sprintf("%s%s%s", root, kcode.AVATAR_SAVE_PARENT_PATH, objs.Avatar)
		_ = os.Remove(url)
	}

	if err := kdaocms.CmsAdminUsersObj.DeleteById(tx, param.ID); err != nil {
		_ = transaction.Rollback()
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}

	if err := kdaocms.CmsAdminUserHasRolesObj.DeleteByAdminId(tx, param.ID); err != nil {
		_ = transaction.Rollback()
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}

	kbase.SendErrorJsonStr(c, kcode.SUCCESS_STATUS, "")

}

//-----------------------------------------------------------------------------------

func (ts *cms) permissionpage(c *gin.Context) {
	objs := kdaocms.CmsAdminPermissionsObj.GetAll(nil)
	kbase.RenderTokenHtml(c, "cms/permission_list.html", gin.H{
		"lists": kutils.TreeMenu.MenuMerge(objs),
	})
}

func (ts *cms) permissionaddpage(c *gin.Context) {
	idStr := kbase.GetParam(c, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Redirect(http.StatusFound, "/permission")
		return
	}
	objs := kdaocms.CmsAdminPermissionsObj.GetAll(nil)
	kbase.RenderTokenHtml(c, "cms/permission_add.html", gin.H{
		"permissionId": id,
		"lists":        kutils.TreeMenu.MenuMerge(objs),
	})
}

type permissionaddBind struct {
	Name   string `form:"name"  validate:"required,max=50" label:"菜单名称"`
	Path   string `form:"path"  validate:"max=128" label:"路径"`
	PID    int64  `form:"pid"  validate:"gte=0" label:"父级菜单编号"`
	Show   int64  `form:"show"  validate:"required,min=1,max=2" label:"导航展示"`
	Modify int64  `form:"modify"  validate:"required,min=1,max=2" label:"权限传递"`
	Record int64  `form:"record"  validate:"required,min=1,max=2" label:"记录日志"`
}

func (ts *cms) permissionadd(c *gin.Context) {
	var param permissionaddBind
	callbackName := kbase.GetParam(c, "callback")
	if err := c.Bind(&param); err != nil {
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}
	if err := kutils.ValidateTranslate(param); err != nil {
		kbase.SendErrorParamsJsonStr(c, kcode.OPERATION_WRONG, err, callbackName)
		return
	}

	if param.PID != 0 && param.Path == "" {
		kbase.SendErrorJsonStr(c, kcode.WRONG_PERMISSION_PATH_EMPTY, callbackName)
		return
	}

	param.Name = strings.TrimSpace(param.Name)
	param.Path = strings.TrimSpace(param.Path)
	re := regexp.MustCompile("^[A-Za-z\\x{4e00}-\\x{9fa5}]*$")
	ok := re.Match([]byte(param.Name))
	if !ok {
		kbase.SendErrorJsonStr(c, kcode.WRONG_NAME_ILLEGAL_CHAR_CHINESE, callbackName)
		return
	}

	//顶级父菜单不能同名
	if param.PID == 0 {
		obj := kdaocms.CmsAdminPermissionsObj.GetByPidName(nil, param.Name, 0)
		if obj.ID > 0 {
			kbase.SendErrorJsonStr(c, kcode.WRONG_FIRST_NAME_EXIST, callbackName)
			return
		}
	}

	if param.Path != "" {
		re = regexp.MustCompile("^/[a-z][a-z0-9-_/]*$")
		ok = re.Match([]byte(param.Path))
		if !ok {
			kbase.SendErrorJsonStr(c, kcode.WRONG_PATH_NAME_ILLEGAL_CHAR, callbackName)
			return
		}
		paths := kdaocms.CmsAdminPermissionsObj.GetByPath(nil, param.Path)
		if paths.ID > 0 {
			kbase.SendErrorJsonStr(c, kcode.WRONG_PATH_EXIST, callbackName)
			return
		}

	}

	if param.Show == 2 {
		param.Show = 0
	}

	if param.Modify == 2 {
		param.Modify = 0
	}

	if param.Record == 2 {
		param.Record = 0
	}

	if _, err := kdaocms.CmsAdminPermissionsObj.Insert(nil, param.Name, param.PID, param.Path, param.Show, param.Modify, param.Record); err != nil {
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}
	kbase.SendErrorJsonStr(c, kcode.SUCCESS_STATUS, "")
}

func (ts *cms) permissioneditpage(c *gin.Context) {
	idStr := kbase.GetParam(c, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		kinit.LogError.Println(err)
		c.Redirect(http.StatusFound, "/permission")
		return
	}
	objs := kdaocms.CmsAdminPermissionsObj.GetAll(nil)
	obj := kdaocms.CmsAdminPermissionsObj.GetById(nil, id)
	kbase.RenderTokenHtml(c, "cms/permission_edit.html", gin.H{
		"lists": kutils.TreeMenu.MenuMerge(objs),
		"obj":   obj,
	})
}

type permissioneditBind struct {
	ID     int64  `form:"id"  validate:"required,gt=0" label:"权限菜单编号"`
	Name   string `form:"name"  validate:"required,max=50" label:"菜单名称"`
	Path   string `form:"path"  validate:"max=128" label:"路径"`
	PID    int64  `form:"pid"  validate:"gte=0" label:"父级菜单编号"`
	Show   int64  `form:"show"  validate:"required,min=1,max=2" label:"导航展示"`
	Modify int64  `form:"modify"  validate:"required,min=1,max=2" label:"权限传递"`
	Record int64  `form:"record"  validate:"required,min=1,max=2" label:"记录日志"`
}

func (ts *cms) permissionedit(c *gin.Context) {
	var param permissioneditBind
	callbackName := kbase.GetParam(c, "callback")
	if err := c.Bind(&param); err != nil {
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}
	if err := kutils.ValidateTranslate(param); err != nil {
		kbase.SendErrorParamsJsonStr(c, kcode.OPERATION_WRONG, err, callbackName)
		return
	}

	if param.PID != 0 && param.Path == "" {
		kbase.SendErrorJsonStr(c, kcode.WRONG_PERMISSION_PATH_EMPTY, callbackName)
		return
	}

	param.Name = strings.TrimSpace(param.Name)
	param.Path = strings.TrimSpace(param.Path)
	re := regexp.MustCompile("^[A-Za-z\\x{4e00}-\\x{9fa5}]*$")
	ok := re.Match([]byte(param.Name))
	if !ok {
		kbase.SendErrorJsonStr(c, kcode.WRONG_NAME_ILLEGAL_CHAR_CHINESE, callbackName)
		return
	}

	//顶级父菜单不能同名
	if param.PID == 0 {
		obj := kdaocms.CmsAdminPermissionsObj.GetByPidName(nil, param.Name, 0)
		if obj.ID > 0 && obj.ID != param.ID {
			kbase.SendErrorJsonStr(c, kcode.WRONG_FIRST_NAME_EXIST, callbackName)
			return
		}
	}

	if param.Path != "" {
		re = regexp.MustCompile("^/[a-z][a-z0-9-_/]*$")
		ok = re.Match([]byte(param.Path))
		if !ok {
			kbase.SendErrorJsonStr(c, kcode.WRONG_PATH_NAME_ILLEGAL_CHAR, callbackName)
			return
		}
		paths := kdaocms.CmsAdminPermissionsObj.GetByPath(nil, param.Path)
		if paths.ID > 0 && paths.ID != param.ID {
			kbase.SendErrorJsonStr(c, kcode.WRONG_PATH_EXIST, callbackName)
			return
		}
	}

	if param.Show == 2 {
		param.Show = 0
	}

	if param.Modify == 2 {
		param.Modify = 0
	}

	if param.Record == 2 {
		param.Record = 0
	}

	if err := kdaocms.CmsAdminPermissionsObj.UpdateById(nil, param.ID, param.Name, param.PID, param.Path, param.Show, param.Modify, param.Record); err != nil {
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}
	kbase.SendErrorJsonStr(c, kcode.SUCCESS_STATUS, "")

}

type permissiondelBind struct {
	ID int64 `form:"id"  validate:"required,gt=0" label:"权限菜单编号"`
}

func (ts *cms) permissiondel(c *gin.Context) {
	var param permissiondelBind
	callbackName := kbase.GetParam(c, "callback")
	if err := c.Bind(&param); err != nil {
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}
	if err := kutils.ValidateTranslate(param); err != nil {
		kbase.SendErrorParamsJsonStr(c, kcode.OPERATION_WRONG, err, callbackName)
		return
	}

	objs := kdaocms.CmsAdminPermissionsObj.GetAll(nil)
	lists := kutils.TreeMenu.DelMergeId(param.ID, objs)

	ids := []int64{param.ID}
	for _, v := range lists {
		ids = append(ids, v.ID)
	}

	transaction := kgorm.NewTransaction()
	tmp, _ := kinit.GetMysqlConnect("")
	tx := transaction.Begin(tmp)
	defer func() {
		_ = transaction.Defer()
	}()

	if err := kdaocms.CmsAdminPermissionsObj.DeleteByIds(tx, ids); err != nil {
		_ = transaction.Rollback()
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}

	if err := kdaocms.CmsAdminRoleHasPermissionsObj.DeleteByPermissionIds(tx, ids); err != nil {
		_ = transaction.Rollback()
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}

	kbase.SendErrorJsonStr(c, kcode.SUCCESS_STATUS, "")

}

//-----------------------------------------------------------------------------------

func (ts *cms) permissionsofrolepage(c *gin.Context) {
	idStr := kbase.GetParam(c, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		kinit.LogError.Println(err)
		c.Redirect(http.StatusFound, "/role")
		return
	}

	//zTree父子关联选项
	types := kbase.GetParam(c, "selectCheckBoxType")
	if types == "" {
		types = "1"
	}

	userId := c.GetInt64("user")
	if id == 1 && userId != 1 {
		c.Redirect(http.StatusFound, "/nopermission")
		return
	}
	kbase.RenderTokenHtml(c, "cms/role_permission.html", gin.H{"role_id": id, "selectCheckBoxType": types})
}

type getpermissionsofroleBind struct {
	RoleId int64 `form:"role_id"  validate:"required,gt=0" label:"角色编号"`
}

func (ts *cms) getpermissionsofrole(c *gin.Context) {
	var param getpermissionsofroleBind
	var objs []kdaocms.CmsAdminPermissions
	callbackName := kbase.GetParam(c, "callback")
	if err := c.Bind(&param); err != nil {
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}
	if err := kutils.ValidateTranslate(param); err != nil {
		kbase.SendErrorParamsJsonStr(c, kcode.OPERATION_WRONG, err, callbackName)
		return
	}

	userId := c.GetInt64("user")
	if param.RoleId == 1 && userId != 1 {
		kbase.SendErrorJsonStr(c, kcode.WRONG_SUPER_ROLE_OPERATION, callbackName)
		return
	}

	if userId == 1 {
		objs = kdaocms.CmsAdminPermissionsObj.GetAll(nil)
	} else {
		objs = kdaocms.CmsAdminPermissionsObj.GetModifyAll(nil, 1)
	}
	lists := kutils.TreeMenu.MenuMerge(objs)
	roles := kdaocms.CmsAdminRoleHasPermissionsObj.GetMapByRoleId(nil, param.RoleId)

	block := make([]kdaocms.PermissionMenu, 0)
	for _, v := range lists {
		tmp := kdaocms.PermissionMenu{
			ID:      v.ID,
			Name:    v.Name,
			Pid:     v.Pid,
			Level:   v.Level,
			Open:    false,
			Checked: false,
		}
		if _, ok := roles[v.ID]; ok {
			tmp.Checked = true
		}
		if v.Pid == 0 {
			tmp.Open = true
		}
		block = append(block, tmp)
	}

	kbase.SendErrorOriginJsonStr(c, kcode.SUCCESS_STATUS, gin.H{
		"lists": block,
	}, callbackName)
}

type permissionsofrolesaveBind struct {
	RoleId      int64   `form:"role_id"  validate:"required,gt=0" label:"角色编号"`
	Permissions []int64 `form:"permissions"  validate:"-" label:"权限编号集合"`
}

func (ts *cms) permissionsofrolesave(c *gin.Context) {
	var param permissionsofrolesaveBind
	callbackName := kbase.GetParam(c, "callback")
	if err := c.Bind(&param); err != nil {
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}
	if err := kutils.ValidateTranslate(param); err != nil {
		kbase.SendErrorParamsJsonStr(c, kcode.OPERATION_WRONG, err, callbackName)
		return
	}

	userId := c.GetInt64("user")
	if param.RoleId == 1 && userId != 1 {
		kbase.SendErrorJsonStr(c, kcode.WRONG_SUPER_ROLE_OPERATION, callbackName)
		return
	}

	count := len(param.Permissions)
	oldPermission := kdaocms.CmsAdminRoleHasPermissionsObj.GetMapByRoleId(nil, param.RoleId)
	countP := len(oldPermission)
	if count == 0 {
		if param.RoleId == 1 {
			kbase.SendErrorJsonStr(c, kcode.WRONG_PERMISSION_EMPTY, "")
			return
		}
		if countP > 0 {
			if err := kdaocms.CmsAdminRoleHasPermissionsObj.DeleteByRoleId(nil, param.RoleId); err != nil {
				kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
				return
			}
		}
	} else {
		if countP > 0 {
			existsIds := make(map[int64]int64, 0)
			deleteIds := make([]int64, 0)
			saveIds := make([]int64, 0)
			for _, v := range param.Permissions {
				if _, ok := oldPermission[v]; !ok {
					saveIds = append(saveIds, v)
				} else {
					existsIds[v] = v
				}
			}
			for k := range oldPermission {
				if _, ok := existsIds[k]; !ok {
					deleteIds = append(deleteIds, k)
				}
			}
			if len(deleteIds) > 0 {
				if err := kdaocms.CmsAdminRoleHasPermissionsObj.DeleteByRoleIdPermissionIds(nil, param.RoleId, deleteIds); err != nil {
					kinit.LogError.Println(err)
					kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
					return
				}
			}
			param.Permissions = saveIds
		}
		if len(param.Permissions) > 0 {
			if err := kdaocms.CmsAdminRoleHasPermissionsObj.InsertByRoleId(nil, param.RoleId, param.Permissions); err != nil {
				kinit.LogError.Println(err)
				kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
				return
			}
		}
	}

	kbase.SendErrorJsonStr(c, kcode.SUCCESS_STATUS, "")
}

//-----------------------------------------------------------------------------------

func (ts *cms) logrecordpage(c *gin.Context) {
	param := kbase.GetParam(c, "searchName")
	count := kdaocms.CmsAdminOptionLogObj.CountByUserName(nil, param)
	params := map[string]interface{}{
		"searchName": param,
	}
	paginate, toUrl, toPage, pageSize := kutils.Paginate(c, count, params)
	objs := kdaocms.CmsAdminOptionLogObj.GetByUserName(nil, count, param, int64(toPage), int64(pageSize))
	countNum := len(objs)
	for i := 0; i < countNum; i++ {
		objs[i].CreatedAt = kutils.FormatTime(objs[i].CreatedAt)
	}
	kbase.RenderTokenHtml(c, "cms/log_list.html", gin.H{
		"lists":      objs,
		"searchName": param,
		"paginate":   paginate,
		"toUrl":      toUrl,
		"count":      countNum,
	})
}

type logrecorddelBind struct {
	Ids []int64 `form:"ids"  validate:"gt=0" label:"勾选删除的日志记录"`
}

func (ts *cms) logrecorddel(c *gin.Context) {
	var param logrecorddelBind
	callbackName := kbase.GetParam(c, "callback")
	if err := c.Bind(&param); err != nil {
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}
	if err := kutils.ValidateTranslate(param); err != nil {
		kbase.SendErrorParamsJsonStr(c, kcode.OPERATION_WRONG, err, callbackName)
		return
	}

	if len(param.Ids) == 0 {
		kbase.SendErrorJsonStr(c, kcode.WRONG_LOGRECORD_NO_CHECK, "")
		return
	}

	if err := kdaocms.CmsAdminOptionLogObj.DeleteByIds(nil, param.Ids); err != nil {
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}

	kbase.SendErrorJsonStr(c, kcode.SUCCESS_STATUS, "")
}

//-----------------------------------------------------------------------------------

func (ts *cms) tables(c *gin.Context) {
	dbName := kbase.GetParam(c, "dbName")
	var db string
	var tmpDbId int64
	if dbName != "" {
		tmpArr := strings.Split(dbName, "_")
		tmpArrLen := len(tmpArr)
		if tmpArrLen >= 1 {
			db = tmpArr[0]
		}
		if tmpArrLen >= 2 {
			tmpDbId, _ = strconv.ParseInt(tmpArr[1], 10, 64)
		}
	}

	dbGorMArr := make([]string, 0)
	for k, v := range kinit.GorMMap {
		if val, ok := v[kinit.MASTER_DB]; ok {
			for kk := range val {
				dbGorMArr = append(dbGorMArr, fmt.Sprintf("%s_%d", k, kk))
			}
		}
	}

	lists, _ := kdao.ScanData(nil, db, tmpDbId, "show table status")
	kbase.RenderTokenHtml(c, "cms/tables_list.html", gin.H{
		"lists":  lists,
		"count":  len(lists),
		"dbArr":  dbGorMArr,
		"dbName": dbName,
		"toUrl":  c.Request.URL.Path,
	})
}

type optimizeBind struct {
	TableName string `form:"table"  validate:"required" label:"数据表名"`
	Engine    string `form:"engine"  validate:"required" label:"引擎"`
}

func (ts *cms) optimize(c *gin.Context) {
	var param optimizeBind
	callbackName := kbase.GetParam(c, "callback")
	if err := c.Bind(&param); err != nil {
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}
	if err := kutils.ValidateTranslate(param); err != nil {
		kbase.SendErrorParamsJsonStr(c, kcode.OPERATION_WRONG, err, callbackName)
		return
	}
	switch strings.ToUpper(param.Engine) {
	case "INNODB":
		_, _ = kdao.ScanData(nil, "", 0, fmt.Sprintf("alter table %s engine = 'InnoDB'", param.TableName))
	case "MYISAM":
		_, _ = kdao.ScanData(nil, "", 0, fmt.Sprintf("optimize table %s", param.TableName))
	default:
		kbase.SendErrorJsonStr(c, kcode.WRONG_TABLE_ENGINE_OPTIMIZE, callbackName)
		return
	}
	kbase.SendErrorJsonStr(c, kcode.SUCCESS_STATUS, "")
}

func (ts *cms) generatepage(c *gin.Context) {
	tables := make([]map[string]interface{}, 0)
	param := kbase.GetParam(c, "searchName")
	dbName := kbase.GetParam(c, "dbName")
	var dbStr string
	var tmpDbId int64
	if dbName != "" {
		tmpArr := strings.Split(dbName, "_")
		tmpArrLen := len(tmpArr)
		if tmpArrLen >= 1 {
			dbStr = tmpArr[0]
		}
		if tmpArrLen >= 2 {
			tmpDbId, _ = strconv.ParseInt(tmpArr[1], 10, 64)
		}
	}
	if param != "" {
		sql := fmt.Sprintf("select table_name as name from information_schema.tables where table_schema='%s'", param)
		tables, _ = kdao.ScanData(nil, dbStr, tmpDbId, sql)
	}
	db, _ := kdao.ScanData(nil, dbStr, tmpDbId, "select schema_name as name from information_schema.schemata")
	kbase.RenderTokenHtml(c, "cms/generate_sql.html", gin.H{
		"lists":      tables,
		"db":         db,
		"searchName": param,
		"dbName":     dbName,
	})
}

type generateBind struct {
	DbName    string `form:"dbName"  validate:"-" label:"数据库节点"`
	Db        string `form:"db"  validate:"required" label:"数据库名"`
	TableName string `form:"table"  validate:"required" label:"数据表名"`
	IsSplit   int64  `form:"split"  validate:"required,min=1,max=2" label:"分库选项"`
	IsDivide  int64  `form:"divide"  validate:"required,min=1,max=2" label:"分库取余选项"`
	IsRead    int64  `form:"read"  validate:"required,min=1,max=2" label:"读写分离选项"`
}

func (ts *cms) generate(c *gin.Context) {
	var param generateBind
	callbackName := kbase.GetParam(c, "callback")
	if err := c.Bind(&param); err != nil {
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}
	if err := kutils.ValidateTranslate(param); err != nil {
		kbase.SendErrorParamsJsonStr(c, kcode.OPERATION_WRONG, err, callbackName)
		return
	}

	splitBool, divideBool, readBool := false, false, false
	if param.IsSplit == 1 {
		splitBool = true
		if param.IsDivide == 1 {
			divideBool = true
		}
	}

	if param.IsRead == 1 {
		readBool = true
	}

	var dbStr string
	var tmpDbId int64
	if param.DbName != "" {
		tmpArr := strings.Split(param.DbName, "_")
		tmpArrLen := len(tmpArr)
		if tmpArrLen >= 1 {
			dbStr = tmpArr[0]
		}
		if tmpArrLen >= 2 {
			tmpDbId, _ = strconv.ParseInt(tmpArr[1], 10, 64)
		}
	}

	str := kutils.GenerateSql.Run(dbStr, tmpDbId, param.Db, param.TableName, splitBool, divideBool, readBool)
	kbase.SendErrorOriginJsonStr(c, kcode.SUCCESS_STATUS, str, "")
}

//-----------------------------------------------------------------------------------
