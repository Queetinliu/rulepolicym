package pkg

import (
	"net/http"
)

type Rule struct {
	Active_At   int               `json:"active_at"`
	Annotations map[string]string `json:"annotations"`
	Code        string            `json:"code"`
	Create_At   int               `json:"create_at"`
	Create_By   string            `json:"create_by"`
	Expr        string            `json:"expr"`
	For         string            `json:"for"`
	Group_Id    string            `json:"group_id"`
	Group_Name  string            `json:"group_name"`
	Id          string            `json:"id"`
	Keywords    string            `json:"keywords"`
	Labels      map[string]string `json:"labels"`
	Name        string            `json:"name"`
	Ns          string            `json:"ns"`
	Operated_By string            `json:"operated_by"`
	Severity    string            `json:"severity"`
	State       string            `json:"state"`
	Target      string            `json:"target"`
	Type        string            `json:"type"`
	Updated_At  int               `json:"updated_at"`
}

type RuleGroupResp struct {
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

func ListRules(client *http.Client, accesstoken, authtoken, rulegroupid string) ([]Rule, error) {

	urlpath, err := GetApiMonitorUrl("rule-groups", rulegroupid)

	//urlpath, err := url.JoinPath(os.Getenv(Dce_Url_Env_Name), Api_Monitor_Path_Prefix, "rule-groups", rulegroupid)
	if err != nil {
		return nil, err
	}
	listrulesreq := Request{
		Client:  client,
		Method:  "GET",
		Url:     urlpath,
		Headers: apimonitorheader(accesstoken, authtoken),
	}
	listrulesresp, err := NewRequest[RuleGroupResp](listrulesreq)
	if err != nil {
		return nil, err
	}
	return listrulesresp.Rules, nil
}

func DeleteRule(client *http.Client, accesstoken, authtoken, rulename, rulegroupname string) error {
	rulegroupid, err := GetRuleGroupId(client, accesstoken, authtoken, rulegroupname)
	if err != nil {
		return err
	}
	listrules, err := ListRules(client, accesstoken, authtoken, rulegroupid)
	if err != nil {
		return err
	}

	var ruleid string
	for _, rule := range listrules {
		if rule.Name == rulename {
			ruleid = rule.Id
		}
	}
	urlpath, err := GetApiMonitorUrl("rule-groups", rulegroupid, "rules", ruleid)
	if err != nil {
		return err
	}
	//urlpath, err := url.JoinPath(os.Getenv(Dce_Url_Env_Name), Api_Monitor_Path_Prefix, "rule-groups", rulegroupid, "rules", ruleid)

	deleterulereq := Request{
		Client: client,
		Method:  "DELETE",
		Url:     urlpath,
		Headers: apimonitorheader(accesstoken, authtoken),
	}
	_, err = NewRequest[struct{}](deleterulereq)
	if err != nil {
		return err
	}
	return nil

}
