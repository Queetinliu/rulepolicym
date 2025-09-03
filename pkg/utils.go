package pkg

import (
	//"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	Sso_Path         = "/dce/sso/login"
	Dce_Url_Env_Name = "DCE_URL"
)

type Printer interface {
	Print()
}

func HttpClient() *http.Client {
	client := &http.Client{Timeout: 30 * time.Second}
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

type Request struct {
	Client      *http.Client
	Method      string
	Url         string
	Headers     map[string]string
	QueryParams map[string]interface{}
	Body        io.Reader
}

type Response interface {
	SsoResp | AuthResp | RuleGroups | RuleGroup | struct{}
}

func NewRequest[R Response](request Request) (R, error) {
	// do a request
	var respbody R

	// client := &http.Client{}
	// print the request.Body for debug.
	// if request.Body != nil {
	// 	bodyBytes, err := io.ReadAll(request.Body)
	// 	if err != nil {
	// 		return respbody, fmt.Errorf("read request body with err:%w", err)
	// 	}
	// 	fmt.Println(string(bodyBytes))
	// 	// After reading, we must replace the now-empty body with a new reader.
	// 	request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	// }

	req, err := http.NewRequest(request.Method, request.Url, request.Body)
	if err != nil {
		return respbody, fmt.Errorf("create new request to %s with err:%w", request.Url, err)
	}
	for k, v := range request.Headers {
		req.Header.Add(k, v)
	}
	q := req.URL.Query()
	for k, v := range request.QueryParams {
		switch v := v.(type) {
		case int:
			q.Add(k, fmt.Sprintf("%d", v))
		case int64:
			q.Add(k, fmt.Sprintf("%d", v))
		case bool:
			q.Add(k, fmt.Sprintf("%v", v))
		case string:
			q.Add(k, v)
		default:
			q.Add(k, fmt.Sprintf("%v", v))
		}
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
	if _, ok := any(respbody).(struct{}); ok {
		return respbody, nil
	}

	err = json.Unmarshal(bodyText, &respbody)
	if err != nil {
		return respbody, err
	}
	return respbody, nil
}
