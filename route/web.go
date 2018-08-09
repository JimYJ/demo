package route

import (
	"666sites/app"
	"666sites/app/auth"
	"666sites/app/sys"
	"666sites/common"
	"666sites/route/middleware"
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
	menu := bms.Group("/menu", middleware.CheckUserMenu("/menu"))
	menu.GET("/", sys.Menu)
	// 菜单删除
	menu.GET("/del", sys.DelMenu)
	// 新增菜单
	menu.POST("/add", sys.AddMenu)
	// 编辑菜单
	menu.POST("/edit", sys.EditMenu)
	// 菜单排序
	menu.GET("/sort", sys.ChangeMenuSort)
	// 获取全部菜单(用于角色权限管理)
	menu.GET("/list", sys.ChangeMenuSort)

	// 后台角色管理页
	role := bms.Group("/role", middleware.CheckUserMenu("/role"))
	role.GET("/", sys.Role)
	// 后台角色删除
	role.GET("/del", sys.DelRole)
	// 新增后台角色
	role.POST("/add", sys.AddRole)
	// 编辑后台角色
	role.POST("/edit", sys.EditRole)
	// 获取管理用户岗位
	role.GET("/menulist", sys.GetRoleMenu)
	// 管理用户岗位
	role.POST("/bindmenu", sys.RoleBindMenu)

	// 后台管理用户管理页
	admin := bms.Group("/admin", middleware.CheckUserMenu("/admin"))
	admin.GET("/", sys.AdminUser)
	// 后台管理用户删除
	admin.GET("/del", sys.DelAdminUser)
	// 新增后台管理用户
	admin.POST("/add", sys.AddAdminUser)
	// 编辑后台管理用户
	admin.POST("/edit", sys.EditAdminUser)
	// 获取管理用户岗位
	admin.GET("/rolelist", sys.GetAdminRole)
	// 管理用户岗位
	admin.POST("/bindrole", sys.AdminBindRole)

	// router.RunTLS("127.0.0.1:443", sslcert, sslkey)
	router.Run(":80")
}
