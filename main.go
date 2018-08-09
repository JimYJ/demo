package main

import (
	"666sites/common"
	"666sites/route"
	// "log"
)

func main() {
	inits()
}

func inits() {
	// log.SetFlags(log.Lshortfile | log.LstdFlags)
	common.GetConfig()
	common.InitMysql()
	common.GetMysqlConn()
	route.Web()
}
