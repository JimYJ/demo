package route

import (
	"bms/app"
	"bms/app/auth"
	"bms/app/sys"
	"bms/common"
	"bms/route/middleware"
	"github.com/gin-gonic/gin"
)

var (
	api *gin.RouterGroup
)

// Web 路由
func Web() {
	// gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Static("/bms/assets", "./statics/assets")
	router.LoadHTMLGlob("statics/html/*")
	// router.Use(middleware.Cors())

	api = router.Group("/api")

	// ----------------------- web 路径 ---------------------------
	bms := router.Group(common.BmsPath)
	bmsLogin := router.Group(common.BmsPath)
	bms.Use(middleware.TokenAuth())

	// 登录页
	bmsLogin.GET("/login", app.Login)
	// 验证登录
	bmsLogin.POST("/checklogin", auth.Login)
	// 登出
	bms.GET("/logout", auth.Logout)
	// 首页
	bms.GET("/", middleware.CheckUserMenu("/"), app.Index)

	// 菜单管理页
	bms.GET("/menu", middleware.CheckUserMenu("/menu"), sys.Menu)
	// 菜单删除
	bms.GET("/delmenu", sys.DelMenu)
	// 新增菜单
	bms.POST("/addmenu", sys.AddMenu)
	// 编辑菜单
	bms.POST("/editmenu", sys.EditMenu)
	// 菜单排序
	bms.GET("/menusort", sys.ChangeMenuSort)
	// 获取全部菜单(用于角色权限管理)
	bms.GET("/menulist", sys.ChangeMenuSort)

	// 后台角色管理页
	bms.GET("/role", middleware.CheckUserMenu("/role"), sys.Role)
	// 后台角色删除
	bms.GET("/delrole", sys.DelRole)
	// 新增后台角色
	bms.POST("/addrole", sys.AddRole)
	// 编辑后台角色
	bms.POST("/editrole", sys.EditRole)
	// 获取管理用户岗位
	bms.GET("/rolemenulist", sys.GetRoleMenu)
	// 管理用户岗位
	bms.POST("/rolebindmenu", sys.RoleBindMenu)

	// 后台管理用户管理页
	bms.GET("/admin", middleware.CheckUserMenu("/admin"), sys.AdminUser)
	// 后台管理用户删除
	bms.GET("/deladmin", sys.DelAdminUser)
	// 新增后台管理用户
	bms.POST("/addadmin", sys.AddAdminUser)
	// 编辑后台管理用户
	bms.POST("/editadmin", sys.EditAdminUser)
	// 获取管理用户岗位
	bms.GET("/adminrolelist", sys.GetAdminRole)
	// 管理用户岗位
	bms.POST("/adminbindrole", sys.AdminBindRole)

	// router.RunTLS("127.0.0.1:443", sslcert, sslkey)
	router.Run(":80")
}
