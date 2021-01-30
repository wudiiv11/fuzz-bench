package endpoint

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"time"
	"userservice/db"
	"userservice/util"
)

type UserService interface {
	Signup(username string, passwd string) bool
	Login(username string, passwd string) (string, error)
}

type UserServiceImpl struct{}

func (acceptor UserServiceImpl) Signup(username string, passwd string) bool {
	encPwd := util.SHA1([]byte(passwd))
	return db.Signup(username, encPwd)
}

func (acceptor UserServiceImpl) Login(username string, passwd string) (string, error) {
	encPwd := util.SHA1([]byte(passwd))
	if !db.LogIn(username, encPwd) {
		return "", errors.New("登录失败")
	}
	token := genToken(username)
	db.UpdateToken(username, token)
	return token, nil
}

func genToken(username string) string {
	// todo: 换成JWT库生成的token
	// 40位字符, md5(username + timestamp + token_salt) + timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + "hash_salt"))
	return tokenPrefix + ts[:8]
}

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Err   string `json:"err,omitempty"`
}

type SignupResponse struct {
	Suc bool `json:"suc"`
}

func MakeUserLoginEndpoint(service UserService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UserRequest)
		token, err := service.Login(req.Username, req.Password)
		if err != nil {
			return nil, err
		}
		return util.NewRespMsg(0, "登录成功", token), nil
	}
}

func MakeUserSignupEndpoint(service UserService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UserRequest)
		suc := service.Signup(req.Username, req.Password)
		if suc {
			return util.NewRespMsg(0, "注册成功", ""), nil
		}
		return util.NewRespMsg(-1, "注册失败", ""), nil
	}
}
