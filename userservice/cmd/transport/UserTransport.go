package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"userservice/endpoint"
	"userservice/util"
)

func DecodeUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	decoder := json.NewDecoder(r.Body)
	request := endpoint.UserRequest{}
	err := decoder.Decode(&request)
	if err != nil {
		return request, err
	}
	return request, nil
}

func EncodeUserResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-type", "application/json")
	resp := response.(*util.RespMsg)
	_, err := w.Write(util.JsonBytes(resp))
	return err
}
