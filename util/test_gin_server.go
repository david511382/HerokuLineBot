package util

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

type TestGinServer struct {
	router        *gin.Engine
	requestSetter func(req *http.Request) error
}

func NewTestServer(router *gin.Engine) *TestGinServer {
	return &TestGinServer{
		router: router,
	}
}

func (s TestGinServer) Get(uri string, param interface{}) *http.Response {
	req := s.getRequest(uri, param)
	if s.requestSetter != nil {
		s.requestSetter(req)
	}
	return s.sendRequest(req)
}

func (s TestGinServer) getRequest(uri string, param interface{}) *http.Request {
	uri = QueryString(uri, param)
	return httptest.NewRequest(string(GET), uri, nil)
}

func (s TestGinServer) PostJson(uri string, param interface{}) *http.Response {
	req := s.jsonRequest(uri, POST, param)
	if s.requestSetter != nil {
		s.requestSetter(req)
	}
	return s.sendRequest(req)
}

func (s TestGinServer) jsonRequest(uri string, method HttpMethod, param interface{}) *http.Request {
	body, _ := json.Marshal(param)
	// 構造post請求
	req := httptest.NewRequest(string(method), uri, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	return req
}

func (s *TestGinServer) SetRequest(requestSetter func(req *http.Request) error) {
	s.requestSetter = requestSetter
}

func (s TestGinServer) sendRequest(req *http.Request) *http.Response {
	// 初始化響應
	w := httptest.NewRecorder()

	// 呼叫相應的handler介面
	s.router.ServeHTTP(w, req)

	// 提取響應
	response := w.Result()

	return response
}
