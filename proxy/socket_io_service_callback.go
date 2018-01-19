package main

import (
	_ "reflect"
	"github.com/googollee/go-socket.io"
)

const (
	SOCKET_IO_SERVICE_MSG_CONNECTION = "connection"
	SOCKET_IO_SERVICE_MSG_DISCONNECT = "disconnect"
	SOCKET_IO_SERVICE_MSG_ERROR      = "error"
)

func (this *SocketIoService) InitCallBack() {
	this.RegistEvent(SOCKET_IO_SERVICE_MSG_CONNECTION, this.Connection)
	this.RegistEvent(SOCKET_IO_SERVICE_MSG_DISCONNECT, this.Disconnection)
	this.RegistEvent(SOCKET_IO_SERVICE_MSG_ERROR, this.Error)
}

func (this *SocketIoService) Connection(so socketio.Socket) {
	this.connection_chan_pool <- so
}

func (this *SocketIoService) Disconnection(so socketio.Socket) {
	this.disconnection_chan_pool <- so
}

func (this *SocketIoService) Error(so socketio.Socket) {
	this.error_chan_pool <- so
}
