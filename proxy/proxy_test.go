package main

import (
	"net"
	"testing"
)

func Test_Addr(t *testing.T) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, addr := range addrs {
		t.Logf(addr.String())
		t.Logf(addr.Network())
	}

}
