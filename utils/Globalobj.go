package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"zinx-demo/ziface"

	"gopkg.in/yaml.v2"
)

type GlobalConfig struct {
	TcpServer ziface.IServer //当前zinx全局的server对象
	Host      string         `yaml:"Host"`    //当前服务器主机监听IP
	TcpPort   int            `yaml:"TcpPort"` //当前服务器主机监听的端口号
	Name      string         `yaml:"Name"`    //当前服务器的名称

	Version          string //当前zinx版本号
	MaxConn          int    `yaml:"MaxConn"` //服务器最大连接数
	MaxPackageSize   uint32 //当前框架数据包最大值
	WorkerPoolSize   uint32 //当前业务工作Worker池的goroutine数量
	MaxWorkerTaskLen uint32 //该框架允许的最大worker池数量
}

// 从cfg/zinx.json加载配置参数
func (g *GlobalConfig) Reload() {
	ReadJson()
	//ReadYaml()
}

func ReadYaml() {
	fmt.Println("Read yaml...")
	conf := &GlobalConfig{}
	if file, err := os.Open("./cfg/zinx.yaml"); err != nil {
		fmt.Println("os open error:", err)
		return
	} else {
		err := yaml.NewDecoder(file).Decode(conf)
		if err != nil {
			fmt.Println("yaml decoder error:", err)
			return
		}
	}

	fmt.Println("conf:", conf)

}
func ReadJson() {
	//读取文件
	data, err := os.ReadFile("./cfg/zinx.json")
	if err != nil {
		fmt.Println("ReadFile error:", err)
		return
	}
	//fmt.Println("data:", string(data))
	//解析json文件
	error := json.Unmarshal(data, &GlobalConf)
	if error != nil {
		fmt.Println("json unmarshal error:", error)
		return
	}
}

var GlobalConf *GlobalConfig

func init() {
	// 如果配置文件未加载，默认值
	GlobalConf = &GlobalConfig{
		Name:             "ZinxServerApp",
		Version:          "V1.0",
		TcpPort:          8999,
		Host:             "0.0.0.0",
		MaxConn:          1000,
		MaxPackageSize:   2048,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 10,
	}

	//加载配置文件
	GlobalConf.Reload()

}
