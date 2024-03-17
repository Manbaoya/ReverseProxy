package util

import (
	"fmt"
	"github.com/go-ini/ini"
	"os"
)

var ProxyConfig map[string]string

type EnvConfig *os.File

func init() {
	ProxyConfig = make(map[string]string)
	EnvConfig, err := ini.Load("./env.ini")
	if err != nil {
		fmt.Println(err)
		return
	}
	proxy, err := EnvConfig.GetSection("proxy") //获取父分区
	if err != nil {
		fmt.Println(err)
		return
	}
	secs := proxy.ChildSections() //获取子分区
	for _, sec := range secs {
		path, _ := sec.GetKey("Path")
		pass, _ := sec.GetKey("Pass")
		if path != nil && pass != nil {
			ProxyConfig[path.Value()] = pass.Value()
		}
	}
}
