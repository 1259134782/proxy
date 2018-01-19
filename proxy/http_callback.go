package main

import (
	"app_im/protocol"
	"encoding/base64"
	"encoding/json"

	"github.com/astaxie/beego"
	log "github.com/cihub/seelog"
)

type ProxyController struct {
	// token.BaseController
	beego.Controller
}

func (this *ProxyController) ParsePara(req interface{}) error {
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

func (this *ProxyController) ProxyPush() {
	var res protocol.PushResJsonProto
	//var req protocol.ReqPushJsonProto
	defer func() {
		this.Data["json"] = res
		this.ServeJSON()
	}()

	log.Infof("req:%v", string(this.Ctx.Input.RequestBody))
	//	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &req); err != nil {
	//		res.Code = "-1"
	//		return
	//	}

	var pushIds []string
	if err := json.Unmarshal([]byte(this.GetString("pushId")), &pushIds); err != nil {
		log.Infof("*****err:%v******", err)
		res.Code = "-1"
		return
	}

	data := this.GetString("data")
	byte1, _ := base64.StdEncoding.DecodeString(data)
	var general interface{}
	if err := json.Unmarshal(byte1, &general); err != nil {
		log.Infof("*****err:%v******", err)
		res.Code = "-1"
		return
	} else {
		log.Infof("*****general:%v******", general)
	}

	for _, pushId := range pushIds {
		log.Infof("pushId:%v, data:%v ", pushId, data)
		send_msg_operater_to_client := OperaterDataConnMapProto{
			Id:       pushId,
			Socket:   nil,
			Operater: SEND_MESSAGE_TO_CLIENT,
			Data:     general,
		}

		GetGloableIoService().DispatchConnsOperator(send_msg_operater_to_client)
	}

	return
}

func (this *ProxyController) ProxyNotify() {
	var res protocol.NotifyResJsonProto
	var req protocol.ReqNotifyJsonProto

	defer func() {
		this.Data["json"] = res
		this.ServeJSON()
	}()

	log.Infof("req:%v", string(this.Ctx.Input.RequestBody))

	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &req); err != nil {
		res.Code = "-1"
		return
	}

	log.Infof("pushId:%v, data:%v ", req.PushId, req.Data)
	return
}
