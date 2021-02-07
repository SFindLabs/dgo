package middleware

import (
	kdaocms "dgo/work/dao/cms"
	kutils "dgo/work/utils"
	"github.com/gin-gonic/gin"
)

func CheckPermissionIsOk(c *gin.Context) (bool, int) {
	userId := c.GetInt64("user")
	if userId == 0 {
		return false, 0
	}
	roleIds := kdaocms.CmsAdminRolesObj.GetIdRoleHasPermissions(nil, c.Request.URL.Path)
	roleId := kdaocms.CmsAdminUsersObj.GetIdUserHasRole(nil, userId)
	//处理删除账号
	if roleId.RoleId == 0 {
		return false, 1
	}
	count := len(roleIds)
	for i := 0; i < count; i++ {
		if roleIds[i].RoleId == roleId.RoleId {
			return true, 0
		}
	}
	return false, 0
}

func CheckUserBanIsOk(id interface{}) bool {
	userId := kutils.AssertInt64(id)
	status := kdaocms.CmsAdminUsersObj.GetStatusById(nil, userId)
	//处理删除账号
	if len(status) == 0 {
		return false
	}
	if status[0] != 1 {
		return false
	}
	return true
}
