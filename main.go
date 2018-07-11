package main

import (
	"bms/common"
	"bms/route"
	"log"
)

func main() {
	inits()
}

func inits() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	common.GetConfig()
	common.InitMysql()
	route.Web()
}
