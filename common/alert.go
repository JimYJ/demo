package common

import (
	"strconv"
)

//alert title
const (
	AlertWarning = iota
	AlertError
	AlertFail
)

// alert Content
const (
	AlertParamsError = iota
	AlertSaveFail
	AlertDelFail
	AlertGetDateFail
	AlertLoginFail
	AlertUserError
	AlertPassError
	AlertDBFail
	AlertCheckTokenError
	AlertPathVisitError
)

var (
	alertTitle   = []string{"警告", "错误", "失败"}
	alertContent = []string{
		"提交参数错误!",
		"提交/保存失败,具体请查看日志.",
		"删除失败,具体请查看日志.",
		"获取数据失败,或无数据记录,具体请查看日志.",
		"登录失败,请检查账户和密码.",
		"用户名只允许4-12位数字+大小写字母+下划线组成",
		"密码必须6位以上",
		"数据库查询失败,请联系管理员查看日志",
		"登录失效,登陆已超时或账户已经在其他地方登录",
		"你无权访问本页面,请重新登录",
	}
)

// GetAlertMsg 获取消息
func GetAlertMsg(t, c string) (string, string) {
	if t == "" || c == "" {
		return "", ""
	}
	ti, err := strconv.Atoi(t)
	ci, err2 := strconv.Atoi(c)
	if err != nil || err2 != nil {
		return "", ""
	}
	if ti >= len(alertTitle) || ci >= len(alertContent) {
		return "", ""
	}
	return alertTitle[ti], alertContent[ci]
}
