package common

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin/binding"
)

// support ids[]=0&ids[]=1
type ArrayQueryBinding struct{}

func NewArrayQueryBinding() ArrayQueryBinding {
	return ArrayQueryBinding{}
}

func (ArrayQueryBinding) Name() string {
	return "query"
}

func (ArrayQueryBinding) Bind(req *http.Request, obj interface{}) error {
	urlValues := req.URL.Query()
	for key, vs := range urlValues {
		values := make([]string, 0)
		for _, v := range vs {
			if v == "" {
				continue
			}

			splitDatas := strings.Split(v, ",")
			values = append(values, splitDatas...)
		}

		urlValues.Del(key)

		key = strings.Replace(key, "[]", "", -1)
		for _, v := range values {
			urlValues.Add(key, v)
		}
	}

	req.URL.RawQuery = urlValues.Encode()

	return binding.Query.Bind(req, obj)
}
