package lent

import (
	"net/http"
	"testing"
)

func Test_net(t *testing.T) {
	_, code, err := HttpGet("http://httpbin.org/ip", "", nil)
	if err != nil || code != http.StatusOK {
		t.Error("status code do not conform to expected.")
	}
	_, code, err = HttpGet("http://10.2.48.241:9010/", "", nil, BasicAuth{
		"admin",
		"admin1",
	})
	if err != nil || code != http.StatusOK {
		t.Error("status code do not conform to expected.")
	}

	body, code, err := HttpPost("http://httpbin.org/post", Param{
		"user": "lipm",
		"from": "china",
	})
	if err != nil || code != http.StatusOK {
		t.Error("status code do not conform to expected.")
	} else {
		t.Log(string(body))
	}
}
