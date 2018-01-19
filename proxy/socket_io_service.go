package main

import (
	"errors"
	"fmt"
	_ "reflect"

	log "github.com/cihub/seelog"
	"github.com/googollee/go-socket.io"
)

const (
	CHAN_POOL_MAX = 4096
)

type ConnMapOperater int

const (
	ADD_CONN_TO_MAP        ConnMapOperater = 1
	DEL_CONN_TO_MAP                        = 2
	SEND_MESSAGE_TO_CLIENT                 = 4
)

type OperaterDataConnMapProto struct {
	Id       string
	Socket   *SocketIoSocket
	Operater ConnMapOperater
	Data     interface{}
}

type SocketIoService struct {
	socket_io_server        *socketio.Server
	socket_io_conns         map[string]*SocketIoSocket
	event_map               map[string]SocketIoServiceHandle
	connection_chan_pool    chan socketio.Socket
	disconnection_chan_pool chan socketio.Socket
	error_chan_pool         chan socketio.Socket
	conns_chan_pool         chan OperaterDataConnMapProto
}

type SocketIoServiceHandle func(so socketio.Socket)

func NewSocketIoService() (service *SocketIoService) {
	server_, err := socketio.NewServer(nil)
	if err != nil {
		log.Errorf("Init socket.io server faild. err:%v", err)
	}

	server_.SetMaxConnection(1024 * 10 * 10 * 10)

	service = &SocketIoService{
		socket_io_server:        server_,
		socket_io_conns:         make(map[string]*SocketIoSocket),
		event_map:               make(map[string]SocketIoServiceHandle),
		connection_chan_pool:    make(chan socketio.Socket, CHAN_POOL_MAX),
		disconnection_chan_pool: make(chan socketio.Socket, CHAN_POOL_MAX),
		error_chan_pool:         make(chan socketio.Socket, CHAN_POOL_MAX),
		conns_chan_pool:         make(chan OperaterDataConnMapProto, CHAN_POOL_MAX),
	}

	//初始化工作
	service.Init()
	//TODO 开始工作的位置可能需要调整
	service.loop()
	return
}

func (this *SocketIoService) GetServer() *socketio.Server {
	return this.socket_io_server
}

func (this *SocketIoService) toDoConnection(so socketio.Socket) {
	socket_io_socket_ := NewSocketIoSocket(this, so)
	socket_io_socket_.Init()
	//log.Infof("on connection")
	so.Join("chatRoom")
}

func (this *SocketIoService) toDoDisconnection(so socketio.Socket) {
}

func (this *SocketIoService) toDoError(so socketio.Socket) {
}

func (this *SocketIoService) SetServiceEvent(event_name string, handle SocketIoServiceHandle) (err error) {
	if _, ok := this.event_map[event_name]; ok {
		err_str := fmt.Sprintf("had alread regist Event %v", event_name)
		err = errors.New(err_str)
	} else {
		this.event_map[event_name] = handle
	}

	this.socket_io_server.On(event_name, handle)

	return
}

func (this *SocketIoService) GetServiceEvent(event_name string) (handle SocketIoServiceHandle, err error) {
	if handle_, ok := this.event_map[event_name]; ok {
		handle = handle_
	} else {
		err_str := fmt.Sprintf("had not regist Event %v", event_name)
		err = errors.New(err_str)
		handle = nil
	}

	return
}

func (this *SocketIoService) RegistEvent(event_name string, handle SocketIoServiceHandle) (err error) {
	err = this.SetServiceEvent(event_name, handle)
	return
}

func (this *SocketIoService) Init() {
	this.InitCallBack()
}

func (this *SocketIoService) DispatchConnsOperator(val OperaterDataConnMapProto) {
	this.conns_chan_pool <- val
}

func (this *SocketIoService) toDoConnsOperaterDispatch(val OperaterDataConnMapProto) {
	id := val.Id
	so := val.Socket
	switch val.Operater {
	case ADD_CONN_TO_MAP:
		this.socket_io_conns[id] = so
	case DEL_CONN_TO_MAP:
		delete(this.socket_io_conns, id)
	case SEND_MESSAGE_TO_CLIENT:
		socket_ := this.socket_io_conns[id]
		socket_.DispatchSendMsg(val.Data)

	}

	for k, v := range this.socket_io_conns {
		log.Infof("******** k:%v,v:%v**********", k, v)
	}
}

func (this *SocketIoService) loop() {
	go func() {
		for {
			select {
			case connection_chan_ := <-this.connection_chan_pool:
				this.toDoConnection(connection_chan_)
			case disconnection_chan_ := <-this.disconnection_chan_pool:
				this.toDoDisconnection(disconnection_chan_)
			case err_chan_ := <-this.error_chan_pool:
				this.toDoError(err_chan_)
			case conns_operater_chan_ := <-this.conns_chan_pool:
				this.toDoConnsOperaterDispatch(conns_operater_chan_)
			}
		}
	}()

	return
}
