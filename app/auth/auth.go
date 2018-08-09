package auth

import (
	"666sites/app/respond"
	"666sites/common"
	"666sites/db/system"
	"666sites/service/log"
	"github.com/gin-gonic/gin"
	"regexp"
	"strconv"
	"time"
)

// Login 登陆验证
func Login(c *gin.Context) {
	user := c.PostForm("user")
	pass := c.PostForm("password")
	matchUser, _ := regexp.MatchString("^[0-9a-zA-Z_]{4,12}$", user)
	if !matchUser {
		respond.RediErr("login", common.AlertFail, common.AlertLoginFail, c)
		return
	}
	if len(pass) < 6 {
		respond.RediErr("login", common.AlertFail, common.AlertLoginFail, c)
		return
	}
	pass = common.SHA1(pass)
	uinfo, err := system.CheckPass(user, pass)
	id := uinfo["id"]
	if err == 500 {
		respond.RediErr("login", common.AlertFail, common.AlertDBFail, c)
		return
	} else if err == 401 {
		log.Println(3)
		respond.RediErr("login", common.AlertFail, common.AlertLoginFail, c)
		return
	}
	t := time.Now().UnixNano()
	timestamp := strconv.FormatInt(t, 10)
	ip := []byte(c.ClientIP())
	uid := []byte(strconv.FormatInt(id.(int64), 10))
	token := common.CreateToken(ip, uid, []byte(timestamp))
	cache := common.GetCache()
	tokeninfo := common.GetTokenCache(string(uid), timestamp, user)
	cache.Set(token, tokeninfo, common.TokenTimeOut)
	common.SingleLogin(token)
	common.SetCookie(c, "c", token)
	system.SetMenu(token)
	respond.RediSuccess("/", c)
}

// Logout 登出
func Logout(c *gin.Context) {
	var token string
	if cookie, err := c.Request.Cookie("c"); err == nil {
		token = cookie.Value
	} else {
		respond.RediSuccess("/login", c)
		return
	}
	cache := common.GetCache()
	cache.Delete(token)
	system.SetMenu(token)
	respond.RediSuccess("/login", c)
}
