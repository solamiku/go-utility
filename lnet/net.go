package lent

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/solamiku/go-utility/compress"
)

type Cookies map[string]string
type Param map[string]string

/*
	basic http authentication
	https://en.wikipedia.org/wiki/Basic_access_authentication
*/

type BasicAuth struct {
	User string
	Pass string
}

// return basic authentication string
// base64encode(username+passwrod)
func (ba *BasicAuth) String() string {
	sKey := ba.User + ":" + ba.Pass
	r := string(compress.Base64encode([]byte(sKey)))
	return r
}

func HttpGet(sUrl, body string, cookies Cookies, basicAuth ...BasicAuth) (string, int, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", sUrl, strings.NewReader(body))
	if err != nil {
		return "", 0, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//set cookie
	sCookie := ""
	for key, val := range cookies {
		if len(sCookie) > 0 {
			sCookie += "&"
		}
		sCookie += (key + "=" + val)
	}
	req.Header.Set("Cookie", sCookie)
	if len(basicAuth) > 0 {
		req.Header.Set("Authorization", "Basic "+basicAuth[0].String())
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	rbody, err := ioutil.ReadAll(resp.Body)
	return string(rbody), resp.StatusCode, nil
}

func HttpPost(sUrl string, params Param) (body []byte, statusCode int, err error) {
	hc := http.Client{}

	form := url.Values{}
	for k, v := range params {
		form.Set(k, v)
	}
	req, err := http.NewRequest("POST", sUrl, strings.NewReader(form.Encode()))
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := hc.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	statusCode = resp.StatusCode
	return
}

func HttpBasicPost(sUrl string, cookies Cookies, params Param, basicAuth BasicAuth) (string, int, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	form := url.Values{}
	for k, v := range params {
		form.Set(k, v)
	}
	req, err := http.NewRequest("POST", sUrl, strings.NewReader(form.Encode()))
	if err != nil {
		return "", 0, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	//set cookie
	sCookie := ""
	for key, val := range cookies {
		if len(sCookie) > 0 {
			sCookie += "&"
		}
		sCookie += (key + "=" + val)
	}
	req.Header.Set("Cookie", sCookie)

	req.Header.Set("Authorization", "Basic "+basicAuth.String())

	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	rbody, err := ioutil.ReadAll(resp.Body)
	return string(rbody), resp.StatusCode, nil
}
