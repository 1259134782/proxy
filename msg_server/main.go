package main

import (
	_ "reflect"
	"vipcomm"
	"vipcomm/mylog"

	"github.com/astaxie/beego"
	log "github.com/cihub/seelog"
)

func main() {
	// 日志
	mylog.InitLog()
	defer mylog.FlushLog()

	// 初始化服务
	if err := vipcomm.InitServer(); err != nil {
		log.Errorf("Init server faild. err:%v", err)
		return
	}

	// 供业务进程调用的接口
	// pushClientUrl = addr + "/api/push"
	// sysNotifyUrl = addr + "/api/notification"

	// json 格式  有两个字段 pushId  及 data
	beego.Router("/api/push", &MgsController{}, "post:MsgPush")
	beego.Router("/api/notification", &MgsController{}, "post:MsgNotification")

	// 供proxy调用
	beego.Router("/proxy/register", &MgsController{}, "post:ProxyRegister")

	beego.Router("/proxy/user_connect", &MgsController{}, "post:ProxyUserConnect")
	beego.Router("/proxy/user_disconnect", &MgsController{}, "post:ProxyUserDisconnect")

	beego.Run()

}
