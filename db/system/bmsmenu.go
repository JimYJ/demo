package system

import (
	"666sites/common"
	"666sites/service/log"
	"fmt"
)

// GetMenulist 获取多级菜单
func GetMenulist() ([]map[string]interface{}, error) {
	list, err := GetMainMenu(0)
	if err != nil {
		return nil, err
	}
	rs := make([]map[string]interface{}, len(list))
	for i := 0; i < len(list); i++ {
		rs[i] = make(map[string]interface{})
		for k, v := range list[i] {
			rs[i][k] = v
		}
		temp, err := GetMainMenu(list[i]["id"].(int64))
		if err != nil {
			log.Println(err)
			rs[i]["list"] = nil
		} else {
			rs[i]["list"] = temp
		}
	}
	return rs, nil
}

// GetMainMenu 获取分类菜单 0父级
func GetMainMenu(parentid int64) ([]map[string]interface{}, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetResults("select id,name,path,icon from bms_menu where deleted = ? and parentid = ? ORDER BY sort", 0, parentid)
}

// SetMenu 从数据库加载菜单到缓存
func SetMenu(token string) {
	cache := common.GetCache()
	list, err := GetMenulist()
	if err != nil {
		log.Println(err)
	}
	ulist := GetUserMenuList(token)
	for i := 0; i < len(list); i++ {
		if list[i]["list"] != nil {
			sublist := list[i]["list"].([]map[string]interface{})
			if sublist != nil {
				templist := make([]map[string]interface{}, 0)
				for j := 0; j < len(sublist); j++ {
					_, ok := ulist[sublist[j]["path"].(string)]
					if ok {
						sublist[j]["path"] = fmt.Sprintf("%s%s", common.BmsPath, sublist[j]["path"])
						templist = append(templist, sublist[j])
					}
				}
				list[i]["list"] = templist
			} else {
				list[i]["list"] = nil
			}
		}
	}
	cache.Set(common.Sysmenu, list, -1)
}

// GetMenu 获取菜单用于HTML渲染
func GetMenu() []map[string]interface{} {
	cache := common.GetCache()
	menu, err := cache.Get(common.Sysmenu)
	if !err {
		menu = nil
	}
	if menu == nil {
		return nil
	}
	return menu.([]map[string]interface{})
}

//  获取分类菜单
func getAllMenu(parentid string) ([]map[string]interface{}, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetResults("select id,name,path,parentid,icon,createtime,updatetime,deleted from bms_menu where deleted = ? and parentid = ? ORDER BY sort", 0, parentid)
}

//  获取子菜单数量
func getSubMenuCount(parentid int64) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	rs, err := mysqlConn.GetVal("select count(*) from bms_menu where deleted = ? and parentid = ? ORDER BY sort", 0, parentid)
	return rs.(int64), err
}

// GetAllMenu 处理菜单详细列表
func GetAllMenu(parentid string) []map[string]interface{} {
	list, err := getAllMenu(parentid)
	if err != nil {
		return nil
	}
	mysqlConn := common.GetMysqlConn()
	for i := 0; i < len(list); i++ {
		if list[i]["deleted"] == 1 {
			list[i]["delete"] = "是"
		} else {
			list[i]["delete"] = "否"
		}
		if list[i]["parentid"].(int64) != 0 {
			name, err := mysqlConn.GetVal("select name from bms_menu where id = ?", list[i]["parentid"].(int64))
			if err != nil {
				log.Println(err)
			}
			if name != nil {
				list[i]["parentname"] = name.(string)
			} else {
				list[i]["parentname"] = "父级菜单"
				count, err := getSubMenuCount(list[i]["id"].(int64))
				if err != nil {
					log.Println(err)
					list[i]["subcount"] = 0
				} else {
					list[i]["subcount"] = count
				}
			}
		} else {
			list[i]["parentname"] = "父级菜单"
			count, err := getSubMenuCount(list[i]["id"].(int64))
			if err != nil {
				log.Println(err)
				list[i]["subcount"] = 0
			} else {
				list[i]["subcount"] = count
			}
		}
	}
	return list
}

// DelMenu 删除菜单
func DelMenu(id, nowTime string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Delete("update bms_menu set deleted = ?,updatetime = ? where id = ?", 1, nowTime, id)
}

// AddMenu 新增菜单
func AddMenu(name, path, parentid, icon, nowTime string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Insert("insert into bms_menu set name = ?,path = ?,parentid = ?,icon = ?,createtime = ?,updatetime = ?", name, path, parentid, icon, nowTime, nowTime)
}

// EditMenu 编辑菜单
func EditMenu(name, path, parentid, icon, nowTime, id string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Insert("update bms_menu set name = ?,path = ?,parentid = ?,icon = ?,updatetime = ? where id = ?", name, path, parentid, icon, nowTime, id)
}

// ChangeSort 更改排序 upordown:true up|false down
func ChangeSort(id, parentid int64, upordown bool) bool {
	list, err := GetMainMenu(parentid)
	if err != nil {
		log.Println(err)
		return false
	}

	var j int
	for i := 0; i < len(list); i++ {
		if list[i]["id"] == id {
			j = i
			break
		}
	}
	if upordown {
		if j == 0 {
			return false
		}
		list[j]["id"], list[j-1]["id"] = list[j-1]["id"], list[j]["id"]
	} else {
		if j == len(list)-1 {
			return false
		}
		list[j]["id"], list[j+1]["id"] = list[j+1]["id"], list[j]["id"]
	}
	mysqlConn := common.GetMysqlConn()
	tx, err := mysqlConn.Begin()
	if err != nil {
		log.Println("begin tx fail:", err)
		return false
	}
	for i := 0; i < len(list); i++ {
		_, err = tx.Update("update bms_menu set sort = ? where id = ?", i, list[i]["id"])
		if err != nil {
			break
		}
	}
	if err != nil {
		tx.Rollback()
		log.Println("change sort fail:", err)
		return false
	}
	tx.Commit()
	return true
}
