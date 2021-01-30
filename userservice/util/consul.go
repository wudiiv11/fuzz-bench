package util

import (
	"github.com/hashicorp/consul/api"
	"log"
)

var (
	consulClient *api.Client
)

func init() {
	config := api.DefaultConfig()
	// todo: 如何从网络中抓取注册中心的地址
	config.Address = "127.0.0.1:8500"
	client, err := api.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}
	consulClient = client
}

func RegService() {
	reg := api.AgentServiceRegistration{}
	reg.Address = "192.168.0.104"
	reg.ID = "userservice1"
	reg.Name = "userservice"
	reg.Port = 8081
	reg.Tags = []string{"primary"}

	check := api.AgentServiceCheck{}
	check.Interval = "5s"
	check.HTTP = "http://192.168.0.104:8081/health"

	reg.Check = &check

	err := consulClient.Agent().ServiceRegister(&reg)
	if err != nil {
		log.Fatal(err)
	}
}

func DeRegister() {
	err := consulClient.Agent().ServiceDeregister("userservice1")
	if err != nil {
		log.Fatal(err)
	}
}
