package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
// Api_Monitor_Url = "http://30.1.64.241/dce/proxy/clusters/17263D42-22D8-45B1-45BB-D9451377FB56/plugin/29002/api/monitor/"
// ApiMonitor_Path_Prefix = "dce/proxy/clusters/17263D42-22D8-45B1-45BB-D9451377FB56/plugin/29002/api/monitor/"
)

func ListRuleGroups(client *http.Client, accesstoken, authtoken string) (RuleGroups, error) {
	urlpath, err := GetApiMonitorUrl("rule-groups")
	if err != nil {
		return RuleGroups{}, err
	}

	listrulegroup_req := Request{
		Client: client,
		Method:  "GET",
		Url:     urlpath,
		Headers: apimonitorheader(accesstoken, authtoken),
		QueryParams: map[string]string{
			"page": "1",
			"size": "50",
			"ns":   "SYSTEM",
		},
	}
	respbody, err := NewRequest[RuleGroups](listrulegroup_req)
	if err != nil {
		return RuleGroups{}, err
	}
	return respbody, nil
}

type RuleGroups struct {
	Data    []RuleGroup `json:"data"`
	Keyword string           `json:"keyword"`
	Page    int              `json:"page"`
	Size    int              `json:"size"`
	Total   int              `json:"total"`
}

type RuleGroup struct {
	Active_At   int               `json:"active_at"`
	Annotation  map[string]string `json:"annotation"`
	Built_in    bool              `json:"built_in"`
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


func GetRuleGroupId(client *http.Client, accesstoken, authtoken, name string) (string, error) {
	rulegroupresp, err := ListRuleGroups(client, accesstoken, authtoken)
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

func CopyRuleGroup(client *http.Client, accesstoken, authtoken, source, dest string) error {
	listrulegroups, err := ListRuleGroups(client, accesstoken, authtoken)
	if err != nil {
		return err
	}
	var sourcerulegroup, destrulegroup RuleGroup
	for _, rulegroup := range listrulegroups.Data {
		if rulegroup.Name == source {
			sourcerulegroup = rulegroup
		}
	}
	destrulegroup = sourcerulegroup
	destrulegroup.Name = dest
	destrulegroup.Built_in = false
	if destrulegroup.Notify_Type == "" {
		destrulegroup.Notify_Type = "script"
	}
	for i, rules := range destrulegroup.Rules {
		if rules.Severity != "critical" {
			destrulegroup.Rules[i].Severity = "critical"
			_, ok := destrulegroup.Rules[i].Labels["severity"]
			if ok {
				destrulegroup.Rules[i].Labels["severity"] = "critical"
			} else {
				destrulegroup.Rules[i].Labels["severity"] = "critical"
			}
			_, ok = destrulegroup.Rules[i].Annotations["ruleGroupName"]
			if ok {
				destrulegroup.Rules[i].Annotations["ruleGroupName"] = dest
			} else {
				destrulegroup.Rules[i].Annotations["ruleGroupName"] = dest
			}
		}

	}
	destgroupdata, err := json.Marshal(destrulegroup)
	if err != nil {
		return fmt.Errorf("marshal the destgroup to json with err:%w", err)
	}
	urlpath, err := GetApiMonitorUrl("rule-groups")
	if err != nil {
		return err
	}

	createrulegroupreq := Request{
		Client: client,
		Method:  "POST",
		Url:     urlpath,
		Headers: apimonitorheader(accesstoken, authtoken),
		Body:    bytes.NewBuffer(destgroupdata),
	}
	_, err = NewRequest[RuleGroup](createrulegroupreq)
	if err != nil {
		return fmt.Errorf("create new rulegroup with err:%w", err)
	}
	return nil

}

// type DeleteRuleGroupResp struct {
// }

func DeleteRuleGroup(client *http.Client, accesstoken, authtoken, rulegroupname string) error {
	rulegroupid, err := GetRuleGroupId(client, accesstoken, authtoken, rulegroupname)
	if err != nil {
		return err
	}

	urlpath, err := GetApiMonitorUrl("rule-groups", rulegroupid)
	if err != nil {
		return err
	}

	deleterulegroupreq := Request{
		Client: client,
		Method:  "DELETE",
		Url:     urlpath,
		Headers: apimonitorheader(accesstoken, authtoken),
	}
	_, err = NewRequest[struct{}](deleterulegroupreq)
	if err != nil {
		return err
	}
	return nil

}
