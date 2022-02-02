package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
)

type HttpMethod string

const (
	GET    HttpMethod = "GET"
	POST   HttpMethod = "POST"
	PATCH  HttpMethod = "PATCH"
	PUT    HttpMethod = "PUT"
	DELETE HttpMethod = "DELETE"
)

func SendFormRequest(uri string, method HttpMethod, param map[string]string, respP interface{}) ([]byte, error) {
	req, err := FormRequest(uri, method, param)
	if err != nil {
		return nil, err
	}
	return SendRequest(req, respP)
}

func SendJsonRequest(uri string, method HttpMethod, param, respP interface{}) ([]byte, error) {
	fmt.Printf("url: %s\nmethod: %s\njson: %v\n", uri, string(method), param)
	req, err := JsonRequest(uri, method, param)
	if err != nil {
		return nil, err
	}
	return SendRequest(req, respP)
}

func SendGetRequest(uri string, param, respP interface{}) ([]byte, error) {
	req, err := GetRequest(uri, param)
	if err != nil {
		return nil, err
	}
	return SendRequest(req, respP)
}

func GetRequest(uri string, param interface{}) (*http.Request, error) {
	uri = QueryString(uri, param)
	return http.NewRequest(string(GET), uri, nil)
}

func FormRequest(uri string, method HttpMethod, param interface{}) (*http.Request, error) {
	params, ok := param.(map[string]string)
	if !ok {
		params = make(map[string]string)
		jsonData, err := json.Marshal(param)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(jsonData, &params)
		if err != nil {
			return nil, err
		}
	}

	form := url.Values{}
	for key, value := range params {
		if value == "" {
			continue
		}
		form.Add(key, value)
	}

	// 構造post請求
	req, err := http.NewRequest(string(method), uri, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}

func JsonRequest(uri string, method HttpMethod, param interface{}) (*http.Request, error) {
	body, _ := json.Marshal(param)
	return RawJsonRequest(uri, method, body)
}

func RawJsonRequest(uri string, method HttpMethod, body []byte) (*http.Request, error) {
	// 構造post請求
	req, err := http.NewRequest(string(method), uri, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func QueryString(uri string, param interface{}) string {
	if param == nil {
		return uri
	}

	q, err := query.Values(param)
	if err == nil {
		qs := q.Encode()
		uri = fmt.Sprintf("%s?%s", uri, qs)
	}

	return uri
}

func ResponseParser(resp *http.Response, respP interface{}) ([]byte, error) {
	body, err := ReadBody(resp, respP)
	if err != nil {
		return nil, err
	}

	httpCode := resp.StatusCode
	if httpCode != http.StatusOK {
		return nil, fmt.Errorf("httpcode:%d not 200, %s", httpCode, string(body))
	}

	return body, nil
}

func ReadBody(resp *http.Response, respP interface{}) ([]byte, error) {
	// 讀取響應body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if respP != nil && !IsNilInterfaceObject(respP) {
		return body, json.Unmarshal(body, respP)
	}

	return body, nil
}

func SendRequest(req *http.Request, respP interface{}) ([]byte, error) {
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return ResponseParser(res, respP)
}
