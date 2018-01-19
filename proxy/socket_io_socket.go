package main

import (
	_ "reflect"

	"github.com/googollee/go-socket.io"
)

type SocketIoSocketHandle func(general interface{})

type SocketIoSocket struct {
	parent_socket_io_service *SocketIoService
	socket_io_socket         socketio.Socket
	PushId                   string
}

func NewSocketIoSocket(parent_socket_io_service_ *SocketIoService, so socketio.Socket) (socket_ *SocketIoSocket) {

	socket_ = &SocketIoSocket{
		parent_socket_io_service: parent_socket_io_service_,
		socket_io_socket:         so,
	}

	return

}

func (this *SocketIoSocket) GetSocketIoService() *SocketIoService {
	return this.parent_socket_io_service
}

func (this *SocketIoSocket) SetSocketEvent(event_name string, handle SocketIoSocketHandle) (err error) {
	this.socket_io_socket.On(event_name, handle)
	return
}

func (this *SocketIoSocket) RegistEvent(event_name string, handle SocketIoSocketHandle) (err error) {
	err = this.SetSocketEvent(event_name, handle)
	return
}

func (this *SocketIoSocket) Send(event_name string, data interface{}) (err error) {
	err = this.socket_io_socket.Emit(event_name, data)
	return
}
func (this *SocketIoSocket) Broadcast(room string, event_name string, data interface{}) (err error) {
	err = this.socket_io_socket.BroadcastTo(room, event_name, data)
	return
}
func (this *SocketIoSocket) Init() {
	this.InitCallBack()
}
