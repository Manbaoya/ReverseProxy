package util

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

// 拷贝头
func CloneHead(target *http.Header, origination http.Header) {
	for k, v := range origination {
		target.Set(k, v[0])
	}

}
func RequestUrl1(rw http.ResponseWriter, req *http.Request, url string) {
	newreq, err := http.NewRequest(req.Method, "http://localhost:9091", req.Body)
	if err != nil {
		rw.Write([]byte("获取页面失败"))
		return
	}
	tp := &http.Transport{
		DialContext: (&net.Dialer{ //连接超时
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ResponseHeaderTimeout: 20 * time.Second, //响应超时
	}
	CloneHead(&newreq.Header, req.Header)
	req.Header.Add("x-forwarded-for", req.RemoteAddr)
	newResponse, err := tp.RoundTrip(newreq)
	if err != nil {
		fmt.Println(err)
		return
	}
	getHeader := rw.Header()
	CloneHead(&getHeader, newResponse.Header) //拷贝头
	rw.WriteHeader(newResponse.StatusCode)    //写入状态码
	defer newResponse.Body.Close()
	response, _ := io.ReadAll(newResponse.Body)
	rw.Write(response)
	return
}
