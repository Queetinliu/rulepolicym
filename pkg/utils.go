package pkg

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	Sso_Url = "http://30.1.64.241/dce/sso/login"
)

type SsoResp struct {
	Access_Token string `json:"access_token"`
}

func apimonitorheader(accesstoken, authtoken string) map[string]string {
	return map[string]string{
		"Authorization": authtoken,
		"Cookie":        fmt.Sprintf("X-DCE-Access-Token=%s", accesstoken),
	}
}

func GetAccessToken() (SsoResp, error) {
	// from the OS environment to get username and password
	username := os.Getenv("DCE_USERNAME")
	password := os.Getenv("DCE_PASSWORD")
	userpass := username + ":" + password
	userpassbase64 := base64.StdEncoding.EncodeToString([]byte(userpass))

	ssoreq := Request{
		Method: "POST",
		Url:    Sso_Url,
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Basic %s", userpassbase64),
		},
	}
	respbody, err := NewRequest[SsoResp](ssoreq)
	if err != nil {
		return SsoResp{}, err
	}
	return respbody, nil
}

type AuthResp struct {
	Token   string  `json:"token"`
	Account Account `json:"account"`
}

type Account struct {
	Emails       []string `json:"emails"`
	ID           string   `json:"id"`
	Is_Admin     bool     `json:"is_admin"`
	Language     string   `json:"language"`
	Name         string   `json:"name"`
	Namespaces   []string `json:"namespaces"`
	NickName     string   `json:"nick_name"`
	Phone        string   `json:"phone"`
	Primary_Mail string   `json:"primary_mail"`
	Status       string   `json:"status"`
}

func GetAuthToken(access_token string) (AuthResp, error) {

	getauth_req := Request{
		Method: "POST",
		Url:    Api_Monitor_Url + "auth",
		Headers: map[string]string{
			"X-DCE-ACCESS-TOKEN": access_token,
		},
	}
	respbody, err := NewRequest[AuthResp](getauth_req)
	if err != nil {
		return AuthResp{}, err
	}
	return respbody, nil
}

type Request struct {
	Method      string
	Url         string
	Headers     map[string]string
	QueryParams map[string]string
}

type Response interface {
	SsoResp | AuthResp | RuleGroupResp
}

func NewRequest[R Response](request Request) (R, error) {
	// do a request
	var respbody R
	client := &http.Client{}
	req, err := http.NewRequest(request.Method, request.Url, nil)
	if err != nil {
		return respbody, fmt.Errorf("create new request to %s with err:%w", request.Url, err)

	}
	for k, v := range request.Headers {
		req.Header.Add(k, v)
	}
	q := req.URL.Query()
	for k, v := range request.QueryParams {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	//fmt.Printf("%#v",req.Header)
	resp, err := client.Do(req)
	if err != nil {
		return respbody, fmt.Errorf("do request to %s ,body: %s with err:%w", request.Url, req.Body, err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return respbody, err
	}
	//fmt.Println(string(bodyText))
	err = json.Unmarshal(bodyText, &respbody)
	if err != nil {
		return respbody, err
	}
	return respbody, nil
}
