package common

import (
	"github.com/JimYJ/easysql/mysql"
	"github.com/patrickmn/go-cache"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

var (
	dbhost, dbname, user, pass string
	port                       int
	once                       sync.Once
	c                          *cache.Cache
	host                       string
	// TokenTimeOut 登陆超时
	TokenTimeOut = 6 * time.Hour
	// CacheTimeOut 缓存超时
	CacheTimeOut = 1 * time.Minute
	// AppID ...
	AppID = ""
	// BmsPath 后管系统前缀路径
	BmsPath = "/bms"
	// SecretKey ...
	SecretKey = ""
	sha1salt  = []byte("63d81bc836e86565a5e8668faf1863cbbcd5b392fba28d7d48b39c858b3e4e75")
	//LoginMaxLimit 每分钟登陆请求限制
	LoginMaxLimit = 30
	//LoginGap 计次时间间隔，单位是秒
	LoginGap = 60
)

// 错误信息
var (
	//Err401 认证错误
	Err401 = "Authentication error!"
	//Err401login 登录失败
	Err401login = "User or Pass error!"
	//Err401captcha 验证码错误
	Err401captcha = "captcha is error!"
	//Err401SmsCode 登录失败
	Err401SmsCode = "Sms code is error!"
	//Err402Phone 手机号不正确
	Err402Phone = "Phone number is error!"
	//Err403PhoneIsBind 手机号已绑定
	Err403PhoneIsBind = "Phone is bind!"
	//Err403Unreg 手机号未注册
	Err403Unreg = "Phone is unregistered!"
	//Err403Unreg 手机号未注册
	Err406Unexpected = "request Unexpected!"
	//Err429Frequent 请求过于频繁
	Err429Frequent = "Request too Frequent!"
	//Err500DBrequest 数据库请求错误
	Err500DBrequest = "Database request error!"
	//Err500CannotGetUID 数据库请求错误
	Err500CannotGetUID = "Cannot get UID by token!"
	//Err502SMS 短信发送失败，检查短信平台账户密码
	Err502SMS = "Sms send fail!"
)

// http状态码
var (
	//HTTPAuthErr 认证错误
	HTTPAuthErr = 401
	//HTTPParamErr 请求参数错误
	HTTPParamErr = 402
	//HTTPForbiddenErr 拒绝请求
	HTTPForbiddenErr = 403
	//HTTPForbiddenErr 请求异常
	HTTPUnexpectedErr = 406
	//HTTPFrequentErr 请求过于频繁
	HTTPFrequentErr = 429
	//HTTPSystemErr 系统内部错误(程序，数据库等)
	HTTPSystemErr = 500
	//HTTPExternalErr 系统外部错误(第三方)
	HTTPExternalErr = 502
)

type config struct {
	Mysql     mysqlconf
	AppID     string
	SecretKey string
	Host      string
}

type mysqlconf struct {
	Host   string
	Port   int
	DBname string
	User   string
	Pass   string
}

func (conf *config) getConfig() *config {
	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return nil
	}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return nil
	}
	return conf
}

// GetConfig 获取配置文件
func GetConfig() {
	var conf config
	conf.getConfig()
	dbhost = conf.Mysql.Host
	dbname = conf.Mysql.DBname
	port = conf.Mysql.Port
	user = conf.Mysql.User
	pass = conf.Mysql.Pass
	AppID = conf.AppID
	SecretKey = conf.SecretKey
	host = conf.Host
}

// InitMysql 初始化mysql参数
func InitMysql() {
	mysql.Init(dbhost, port, dbname, user, pass, "utf8mb4", 1000, 1000)
}

// GetMysqlConn 获取mysql连接
func GetMysqlConn() *mysql.MysqlDB {
	mysqlConn, err := mysql.GetMysqlConn()
	if err != nil {
		log.Panicln("mysql conn error:", err)
		return nil
	}
	return mysqlConn
}

// GetCache 获得缓存对象
func GetCache() *cache.Cache {
	once.Do(func() {
		c = cache.New(TokenTimeOut, 10*time.Minute)
	})
	return c
}
