package respond

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// 后台路径
var (
	BmsPath = "/bms"
)

// RedirectBody 响应结构体
type RedirectBody struct {
	path                     string
	alerttitle, alertcontent int
	c                        *gin.Context
}

// Err 错误结果转跳
func (r *RedirectBody) Err() {
	url := fmt.Sprintf("%s%s?t=%d&c=%d", BmsPath, r.path, r.alerttitle, r.alertcontent)
	r.c.Redirect(302, url)
	r.c.Abort()
}

// Success 正常结果转跳
func (r *RedirectBody) Success() {
	url := fmt.Sprintf("%s%s", BmsPath, r.path)
	r.c.Redirect(302, url)
}

// RediErr 错误结果转跳
func RediErr(path string, alerttitle, alertcontent int, c *gin.Context) {
	if string([]rune(path)[0]) != "/" {
		path = fmt.Sprintf("/%s", path)
	}
	r := &RedirectBody{path, alerttitle, alertcontent, c}
	r.Err()
}

// RediSuccess 成功结果转跳
func RediSuccess(path string, c *gin.Context) {
	r := &RedirectBody{path, 0, 0, c}
	r.Success()
}
