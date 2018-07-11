package sys

import (
	"bms/app/respond"
	"bms/common"
	"bms/db/system"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"time"
)

// Menu 菜单管理
func Menu(c *gin.Context) {
	mainlist, _ := system.GetMainMenu("0")
	id := c.DefaultQuery("id", "0")
	title, content := common.GetAlertMsg(c.Query("t"), c.Query("c"))
	// log.Println(id)
	c.HTML(200, "menu.html", gin.H{
		"menu":         system.GetMenu(),
		"list":         system.GetAllMenu(id),
		"mainlist":     mainlist,
		"alerttitle":   title,
		"alertcontext": content,
		"bmspath":      common.BmsPath,
	})
}

// DelMenu 删除菜单
func DelMenu(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		respond.RediErr("menu", common.AlertError, common.AlertParamsError, c)
		return
	}
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	_, err := system.DelMenu(id, nowTime)
	if err != nil {
		log.Println(err)
		respond.RediErr("menu", common.AlertFail, common.AlertDelFail, c)
		return
	}
	system.SetMenu(common.GetTokenByCookie(c))
	respond.RediSuccess("/menu", c)
}

// AddMenu 新增菜单
func AddMenu(c *gin.Context) {
	handelMenu(c, false)
}

// EditMenu 编辑菜单
func EditMenu(c *gin.Context) {
	handelMenu(c, true)
}

func handelMenu(c *gin.Context, isEdit bool) {
	name := c.PostForm("name")
	path := c.PostForm("path")
	icon := c.PostForm("icon")
	parentid := c.PostForm("parentid")
	if name == "" || path == "" || !common.CheckInt(parentid) {
		respond.RediErr("menu", common.AlertError, common.AlertParamsError, c)
		return
	}
	if icon == "" {
		icon = "flaticon-layers"
	}
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	if isEdit {
		id := c.Query("id")
		if _, err := strconv.Atoi(id); err != nil {
			log.Println("menu id error:", err)
			respond.RediErr("menu", common.AlertError, common.AlertParamsError, c)
			return
		}
		_, err := system.EditMenu(name, path, parentid, icon, nowTime, id)
		if err != nil {
			log.Println(err)
			respond.RediErr("menu", common.AlertFail, common.AlertSaveFail, c)
			return
		}
		system.SetMenu(common.GetTokenByCookie(c))
		respond.RediSuccess("/menu", c)
		return
	}
	_, err := system.AddMenu(name, path, parentid, icon, nowTime)
	if err != nil {
		log.Println("add menu fail:", err)
		respond.RediErr("menu", common.AlertFail, common.AlertSaveFail, c)
		return
	}
	system.SetMenu(common.GetTokenByCookie(c))
	respond.RediSuccess("/menu", c)
}

// ChangeMenuSort 改变菜单排序
func ChangeMenuSort(c *gin.Context) {
	id := c.Query("id")
	parentid := c.Query("parentid")
	updown := c.Query("updown")
	var u bool
	if updown == "up" {
		u = true
	} else {
		u = false
	}
	rs := system.ChangeSort(id, parentid, u)
	if !rs {
		respond.RediErr("menu", common.AlertFail, common.AlertSaveFail, c)
		return
	}
	system.SetMenu(common.GetTokenByCookie(c))
	respond.RediSuccess("/menu", c)
}
