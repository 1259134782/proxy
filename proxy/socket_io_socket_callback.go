package main

import (
	"app_im/protocol"
	"encoding/base64"
	"encoding/json"
	"fmt"
	_ "reflect"

	log "github.com/cihub/seelog"
)

const (
	SOCKET_IO_SOCKET_MSG_DISPATCH_RECIVE_MSG = "dispatch_recive_msg"
	SOCKET_IO_SOCKET_MSG_DISPATCH_SEND_MSG   = "push"
	SOCKET_IO_SOCKET_MSG_PUSH_ID             = "pushId"
	SOCKET_IO_SOCKET_MSG_CHAT_MSG            = "chat_message"
	SOCKET_IO_SOCKET_MSG_DISCONNECTION       = "disconnection"
)

func (this *SocketIoSocket) InitCallBack() {
	this.RegistEvent(SOCKET_IO_SOCKET_MSG_DISPATCH_RECIVE_MSG, this.DispatchReciveMsg)
	this.RegistEvent(SOCKET_IO_SOCKET_MSG_DISPATCH_SEND_MSG, this.DispatchSendMsg)
	this.RegistEvent(SOCKET_IO_SOCKET_MSG_PUSH_ID, this.PushIdMsg)
	this.RegistEvent(SOCKET_IO_SOCKET_MSG_CHAT_MSG, this.ChatMgs)
	this.RegistEvent(SOCKET_IO_SOCKET_MSG_DISCONNECTION, this.Disconnection)
}

func (this *SocketIoSocket) DispatchReciveMsg(data interface{}) {
	// TODO logic
	//TODO res

	this.Send(SOCKET_IO_SOCKET_MSG_DISPATCH_RECIVE_MSG, data)
}

func (this *SocketIoSocket) DispatchSendMsg(data interface{}) {
	// TODO logic
	// 下行消息 目前没有此类需求
	//TODO res
	this.Send(SOCKET_IO_SOCKET_MSG_DISPATCH_SEND_MSG, data)
}

func (this *SocketIoSocket) BroadcastTo(room string, data interface{}) {
	// TODO logic
	// 广播消息
	//TODO res
	this.Broadcast(room, SOCKET_IO_SOCKET_MSG_DISPATCH_SEND_MSG, data)
}

func (this *SocketIoSocket) PushIdMsg(general interface{}) {
	// TODO logic
	data := general.(map[string]interface{})
	if id, ok := data["id"].(string); ok {
		this.PushId = id
		add_to_service_conns := OperaterDataConnMapProto{
			Id:       id,
			Socket:   this,
			Operater: ADD_CONN_TO_MAP,
		}
		this.GetSocketIoService().DispatchConnsOperator(add_to_service_conns)
		log.Infof("PushId id %v . general %v", id, general)
	}

	//TODO res
	// logic
	pk := protocol.PushIdProto{"vxh5v5KuKSy1234", "123456"}
	err := this.Send(SOCKET_IO_SOCKET_MSG_PUSH_ID, pk)
	if err != nil {
		log.Errorf("Send pushId to client faild. err:%v", err)
	}
}

func (this *SocketIoSocket) ChatMgs(general interface{}) {
	m := general.(map[string]interface{})
	fmt.Println("********************", m["message"], m["nickName"], "******************")
	message := protocol.MessageProto{}
	message.Type = "chat_message"

	if msg, ok := m["message"].(string); ok {
		message.Msg = msg
	}

	if nickName, ok := m["nickName"].(string); ok {
		message.NickName = nickName
	}

	data, _ := json.Marshal(message)
	var res protocol.ChatResProto
	res.Topic = "chatRoom"
	res.Data = base64.StdEncoding.EncodeToString(data)
	// server.BroadcastTo("chatRoom", "push", res)
	this.GetSocketIoService().GetServer().BroadcastTo("chatRoom", "push", res)
	// err := so.Emit("push", res)
	// if err != nil {
	// fmt.Println("###########", err, "##############")
	// }
	fmt.Println("********************end******************")
}

func (this *SocketIoSocket) Disconnection(data interface{}) {
	log.Infof("%v", "Disconnection")
	del_from_service_conns := OperaterDataConnMapProto{
		Id:       this.PushId,
		Socket:   this,
		Operater: DEL_CONN_TO_MAP,
	}

	this.GetSocketIoService().DispatchConnsOperator(del_from_service_conns)
}
