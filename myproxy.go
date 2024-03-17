package main

import (
	. "ReverseProxy/util"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ProxyHandler struct{}

func (p ProxyHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			rw.WriteHeader(500)
			fmt.Println(err)
		}
	}()

	//随机算法
	//server := LB.SelectByRand()

	//ip_hash算法
	//server := LB.SelectByIpHash(req.RemoteAddr)

	//加权随机算法
	//server := LB.SelectByWeightRand()

	//加权随机算法改良
	//server := LB.SelectByWeightRand2()

	//轮询算法
	//server := LB.RoundRobinByWeight()

	//加权轮询算法
	//server := LB.RoundRobinByWeight2()

	//平滑加权轮询算法
	server := LB.SmoothRoundRobinByWeight()
	target, _ := url.Parse(server.Host)
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ServeHTTP(rw, req)

	//for k, v := range util.ProxyConfig {
	//
	//	ok, err := regexp.MatchString(k, req.URL.Path)
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//	if ok {
	//		//util.RequestUrl1(rw, req, v)
	//		target, _ := url.Parse(v) //目标网站
	//		proxy := httputil.NewSingleHostReverseProxy(target)
	//		proxy.ServeHTTP(rw, req)
	//	}
	//
	//

	//}
}
func main() {
	http.ListenAndServe(":8080", ProxyHandler{})
}
