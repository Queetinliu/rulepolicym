package cmd

import (
	//"bytes"
	"fmt"
)

func ListRuleGroup(accesstoken,authtoken string) (RuleGroupResp, error) {
    //jsonBody := []byte(`{"page": 1, "size": 50,"ns": "SYSTEM"}`)
    //bodyReader := bytes.NewReader(jsonBody)


	listrulegroup_req := Request{
		Method: "GET",
		Url: Api_Monitor_Url+"rule-groups",
		Headers: map[string]string{
			"Authorization": authtoken,
			"Cookie": fmt.Sprintf("X-DCE-Access-Token=%s", accesstoken),
		},
		QueryParams: map[string]string{
			"page": "1",
			"size": "50",
			"ns": "SYSTEM",
		},
	}
	respbody, err := NewRequest[RuleGroupResp](listrulegroup_req)
	if err != nil {
		return RuleGroupResp{}, err
	}
	return respbody, nil
}

type RuleGroupResp struct {
     Data    []RuleGroupData `json:"data"`
	 Keyword string `json:"keyword"`
	 Page    int `json:"page"`
	 Size    int `json:"size"`
	 Total   int `json:"total"`
}

type RuleGroupData struct {
	Active_At int `json:"active_at"`
	Annotation map[string]string `json:"annotation"`
	Create_At int `json:"create_at"`
	Create_By string `json:"create_by"`
	Id string `json:"id"`
	Interval string `json:"interval"`
	Keywords string `json:"keywords"`
	Labels map[string]string `json:"labels"`
	Name string `json:"name"`
	Notifiers []string `json:"notifiers"`
	Notify_Type string `json:"notify_type"`
	Ns string `json:"ns"`
	Operated_By string `json:"operated_by"`
	Paused bool `json:"paused"`
	Rules []Rule `json:"rules"`
	Target string `json:"target"`
	Type string `json:"type"`
	Update_At int `json:"update_at"`
}

type Rule struct {
	Active_At int `json:"active_at"`
	Annotation map[string]string `json:"annotation"`
	Code string `json:"code"`
	Create_At int `json:"create_at"`
	Create_By string `json:"create_by"`
	Expr string `json:"expr"`
	For string `json:"for"`
	Group_Id string `json:"group_id"`
	Group_Name string `json:"group_name"`
	Id string `json:"id"`
	Keywords string `json:"keywords"`
	Labels map[string]string `json:"labels"`
	Name string `json:"name"`
	Ns string `json:"ns"`
	Operated_By string `json:"operated_by"`
	Severity string `json:"severity"`
	State string `json:"state"`
	Target string `json:"target"`
	Type string `json:"type"`
	Updated_At int `json:"updated_at"`
}