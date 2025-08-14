package pkg

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	// Sso_Url = "http://30.1.64.241/dce/sso/login"
	Sso_Path         = "/dce/sso/login"
	Dce_Url_Env_Name = "DCE_URL"
)

type Printer interface {
	print()
}

type SsoResp struct {
	Access_Token string `json:"access_token"`
}


func HttpClient() *http.Client {
    client := &http.Client{Timeout: 10 * time.Second}
    return client
}


func apimonitorheader(accesstoken, authtoken string) map[string]string {
	return map[string]string{
		"Authorization": authtoken,
		"Cookie":        fmt.Sprintf("X-DCE-Access-Token=%s", accesstoken),
	}
}

func GetApiMonitorUrl(pathelements ...string) (string, error) {
	newpathelements := append([]string{os.Getenv("ApiMonitor_Path_Prefix")}, pathelements...)
	return url.JoinPath(os.Getenv(Dce_Url_Env_Name), newpathelements...)
}

func GetAccessToken(client *http.Client) (SsoResp, error) {
	// from the OS environment to get username and password
	dceurl := os.Getenv(Dce_Url_Env_Name)
	Sso_Url, err := url.JoinPath(dceurl, Sso_Path)
	if err != nil {
		return SsoResp{}, fmt.Errorf("url join path %s %s with err:%w", dceurl, Sso_Path, err)
	}
	username := os.Getenv("DCE_USERNAME")
	password := os.Getenv("DCE_PASSWORD")
	userpass := username + ":" + password
	userpassbase64 := base64.StdEncoding.EncodeToString([]byte(userpass))

	ssoreq := Request{
		Client: client,
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

func GetAuthToken(client *http.Client, access_token string) (AuthResp, error) {
	urlpath, err := GetApiMonitorUrl("auth")
	if err != nil {
		return AuthResp{}, err
	}
	getauth_req := Request{
		Client: client,
		Method: "POST",
		Url:    urlpath,
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
	Client      *http.Client
	Method      string
	Url         string
	Headers     map[string]string
	QueryParams map[string]string
	Body        io.Reader
}

type Response interface {
	SsoResp | AuthResp | RuleGroups | RuleGroup | struct{}
}

func NewRequest[R Response](request Request) (R, error) {
	// do a request
	var respbody R
	// client := &http.Client{}
	req, err := http.NewRequest(request.Method, request.Url, request.Body)
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
	client := request.Client
	if client == nil {
		client = http.DefaultClient
	}
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
	// if
	if _, ok := any(respbody).(struct{}); ok {
		return respbody, nil
	}

	err = json.Unmarshal(bodyText, &respbody)
	if err != nil {
		return respbody, err
	}
	return respbody, nil
}
