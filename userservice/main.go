package main

import (
	"fmt"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"userservice/cmd/transport"
	"userservice/endpoint"
	"userservice/util"
)

func main() {

	// 初始化路由
	r := mux.NewRouter()

	userService := endpoint.UserServiceImpl{}
	// 生成登录endpoint
	loginEndpoint := endpoint.MakeUserLoginEndpoint(userService)
	loginHandler := httptransport.NewServer(loginEndpoint, transport.DecodeUserRequest, transport.EncodeUserResponse)
	r.Methods("Post").Path("/user/login").Handler(loginHandler)

	// 生成注册endpoint
	signupEndpoint := endpoint.MakeUserSignupEndpoint(userService)
	signupHandler := httptransport.NewServer(signupEndpoint, transport.DecodeUserRequest, transport.EncodeUserResponse)
	r.Methods("Post").Path("/user/signup").Handler(signupHandler)

	// 生成consul的健康检查接口
	r.Methods("GET").Path("/health").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	errChan := make(chan error)

	// 开启web服务
	// todo: 端口改为从配置中获取
	go func() {
		util.RegService()
		// todo: 随机获取一个可用的端口
		err := http.ListenAndServe(":8081", r)
		if err != nil {
			errChan <- err
		}
	}()

	// 捕获ctrl-c和kill -9信号
	go func() {
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-sigChan)
	}()

	// 程序退出前反注册service
	err := <-errChan
	util.DeRegister()
	log.Println(err)
}
