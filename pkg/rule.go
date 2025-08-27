package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"text/tabwriter"
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

func ListRules(client *http.Client, accesstoken, authtoken, rulegroupid string) (RuleGroup, error) {

	urlpath, err := GetApiMonitorUrl("rule-groups", rulegroupid)

	//urlpath, err := url.JoinPath(os.Getenv(Dce_Url_Env_Name), Api_Monitor_Path_Prefix, "rule-groups", rulegroupid)
	if err != nil {
		return RuleGroup{}, err
	}
	listrulesreq := Request{
		Client:  client,
		Method:  "GET",
		Url:     urlpath,
		Headers: apimonitorheader(accesstoken, authtoken),
	}
	listrulesresp, err := NewRequest[RuleGroup](listrulesreq)
	if err != nil {
		return RuleGroup{}, err
	}
	return listrulesresp, nil
}

func (rg RuleGroup) Print() {
	w := tabwriter.NewWriter(os.Stdout, 16, 8, 0, '\t', 0)
	// Write some data to the Writer.
	fmt.Fprintf(w, "\n %s\t%s\t%s\t%s\t", "规则", "状态", "告警周期", "告警级别")
	for _, rule := range rg.Rules {
		fmt.Fprintf(w, "\n %s\t%s\t%s\t%s\t", rule.Name, rule.State, rule.For, rule.Severity)
	}
	fmt.Fprintf(w, "\n")
	// Flush the Writer to ensure all data is written to the output.
	w.Flush()
}

func DeleteRule(client *http.Client, accesstoken, authtoken, rulename, rulegroupname string) error {
	rulegroupid, err := GetRuleGroupId(client, accesstoken, authtoken, rulegroupname)
	if err != nil {
		return err
	}
	rulegroup, err := ListRules(client, accesstoken, authtoken, rulegroupid)
	if err != nil {
		return err
	}

	var ruleid string
	for _, rule := range rulegroup.Rules {
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
		Client:  client,
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

type PausePayload struct {
	Paused bool `json:"paused"`
}

func ControlRule(client *http.Client, accesstoken, authtoken, ruleid, rulegroupid string, paused bool) error {
	urlpath, err := GetApiMonitorUrl("rule-groups", rulegroupid, "rules", ruleid, "pause")
	if err != nil {
		return err
	}
	requestbody, err := json.Marshal(PausePayload{Paused: paused})
	if err != nil {
		return err
	}
	pauserulereq := Request{
		Client:  client,
		Method:  "PUT",
		Url:     urlpath,
		Headers: apimonitorheader(accesstoken, authtoken),
		Body:    bytes.NewBuffer(requestbody),
	}
	_, err = NewRequest[struct{}](pauserulereq)
	if err != nil {
		return err
	}
	return nil
}
