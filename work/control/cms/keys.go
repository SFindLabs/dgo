package cms

import (
	kinit "dgo/work/base/initialize"
	kroute "dgo/work/base/route"
	kcode "dgo/work/code"
	kbase "dgo/work/control/base"
	kdaocms "dgo/work/dao/cms"
	kutils "dgo/work/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type keys struct {
}

func NewKeys() *keys {
	return &keys{}
}
func (ts *keys) Load() []kroute.RouteWrapStruct {
	m := make([]kroute.RouteWrapStruct, 0)

	//配置管理
	m = append(m, kbase.Wrap("GET", "/keys", ts.keyspage, kbase.MIDDLE_TYPE_CHECK_PERMISSION))
	m = append(m, kbase.Wrap("POST", "/keys/sort", ts.keysort, kbase.MIDDLE_TYPE_CHECK_PERMISSION_AND_CSRF))
	m = append(m, kbase.Wrap("POST", "/keys/delete", ts.keydel, kbase.MIDDLE_TYPE_CHECK_PERMISSION_AND_CSRF))

	m = append(m, kbase.Wrap("GET", "/keys/addpage", ts.keyaddpage, kbase.MIDDLE_TYPE_CHECK_PERMISSION))
	m = append(m, kbase.Wrap("POST", "/keys/add", ts.keyadd, kbase.MIDDLE_TYPE_CHECK_PERMISSION_AND_CSRF))
	m = append(m, kbase.Wrap("GET", "/keys/editpage", ts.keyeditpage, kbase.MIDDLE_TYPE_CHECK_PERMISSION))
	m = append(m, kbase.Wrap("POST", "/keys/edit", ts.keyedit, kbase.MIDDLE_TYPE_CHECK_PERMISSION_AND_CSRF))
	return m
}

//-----------------------------------------------------------------------------------

func (ts *keys) keyspage(c *gin.Context) {
	param := kbase.GetParam(c, "search_name")
	count := kdaocms.CmsKeysObj.CountByKey(nil, param)
	params := map[string]interface{}{
		"search_name": param,
	}
	paginate, toUrl, toPage, pageSize := kutils.Paginate(c, count, params)
	objs := kdaocms.CmsKeysObj.GetAllByKey(nil, param, int64(toPage), int64(pageSize))
	countNum := len(objs)
	for i := 0; i < countNum; i++ {
		objs[i].CreatedAt = kutils.FormatTime(objs[i].CreatedAt)
		objs[i].UpdatedAt = kutils.FormatTime(objs[i].UpdatedAt)
	}
	kbase.RenderTokenHtml(c, "cms/keys_list.html", gin.H{
		"lists":      objs,
		"searchName": param,
		"paginate":   paginate,
		"toUrl":      toUrl,
		"count":      countNum,
	})
}

type keySortBind struct {
	ID      int64 `form:"id"  validate:"required,gt=0" label:"配置编号"`
	SortNum int64 `form:"sort_num"  validate:"-" label:"排序"`
}

func (ts *keys) keysort(c *gin.Context) {
	var param keySortBind
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

	if err := kdaocms.CmsKeysObj.UpdateById(nil, param.ID, map[string]interface{}{"sort_num": param.SortNum}); err != nil {
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}
	kbase.SendErrorJsonStr(c, kcode.SUCCESS_STATUS, "")
}

type keyDelBind struct {
	Ids []int64 `form:"ids"  validate:"gt=0" label:"勾选删除的数据配置"`
}

func (ts *keys) keydel(c *gin.Context) {
	var param keyDelBind
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

	/*if len(param.Ids) == 0 {
		kbase.SendErrorJsonStr(c, kcode.WRONG_LOGRECORD_NO_CHECK, "")
		return
	}*/

	if err := kdaocms.CmsKeysObj.DeleteByIds(nil, param.Ids); err != nil {
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}

	kbase.SendErrorJsonStr(c, kcode.SUCCESS_STATUS, "")
}

func (ts *keys) keyaddpage(c *gin.Context) {
	kbase.RenderTokenHtml(c, "cms/keys_add.html", gin.H{})
}

type keyaddBind struct {
	Name   string `form:"name"  validate:"required,max=64" label:"描述"`
	Keyx1  string `form:"keyx1"  validate:"required,max=64" label:"键1"`
	Keyx2  string `form:"keyx2"  validate:"max=64" label:"键2"`
	Valuex string `form:"valuex"  validate:"required,max=255" label:"值"`
	Status int64  `form:"status"  validate:"required,min=1,max=2" label:"状态"`
}

func (ts *keys) keyadd(c *gin.Context) {
	var param keyaddBind
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

	param.Keyx1 = strings.TrimSpace(param.Keyx1)
	re := regexp.MustCompile("^[A-Za-z_]*$")
	ok := re.Match([]byte(param.Keyx1))
	if !ok {
		kbase.SendErrorJsonStr(c, kcode.WRONG_NAME_ILLEGAL_CHAR_KEY, callbackName)
		return
	}

	if param.Keyx2 != "" {
		param.Keyx2 = strings.TrimSpace(param.Keyx2)
		ok = re.Match([]byte(param.Keyx2))
		if !ok {
			kbase.SendErrorJsonStr(c, kcode.WRONG_NAME_ILLEGAL_CHAR_KEY, callbackName)
			return
		}
	}

	obj := kdaocms.CmsKeysObj.GetByAllKey(nil, param.Keyx1, param.Keyx2)
	if obj.ID > 0 {
		kbase.SendErrorJsonStr(c, kcode.WRONG_KEY_EXIST, callbackName)
		return
	}

	if param.Status == 2 {
		param.Status = 0
	}

	if _, err := kdaocms.CmsKeysObj.Insert(nil, param.Name, param.Keyx1, param.Keyx2, param.Valuex, param.Status); err != nil {
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}
	kbase.SendErrorJsonStr(c, kcode.SUCCESS_STATUS, "")
}

func (ts *keys) keyeditpage(c *gin.Context) {
	idStr := kbase.GetParam(c, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Redirect(http.StatusFound, "/keys")
		return
	}
	objs := kdaocms.CmsKeysObj.GetById(nil, id)
	kbase.RenderTokenHtml(c, "cms/keys_edit.html", gin.H{"obj": objs})
}

type keyeditBind struct {
	ID     int64  `form:"id"  validate:"required,gt=0" label:"配置编号"`
	Name   string `form:"name"  validate:"required,max=64" label:"描述"`
	Keyx1  string `form:"keyx1"  validate:"required,max=64" label:"键1"`
	Keyx2  string `form:"keyx2"  validate:"max=64" label:"键2"`
	Valuex string `form:"valuex"  validate:"required,max=255" label:"值"`
	Status int64  `form:"status"  validate:"required,min=1,max=2" label:"状态"`
}

func (ts *keys) keyedit(c *gin.Context) {
	var param keyeditBind
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

	param.Keyx1 = strings.TrimSpace(param.Keyx1)
	re := regexp.MustCompile("^[A-Za-z_]*$")
	ok := re.Match([]byte(param.Keyx1))
	if !ok {
		kbase.SendErrorJsonStr(c, kcode.WRONG_NAME_ILLEGAL_CHAR_KEY, callbackName)
		return
	}

	if param.Keyx2 != "" {
		param.Keyx2 = strings.TrimSpace(param.Keyx2)
		ok = re.Match([]byte(param.Keyx2))
		if !ok {
			kbase.SendErrorJsonStr(c, kcode.WRONG_NAME_ILLEGAL_CHAR_KEY, callbackName)
			return
		}
	}

	obj := kdaocms.CmsKeysObj.GetByAllKey(nil, param.Keyx1, param.Keyx2)
	if obj.ID > 0 && obj.ID != param.ID {
		kbase.SendErrorJsonStr(c, kcode.WRONG_KEY_EXIST, callbackName)
		return
	}

	if param.Status == 2 {
		param.Status = 0
	}

	if err := kdaocms.CmsKeysObj.UpdateById(nil, param.ID, map[string]interface{}{
		"name":   param.Name,
		"keyx1":  param.Keyx1,
		"keyx2":  param.Keyx2,
		"valuex": param.Valuex,
		"status": param.Status,
	}); err != nil {
		kinit.LogError.Println(err)
		kbase.SendErrorJsonStr(c, kcode.OPERATION_WRONG, callbackName)
		return
	}
	kbase.SendErrorJsonStr(c, kcode.SUCCESS_STATUS, "")

}

//----------------------------------------------------------------------------------------
