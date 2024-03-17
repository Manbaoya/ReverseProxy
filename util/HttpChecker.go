package util

import (
	"fmt"
	"net/http"
	"time"
)

type HttpChecker struct {
	Servers    ServerSlice
	FailMax    int
	SuccessMin int
}

func (this *HttpChecker) Check(timeout time.Duration) {
	client := http.Client{}
	for _, v := range this.Servers {
		res, err := client.Head(v.Host)
		defer res.Body.Close()
		if res != nil {
			if res.StatusCode != 200 {
				this.Fail(v)
				continue
			} else {
				this.Success(v)
				continue
			}
		}
		if err != nil {
			this.Fail(v)
			fmt.Println(err)
			continue
		}

	}
}

func (this *HttpChecker) Fail(server *HttpServer) {
	if server.FailSum >= this.FailMax {
		server.Status = "Down"
	} else {
		server.FailSum++
	}
	server.SuccessSum = 0
	server.FailWeight += server.Weight / 5
	if server.FailWeight > server.Weight {
		server.FailWeight = server.Weight
	}
}
func (this *HttpChecker) Success(server *HttpServer) {
	if server.SuccessSum == this.SuccessMin {
		server.SuccessSum = 0
		server.FailSum = 0
		server.Status = "UP"
		server.FailWeight = 0
	} else {
		if server.FailSum > 0 {
			server.FailSum--
		}
		server.SuccessSum++

	}

}
