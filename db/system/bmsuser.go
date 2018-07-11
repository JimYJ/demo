package system

import (
	"bms/common"
	"fmt"
	"github.com/JimYJ/easysql/mysql"
	"log"
)

//CheckPass 检查密码是否正确
func CheckPass(user, pass string) (map[string]string, int) {
	mysqlConn := common.GetMysqlConn()
	uinfo, err := mysqlConn.GetRow(mysql.Statement, "select id,password from bms_user WHERE deleted = 0 and username = ? order by id", user)
	if err != nil {
		log.Println(err)
		return nil, 500
	}
	oldpass := uinfo["password"]
	if oldpass != pass {
		return nil, 401
	}
	return uinfo, 200
}

// DelAdminUser 删除后台管理用户
func DelAdminUser(id, nowTime string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Delete(mysql.Statement, "update bms_user set deleted = ?,updatetime = ? where id = ?", 1, nowTime, id)
}

// AddAdminUser 新增后台管理用户
func AddAdminUser(name, pass, nowTime string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Insert(mysql.Statement, "insert into bms_user set username = ?,password = ?,createtime = ?,updatetime = ?", name, pass, nowTime, nowTime)
}

// EditAdminUser 编辑后台管理用户
func EditAdminUser(name, pass, nowTime, id string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Update(mysql.Statement, "update bms_user set username = ?,password = ?,updatetime = ? where id = ?", name, pass, nowTime, id)
}

//  获取全部用户
func getAllAdminUser() ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetResults(mysql.Statement, "select id,username,createtime,updatetime,deleted from bms_user where deleted = ? ORDER BY id", 0)
}

// GetAllAdminUser 处理用户详细列表
func GetAllAdminUser() []map[string]string {
	list, err := getAllAdminUser()
	if err != nil {
		return nil
	}
	for i := 0; i < len(list); i++ {
		if list[i]["deleted"] == "1" {
			list[i]["delete"] = "是"
		} else {
			list[i]["delete"] = "否"
		}
	}
	return list
}

// GetAdminRole 获取用户绑定的角色
func GetAdminRole(id string) []map[string]string {
	mysqlConn := common.GetMysqlConn()
	list, err := mysqlConn.GetResults(mysql.Statement, "select roleid from bms_userrole where userid = ? order by id", id)
	if err != nil {
		return nil
	}
	return list
}

// AdminBindRole 绑定管理账户岗位
func AdminBindRole(id, nowTime string, list []string) error {
	mysqlConn := common.GetMysqlConn()
	mysqlConn.TxBegin()
	_, err := mysqlConn.TxDelete(mysql.Statement, "delete from bms_userrole where userid = ?", id)
	if err != nil {
		mysqlConn.TxRollback()
		return err
	}
	for i := 0; i < len(list); i++ {
		if !common.CheckInt(list[i]) {
			break
		}
		_, err = mysqlConn.TxInsert(mysql.Statement, "insert into bms_userrole set  userid = ?,roleid = ?,createtime = ?,updatetime = ?", id, list[i], nowTime, nowTime)
		if err != nil {
			log.Println(err)
			break
		}
	}
	if err != nil {
		mysqlConn.TxRollback()
		return err
	}
	mysqlConn.TxCommit()
	return nil
}

// GetUserMenuList 获得用户允许访问的菜单列表
func GetUserMenuList(token string) map[string]bool {
	uid, _ := common.GetUIDByToken(token)
	rolelist := GetAdminRole(uid)
	userMenuList := make(map[string]bool)
	for i := 0; i < len(rolelist); i++ {
		pathlist := GetRoleMenuPath(rolelist[i]["roleid"])
		for j := 0; j < len(pathlist); j++ {
			userMenuList[pathlist[j]["path"]] = true
		}
	}
	cache := common.GetCache()
	k := fmt.Sprintf("UserMenu:%s", uid)
	cache.Set(k, userMenuList, -1)
	return userMenuList
}

// GetUserMenuListByCache 从缓存获取权限菜单
func GetUserMenuListByCache(token string) map[string]bool {
	uid, _ := common.GetUIDByToken(token)
	cache := common.GetCache()
	k := fmt.Sprintf("UserMenu:%s", uid)
	list, _ := cache.Get(k)
	return list.(map[string]bool)
}
