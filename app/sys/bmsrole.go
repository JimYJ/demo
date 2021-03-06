package sys

import (
	"666sites/app/respond"
	"666sites/common"
	"666sites/db/system"
	"666sites/service/log"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

// Role 菜单管理
func Role(c *gin.Context) {
	// log.Println(id)
	title, content := common.GetAlertMsg(c.Query("t"), c.Query("c"))
	menulist, err := system.GetMenulist()
	if err != nil {
		menulist = nil
	}
	c.HTML(200, "role.html", gin.H{
		"menu":         system.GetMenu(),
		"list":         system.GetAllRole(),
		"menulist":     menulist,
		"alerttitle":   title,
		"alertcontext": content,
		"bmspath":      common.BmsPath,
	})
}

// DelRole 删除菜单
func DelRole(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		respond.RediErr("role", common.AlertError, common.AlertParamsError, c)
		return
	}
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	_, err := system.DelRole(id, nowTime)
	if err != nil {
		log.Println(err)
		respond.RediErr("role", common.AlertError, common.AlertParamsError, c)
		return
	}
	respond.RediSuccess("/role", c)
}

// AddRole 新增菜单
func AddRole(c *gin.Context) {
	handelRole(c, false)
}

// EditRole 编辑菜单
func EditRole(c *gin.Context) {
	handelRole(c, true)
}

func handelRole(c *gin.Context, isEdit bool) {
	name := c.PostForm("name")
	if name == "" {
		respond.RediErr("role", common.AlertError, common.AlertParamsError, c)
		return
	}
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	if isEdit {
		id := c.Query("id")
		if _, err := strconv.Atoi(id); err != nil {
			log.Println("role id error:", err)
			respond.RediErr("role", common.AlertError, common.AlertParamsError, c)
			return
		}
		_, err := system.EditRole(name, nowTime, id)
		if err != nil {
			log.Println(err)
			respond.RediErr("role", common.AlertFail, common.AlertSaveFail, c)
			return
		}
		respond.RediSuccess("/role", c)
		return
	}
	_, err := system.AddRole(name, nowTime)
	if err != nil {
		log.Println("add role fail:", err)
		respond.RediErr("role", common.AlertFail, common.AlertSaveFail, c)
		return
	}
	respond.RediSuccess("/role", c)
}

// GetRoleMenu 获取角色权限列表
func GetRoleMenu(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id == "" || !common.CheckInt(id) {
		c.JSON(200, gin.H{
			"list": nil,
		})
		return
	}
	list := system.GetRoleMenu(id)
	c.JSON(200, gin.H{
		"list": list,
	})
}

// RoleBindMenu 管理用户岗位
func RoleBindMenu(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	rolemenu := c.PostFormArray("rolemenu")
	if id == "" || !common.CheckInt(id) || rolemenu == nil {
		respond.RediErr("role", common.AlertError, common.AlertParamsError, c)
		return
	}
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	err := system.RoleBindMenu(id, nowTime, rolemenu)
	if err != nil {
		log.Println(err)
		respond.RediErr("role", common.AlertFail, common.AlertSaveFail, c)
		return
	}
	cookie, _ := c.Request.Cookie("c")
	token := cookie.Value
	system.SetMenu(token)
	respond.RediSuccess("/role", c)
}
