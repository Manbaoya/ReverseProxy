package util

import (
	"fmt"
	"hash/crc32"
	"math/rand"
	"sort"
	"time"
)

type HttpServer struct { //服务器类
	Host       string
	Weight     int
	NowWeight  int //当前权重
	FailWeight int
	Status     string
	FailSum    int
	SuccessSum int
}
type LoadBalance struct { //负载均衡类
	Servers ServerSlice
	Current int //当前服务坐标
}

//初始化

func NewHttpServer(host string, weight int) *HttpServer {
	return &HttpServer{Host: host, Weight: weight, NowWeight: weight}
}

func NewLoadBalance() *LoadBalance {
	return &LoadBalance{
		Servers: make([]*HttpServer, 0),
	}
}

// 添加服务

func (this *LoadBalance) AddServer(server *HttpServer) {
	this.Servers = append(this.Servers, server)
}

type ServerSlice []*HttpServer

// 排序

func (s ServerSlice) Less(i, j int) bool {
	return s[i].NowWeight > s[j].NowWeight
}
func (s ServerSlice) Len() int { return len(s) }

func (s ServerSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// 随机算法

func (this *LoadBalance) SelectByRand() *HttpServer {
	rand.NewSource(time.Now().UnixNano())
	index := rand.Intn(len(this.Servers))
	return this.Servers[index]
}

// 加权随机算法

func (this *LoadBalance) SelectByWeightRand() *HttpServer {
	rand.NewSource(time.Now().UnixNano())
	index := rand.Intn(len(ServerList))
	fmt.Println(ServerList)
	return this.Servers[ServerList[index]]
}

// 加权随机算法改良版

func (this *LoadBalance) SelectByWeightRand2() *HttpServer {
	var length = len(this.Servers)
	SumList := make([]int, length)
	sum := 0
	for i := 0; i < length; i++ {
		sum += this.Servers[i].Weight
		SumList[i] = sum
	}
	rand.NewSource(time.Now().UnixNano())
	index := rand.Intn(sum)
	for i := 0; i < length; i++ {

		if index < SumList[i] {
			return this.Servers[i]
		}
	}

	fmt.Println(SumList)
	return this.Servers[0]
}

//ip_hash算法

func (this *LoadBalance) SelectByIpHash(ip string) *HttpServer {
	index := int(crc32.ChecksumIEEE([]byte(ip))) % len(this.Servers)
	return this.Servers[index]
}

//轮询算法

func (this *LoadBalance) RoundRobin() *HttpServer {

	server := this.Servers[this.Current]
	this.Current = (this.Current + 1) % len(this.Servers)
	if server.Status == "Down" && this.IsAllDown() {
		return server
	}
	if server.Status == "Down" {
		return this.RoundRobin()
	}

	return server
}

// 加权轮询算法

func (this *LoadBalance) RoundRobinByWeight() *HttpServer {
	server := this.Servers[ServerList[this.Current]]
	this.Current = (this.Current + 1) % len(ServerList)

	return server
}

// 加权轮询算法改良版

func (this *LoadBalance) RoundRobinByWeight2() *HttpServer {
	sumList := make([]int, len(this.Servers))
	var sum = 0
	server := this.Servers[0]
	for i := 0; i < len(this.Servers); i++ {
		realWeight := this.Servers[i].Weight - this.Servers[i].FailWeight
		if realWeight == 0 {
			continue
		}
		sum += realWeight
		sumList[i] = sum
		if this.Current < sumList[i] {
			server = this.Servers[i]
			if this.Current == sum-1 && i != len(this.Servers)-1 {
				this.Current++
			} else {
				this.Current = (this.Current + 1) % sum

			}
			break
		}

	}
	fmt.Println(server.Host, server.FailWeight, server.Weight)
	return server
}

// 平滑加权轮询算法

func (this *LoadBalance) SmoothRoundRobinByWeight() *HttpServer {

	sort.Sort(this.Servers)
	server := this.Servers[0]
	fmt.Println(server.NowWeight)
	sumWeight := this.GetWeightSum()
	this.Servers[0].NowWeight -= sumWeight
	for _, v := range this.Servers {
		v.NowWeight += v.Weight - v.FailWeight
	}
	fmt.Println(server.Host, server.FailWeight, server.Weight)
	return server

}

// 实时获取真实权重和
func (this *LoadBalance) GetWeightSum() int {
	sum := 0
	for _, v := range this.Servers {
		realWeight := v.Weight - v.FailWeight
		if realWeight > 0 {
			sum += realWeight
		}
	}
	return sum
}

// 计时器

func (this *LoadBalance) CheckServers() {
	t := time.NewTicker(10 * time.Second)
	check := &HttpChecker{this.Servers, 10, 3}
	for {
		select {
		case <-t.C:
			{
				check.Check(5 * time.Second)

			}
		}
	}
}
func (this *LoadBalance) IsAllDown() bool {
	var sum = 0
	for _, v := range this.Servers {
		if v.Status == "Down" {
			sum++
		}
	}
	var length = len(this.Servers)
	if sum == length {
		return true
	}
	return false
}

var LB *LoadBalance
var ServerList []int
var SumWeight = 0

func init() {
	LB = NewLoadBalance()
	LB.AddServer(&HttpServer{Host: "http://localhost:9091", Weight: 5, Status: "UP"})
	LB.AddServer(&HttpServer{Host: "http://localhost:9092", Weight: 10, Status: "UP"})

	for k, v := range LB.Servers {

		SumWeight += v.Weight
		if v.Weight > 0 {
			for i := 0; i < v.Weight; i++ {
				ServerList = append(ServerList, k)
			}
		}
	}
	go (func() {
		LB.CheckServers()
	})()

}
