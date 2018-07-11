package respond

import (
	"github.com/gin-gonic/gin"
)

// Body 响应结构体
type Body struct {
	code int
	msg  string
	c    *gin.Context
}

//Err 错误结果响应
func (r *Body) Err() {
	body := make(map[string]string)
	body["msg"] = r.msg
	r.c.JSON(r.code, body)
	r.c.Abort()
}

//Success 错误结果响应
func (r *Body) Success() {
	body := make(map[string]string)
	body["msg"] = r.msg
	r.c.JSON(r.code, body)
	r.c.Abort()
}

// Success 成功的响应
func Success(c *gin.Context) {
	resp := &Body{200, "success", c}
	resp.Success()
}

// Err 错误的响应
func Err(code int, msg string, c *gin.Context) {
	resp := &Body{code, msg, c}
	resp.Success()
}
