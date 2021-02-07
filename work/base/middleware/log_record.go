package middleware

import (
	kdaocms "dgo/work/dao/cms"
	kutils "dgo/work/utils"
	"github.com/gin-gonic/gin"
)

func StoreLogRecord(c *gin.Context) {
	userId := c.GetInt64("user")
	userName := c.GetString("user_name")
	path := c.Request.URL.Path
	obj := kdaocms.CmsAdminPermissionsObj.GetByPath(nil, path)
	if obj.ID > 0 {
		if obj.IsRecord == 1 {
			_, _ = kdaocms.CmsAdminOptionLogObj.Insert(nil, kutils.AssertString(userName), kutils.AssertInt64(userId), path, c.Request.Method, obj.Name)
		}
	}

}
