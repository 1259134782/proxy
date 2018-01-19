package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	simplejson "github.com/bitly/go-simplejson"
	log "github.com/cihub/seelog"
)

var msg_server_http_addr = "http://localhost:5060"
var proxy_regist_url = "/proxy/register"

func postToMsg(http_addr, url string, data []byte) (err error) {
	cli := httplib.Post(http_addr + url)
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

func HeatToMsgServer() {
	this_proxy_ip := beego.AppConfig.String("register_to_msg_server_add")
	this_proxy_port := ""
	if port_, err := beego.AppConfig.Int("httpport"); err != nil {
		log.Errorf("configure error no found proxy port")
		return
	} else {
		this_proxy_port = strconv.Itoa(port_)
	}
	post_data := make(map[string]string)
	post_data["ip"] = this_proxy_ip
	post_data["port"] = this_proxy_port

	json_post_data, _ := json.Marshal(post_data)
	for {
		select {
		case <-time.After(1 * time.Second):
			if err := postToMsg(msg_server_http_addr, proxy_regist_url, json_post_data); err != nil {
				log.Errorf("register error %v", err)
			}
		}
	}
}
