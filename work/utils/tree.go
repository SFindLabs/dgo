package utils

import (
	kdaocms "dgo/work/dao/cms"
)

var TreeMenu *Menu

type Menu struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Pid       int64  `json:"pid"`
	Level     int64  `json:"level"`
	Path      string `json:"path"`
	IsShow    int64  `json:"is_show"`
	IsRecord  int64  `json:"is_record"`
	IsModify  int64  `json:"is_modify"`
	CreatedAt string `json:"created_at"`
}

type delMenuId struct {
	ID int64 `json:"id"`
}

type treeMenuList struct {
	ID        int64           `json:"id"`
	Name      string          `json:"name"`
	Pid       int64           `json:"pid"`
	Level     int64           `json:"level"`
	Path      string          `json:"path"`
	IsShow    int64           `json:"is_show"`
	IsRecord  int64           `json:"is_record"`
	IsModify  int64           `json:"is_modify"`
	CreatedAt string          `json:"created_at"`
	Child     []*treeMenuList `json:"child"`
}

//获取菜单树形切片
func (m *Menu) MenuList(objs []kdaocms.CmsAdminPermissions) []*treeMenuList {
	return m.getMenu(objs, 0, 0)
}

/**
*   合并菜单树形为切片方便遍历取出
* 	param objs []CmsAdminPermissions 所有菜单
*   return 	[]*Menu
*
* for example:
*    id  name      pid  path   is_show	   created_at
*    1	 权限管理 	0             1	    2020-01-10 02:16:36
*    2	 用户管理	    1	/user	  1	    2020-01-10 02:17:10
*
* [
	{
       id: 1
       name: 权限管理
       pid: 0
       level: 0
	   path:
       is_show: 1
       created_at: 2020-01-10 02:16:36
	}
	{
	  id: 2
	  name: 用户管理
	  pid: 1
	  level: 1
	  path: /user
	  is_show: 1
	  created_at: 2020-01-10 02:17:10
	}
 ]
*/
func (m *Menu) MenuMerge(objs []kdaocms.CmsAdminPermissions) []*Menu {
	tmp := make([]*Menu, 0)
	list := m.getMenu(objs, 0, 0)
	for _, v := range list {
		if v.Pid == 0 {
			tmp = append(tmp, &Menu{
				ID:        v.ID,
				Name:      v.Name,
				Pid:       v.Pid,
				Level:     v.Level,
				Path:      v.Path,
				IsShow:    v.IsShow,
				IsRecord:  v.IsRecord,
				IsModify:  v.IsModify,
				CreatedAt: v.CreatedAt,
			})
			if len(v.Child) != 0 {
				tmp = append(tmp, m.concat(v.ID, v.Child)...)
			}
		}
	}
	return tmp
}

/**
*   递归按顺序添加当前父类菜单下的子菜单(用于合并树形为切片)
* 	param id int64  当前父类菜单id
* 	param list []*treeMenuList  当前父类菜单下的树形子菜单
*   return 	[]*Menu
 */
func (m *Menu) concat(id int64, list []*treeMenuList) []*Menu {
	tmp := make([]*Menu, 0)
	for _, v := range list {
		if id == v.Pid {
			tmp = append(tmp, &Menu{
				ID:        v.ID,
				Name:      v.Name,
				Pid:       v.Pid,
				Level:     v.Level,
				Path:      v.Path,
				IsShow:    v.IsShow,
				IsRecord:  v.IsRecord,
				IsModify:  v.IsModify,
				CreatedAt: v.CreatedAt,
			})
			if len(v.Child) != 0 {
				tmp = append(tmp, m.concat(v.ID, v.Child)...)
			}
		}
	}
	return tmp
}

/**
*   获取当前菜单下的所有子菜单id
*   param id int64   当前菜单id
* 	param objs []CmsAdminPermissions 所有菜单
*   return 	[]*DelMenuId
*
* for example:
*    id  name      pid  path   is_show	   created_at
*    1	 权限管理 	0             1	    2020-01-10 02:16:36
*    2	 用户管理	    1	/user	  1	    2020-01-10 02:17:10
*
* [
	{
       id: 1
	}
	{
	  id: 2
	}
 ]
*/
func (m *Menu) DelMergeId(id int64, objs []kdaocms.CmsAdminPermissions) []*delMenuId {
	tmp := make([]*delMenuId, 0)
	list := m.getMenu(objs, id, 0)
	for _, v := range list {
		if v.Pid == id {
			tmp = append(tmp, &delMenuId{
				ID: v.ID,
			})
			if len(v.Child) != 0 {
				tmp = append(tmp, m.concatId(v.ID, v.Child)...)
			}
		}
	}
	return tmp
}

/**
*   递归按顺序添加当前父类菜单下的子菜单id
* 	param id int64  当前父类菜单id
* 	param list []*treeMenuList  当前父类菜单下的树形子菜单
*   return 	[]*DelMenuId
 */
func (m *Menu) concatId(id int64, list []*treeMenuList) []*delMenuId {
	tmp := make([]*delMenuId, 0)
	for _, v := range list {
		if id == v.Pid {
			tmp = append(tmp, &delMenuId{
				ID: v.ID,
			})
			if len(v.Child) != 0 {
				tmp = append(tmp, m.concatId(v.ID, v.Child)...)
			}
		}
	}
	return tmp
}

/**
*  获取菜单树形方法
* 	param objs []CmsAdminPermissions 所有菜单
* 	param pid int64	  父级菜单ID
* 	param level int64 当前层数
*   return 	[]*treeMenuList
*
* for example:
*    id  name      pid  path   is_show	   created_at
*    1	 权限管理 	0             1	    2020-01-10 02:16:36
*    2	 用户管理	    1	/user	  1	    2020-01-10 02:17:10
*
*    {
       id: 1
       name: 权限管理
       pid: 0
       level: 0
	   path:
       is_show: 1
       created_at: 2020-01-10 02:16:36
       child: [
				{
                  id: 2
				  name: 用户管理
				  pid: 1
				  level: 1
				  path: /user
				  is_show: 1
        		  created_at: 2020-01-10 02:17:10
				  child: []
				}
              ]
      }
*/
func (m *Menu) getMenu(objs []kdaocms.CmsAdminPermissions, pid int64, level int64) []*treeMenuList {
	var menu []Menu
	var treeList []*treeMenuList
	count := len(objs)
	for i := 0; i < count; i++ {
		if objs[i].Pid == pid {
			temp := Menu{
				ID:        objs[i].ID,
				Name:      objs[i].Name,
				Pid:       objs[i].Pid,
				Level:     level,
				Path:      objs[i].Path,
				IsShow:    objs[i].IsShow,
				IsRecord:  objs[i].IsRecord,
				IsModify:  objs[i].IsModify,
				CreatedAt: FormatTime(objs[i].CreatedAt),
			}
			menu = append(menu, temp)
		}
	}
	for _, v := range menu {
		level = v.Level + 1
		child := v.getMenu(objs, v.ID, level)
		node := &treeMenuList{
			ID:        v.ID,
			Name:      v.Name,
			Pid:       v.Pid,
			Level:     v.Level,
			Path:      v.Path,
			IsShow:    v.IsShow,
			IsRecord:  v.IsRecord,
			IsModify:  v.IsModify,
			CreatedAt: v.CreatedAt,
		}
		node.Child = child
		treeList = append(treeList, node)
	}
	return treeList
}
