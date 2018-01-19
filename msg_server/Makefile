all: msg_server

msg_server:
	go build -v -gcflags '-N' -o app_im_msg_server_d

clean:
	go clean
	rm -rf app_im_msg_server_d
