package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
)

type web1handler struct {
}

// 获取真实ip
func (web1handler) GetRealIp(req *http.Request) string {
	ips := req.Header.Get("x-forwarded-for")
	if ips != "" {
		ipList := strings.Split(ips, ",")
		if len(ipList) > 0 && ipList[0] != "" {
			return ipList[0]
		}
	}
	return req.RemoteAddr
}
func (web1 web1handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	//auth := req.Header.Get("Authorization")
	//if auth == "" {
	//	rw.Header().Set("WWW-Authenticate", `Basic realm="请输入用户名和密码"`)
	//	rw.WriteHeader(http.StatusUnauthorized)
	//	return
	//}
	//str := strings.Split(auth, " ")
	//if len(str) == 2 && str[0] == "Basic" {
	//	result, err := base64.StdEncoding.DecodeString(str[1])
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	if string(result) == "manbao:123456" {
	//		rw.Write([]byte(fmt.Sprintf("<h1>东皇太一，来自:%s</h1>", web1.GetRealIp(req))))
	//		fmt.Println(web1.GetRealIp(req))
	//		return
	//	}
	//}
	//rw.Write([]byte("用户名或密码错误"))
	rw.Write([]byte("web1"))
}

type web2handler struct {
}

func (web2handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("web2"))
}

func main() {
	c := make(chan os.Signal)
	go (func() {
		http.ListenAndServe(":9091", web1handler{})
	})()

	go (func() {
		http.ListenAndServe(":9092", web2handler{})
	})()

	signal.Notify(c, os.Interrupt)
	s := <-c
	fmt.Println(s)
}
