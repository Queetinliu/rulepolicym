package pkg

import (
	"fmt"
)

const (
	Api_Monitor_Url = "http://30.1.64.241/dce/proxy/clusters/17263D42-22D8-45B1-45BB-D9451377FB56/plugin/29002/api/monitor/"
)

func ListRuleGroup(accesstoken, authtoken string) (RuleGroupResp, error) {
	listrulegroup_req := Request{
		Method: "GET",
		Url:    Api_Monitor_Url + "rule-groups",
		Headers: apimonitorheader(accesstoken, authtoken),
		QueryParams: map[string]string{
			"page": "1",
			"size": "50",
			"ns":   "SYSTEM",
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
	Keyword string          `json:"keyword"`
	Page    int             `json:"page"`
	Size    int             `json:"size"`
	Total   int             `json:"total"`
}

type RuleGroupData struct {
	Active_At   int               `json:"active_at"`
	Annotation  map[string]string `json:"annotation"`
	Create_At   int               `json:"create_at"`
	Create_By   string            `json:"create_by"`
	Id          string            `json:"id"`
	Interval    string            `json:"interval"`
	Keywords    string            `json:"keywords"`
	Labels      map[string]string `json:"labels"`
	Name        string            `json:"name"`
	Notifiers   []string          `json:"notifiers"`
	Notify_Type string            `json:"notify_type"`
	Ns          string            `json:"ns"`
	Operated_By string            `json:"operated_by"`
	Paused      bool              `json:"paused"`
	Rules       []Rule            `json:"rules"`
	Target      string            `json:"target"`
	Type        string            `json:"type"`
	Update_At   int               `json:"update_at"`
}

func GetRuleGroupId(accesstoken, authtoken, name string) (string, error) {
	rulegroupresp, err := ListRuleGroup(accesstoken, authtoken)
	if err != nil {
		return "", err
	}
	for _, rulegroup := range rulegroupresp.Data {
		if rulegroup.Name == name {
			return rulegroup.Id, nil
		}
	}
	return "", fmt.Errorf("not found rulegroup name:%s", name)
}
