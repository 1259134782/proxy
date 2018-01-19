package main

import (
	_ "reflect"
	"net/http"
	"github.com/astaxie/beego"
	log "github.com/cihub/seelog"
	"vipcomm/mylog"
	_"vipcomm"
)

var GloableIoServiceInstance *SocketIoService

func init() {
	GloableIoServiceInstance = NewSocketIoService()
}

func GetGloableIoService() *SocketIoService {
	return GloableIoServiceInstance
}

func main() {
	// 日志
	  mylog.InitLog()
	  defer mylog.FlushLog()

	// 初始化服务
	/*
	if err := vipcomm.InitServer(); err != nil {
			log.Errorf("Init server faild. err:%v", err)
			return
	}*/
	socket_io_service := GetGloableIoService()

	http.Handle("/socket.io/", socket_io_service.GetServer())
	http.Handle("/", http.FileServer(http.Dir("./asset")))

	log.Infof("Serving at localhost:5050...")
	if err := http.ListenAndServe(":5050", nil); err != nil {
	log.Errorf("ListenAndServe happen error:%v", err)
	}

	go HeatToMsgServer()

	// socket.io 的处理
	beego.Handler("/socket.io/", socket_io_service.GetServer())

	//从 msg server 下行的消息
	beego.Router("/msg/push", &ProxyController{}, "post:ProxyPush")
	beego.Router("/msg/notify", &ProxyController{}, "post:ProxyNotify")
	beego.Run()

}
