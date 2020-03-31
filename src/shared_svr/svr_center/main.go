package main

import (
	"common/console"
	"common/console/shutdown"
	"common/file"
	"common/tool/email"
	"conf"
	"flag"
	"gamelog"
	_ "generate_out/rpc/shared_svr/svr_center"
	"netConfig"
	"netConfig/meta"
	"shared_svr/svr_center/logic"
)

const kModuleName = "center"

func main() {
	var svrId int
	flag.IntVar(&svrId, "id", 1, "svrId")
	flag.Parse()

	//初始化日志系统
	gamelog.InitLogger(kModuleName)
	InitConf()

	//设置本节点meta信息
	meta.G_Local = meta.GetMeta(kModuleName, svrId)

	netConfig.RunNetSvr(false)
	logic.MainLoop()
}
func InitConf() {
	var metaCfg []meta.Meta
	file.G_Csv_Map = map[string]interface{}{
		"csv/conf_net.csv":      &metaCfg,
		"csv/conf_svr.csv":      &conf.SvrCsv,
		"csv/email/email.csv":   &email.G_EmailCsv,
		"csv/email/invalid.csv": &email.G_InvalidCsv,
	}
	file.LoadAllCsv()
	meta.InitConf(metaCfg)
	console.Init()
	console.RegShutdown(shutdown.Default)
}
