package main

import (
	"app_im/im_utils"
	"app_im/models"
	"app_im/protocol"
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	simplejson "github.com/bitly/go-simplejson"
	log "github.com/cihub/seelog"
)

type MgsController struct {
	// token.BaseController
	beego.Controller
}

var ProxyRegisters map[string]string

// // proxy 的管理
// var proxy_chan_queue chan

// //用户的管理
// var user_chan_queue chan

func init() {
	ProxyRegisters = make(map[string]string)
}

func (this *MgsController) ParsePara(req interface{}) error {
	//log.Infof("req:%v", string(this.Ctx.Input.RequestBody))

	if json_byte_arr, err := base64.StdEncoding.DecodeString(this.GetString("data")); err != nil {
		log.Errorf("base64 decode error %v", err)
		return err
	} else {
		if err := json.Unmarshal(json_byte_arr, &req); err != nil {
			log.Errorf("Unmarshal req error %v", err)
			return err
		}
	}
	return nil

}

func (this *MgsController) postToProxy(proxy_addr, url string, data []byte) (err error) {
	cli := httplib.Post(proxy_addr + url)
	// 由于是转发，不需要encode
	cli.Param("data", base64.StdEncoding.EncodeToString(data))

	res, err := cli.String()
	if err != nil {
		return err
	}

	js, err := simplejson.NewJson([]byte(res))
	rescode, _ := js.Get("code").String()
	if rescode != "success" {
		return errors.New("Invalid Rescode " + rescode)
	}
	return nil
}

func (this *MgsController) MsgPush() {
	var res protocol.PushResJsonProto
	var req protocol.ReqPushJsonProto

	defer func() {
		this.Data["json"] = res
		this.ServeJSON()
	}()

	// log.Infof("req:%v", string(this.Ctx.Input.RequestBody))
	res.Code = "success"

	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &req); err != nil {
		res.Code = "-1"
		return
	}

	if proxy, err := im_utils.GetProxyByPushId(req.PushId); err != nil {
		log.Errorf("redis not found proxy by pushid ")
	} else {
		this.postToProxy(proxy, "/msg/push", []byte(req.Data))
		log.Infof("pushId:%v, data:%v ", req.PushId, req.Data)
	}

	return
}

func (this *MgsController) MsgNotification() {
	var res protocol.NotificationResJsonProto
	var req protocol.ReqNotificationJsonProto

	defer func() {
		this.Data["json"] = res
		this.ServeJSON()
	}()

	log.Infof("req:%v", string(this.Ctx.Input.RequestBody))
	res.Code = "success"

	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &req); err != nil {
		res.Code = "-1"
		return
	}
	log.Infof("pushId:%v, data:%v ", req.PushId, req.Data)

	if proxy, err := im_utils.GetProxyByPushId(req.PushId); err != nil {
		log.Errorf("redis not found proxy by pushid ")
	} else {
		this.postToProxy(proxy, "/msg/notify", []byte(req.Data))
	}

	log.Infof("NotificationId:%v, data:%v ", req.PushId, req.Data)
	return
}

func (this *MgsController) ProxyUserConnect() {
	var res protocol.ProxyUserConnectionResJsonProto
	var req protocol.ReqProxyUserConnectionJsonProto

	defer func() {
		this.Data["json"] = res
		this.ServeJSON()
	}()

	res.Code = "success"

	if err := this.ParsePara(&req); err != nil {
		log.Errorf("ParsePara err %v", err)
	}

	im_utils.SetProxyByPushId(req.Pushid, req.Proxy)

	log.Infof("add user pushid :%v, proxy :%v ", req.Pushid, req.Proxy)
	return
}

func (this *MgsController) ProxyUserDisconnect() {
	var res protocol.ProxyUserDisConnectionResJsonProto
	var req protocol.ReqProxyUserDisConnectionJsonProto

	defer func() {
		this.Data["json"] = res
		this.ServeJSON()
	}()

	res.Code = "success"

	if err := this.ParsePara(&req); err != nil {
		log.Errorf("ParsePara err %v", err)
	}

	default_str := models.MYSQL_FIELD_DEFAULT_VAL_OF_STRING
	if req.Proxy != default_str {
		log.Infof("user drop line , but the proxy is not ok")
	}

	im_utils.SetProxyByPushId(req.PushId, default_str)

	log.Infof("del user pushid :%v, proxy :%v ", req.PushId, req.Proxy)
	return
}

func (this *MgsController) ProxyRegister() {
	var res protocol.ProxyRegisterResJsonProto
	var req protocol.ReqProxyRegisterJsonProto

	defer func() {
		this.Data["json"] = res
		this.ServeJSON()
	}()

	res.Code = "success"

	if err := this.ParsePara(&req); err != nil {
		log.Errorf("ParsePara err %v", err)
	}

	proxy_ip_port := req.ProxyIP + ":" + req.ProxyPort
	if _, ok := ProxyRegisters[proxy_ip_port]; ok {
		// 已经注册过了
		// 需不需要其他逻辑，待定
		// 例如是否强制更新
	} else {
		ProxyRegisters[proxy_ip_port] = proxy_ip_port
		log.Infof("has new proxy in come : %v", proxy_ip_port)
		log.Infof("ProxyRegister IP:%v, Port:%v ", req.ProxyIP, req.ProxyPort)
	}

	return
}
