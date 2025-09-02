package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"text/tabwriter"
)

func ListRuleGroups(client *http.Client, accesstoken, authtoken string) (RuleGroups, error) {
	urlpath, err := GetApiMonitorUrl("rule-groups")
	if err != nil {
		return RuleGroups{}, err
	}

	listrulegroup_req := Request{
		Client:  client,
		Method:  "GET",
		Url:     urlpath,
		Headers: apimonitorheader(accesstoken, authtoken),
		QueryParams: map[string]interface{}{
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
	Keyword string      `json:"keyword"`
	Page    int         `json:"page"`
	Size    int         `json:"size"`
	Total   int         `json:"total"`
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

func (rgs RuleGroups) Print() {
	w := tabwriter.NewWriter(os.Stdout, 8, 8, 0, '\t', 0)
	// Write some data to the Writer.
	fmt.Fprintf(w, "\n %s\t%s\t", "策略", "告警规则数")
	for _, rulegroup := range rgs.Data {
		fmt.Fprintf(w, "\n %s\t%s\t", rulegroup.Name, strconv.Itoa(len(rulegroup.Rules)))
	}
	fmt.Fprintf(w, "\n")
	// Flush the Writer to ensure all data is written to the output.
	w.Flush()
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
	if destrulegroup.Target == "cluster"{
		destrulegroup.Target = "{\"cluster\":\"cluster\"}"
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
		Client:  client,
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
		Client:  client,
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

func RemoveDuplicateRule(client *http.Client, accesstoken, authtoken, sourcerulegroupname, baserulegroupname string) error {
	sourcerulegroupid, err := GetRuleGroupId(client, accesstoken, authtoken, sourcerulegroupname)
	if err != nil {
		return err
	}
	baserulegroupid, err := GetRuleGroupId(client, accesstoken, authtoken, baserulegroupname)
	if err != nil {
		return err
	}
	baserulemap := make(map[string]string)
	baserulegroup, err := ListRules(client, accesstoken, authtoken, baserulegroupid)
	if err != nil {
		return err
	}
	for _, rule := range baserulegroup.Rules {
		baserulemap[rule.Name] = rule.Expr
	}
	sourcerulegroup, err := ListRules(client, accesstoken, authtoken, sourcerulegroupid)
	if err != nil {
		return err
	}
	for _, rule := range sourcerulegroup.Rules {
		_, ok := baserulemap[rule.Name]
		if ok {
			if baserulemap[rule.Name] != rule.Expr {
				fmt.Printf("rule %s expr is not same in %s and %s,please delete it manually\n", rule.Name, baserulegroupname, sourcerulegroupname)
			} else {
				err = DeleteRule(client, accesstoken, authtoken, rule.Name, sourcerulegroupname)
				if err != nil {
					return err
				}

			}
		} else {
			fmt.Printf("rule %s is not in %s\n", rule.Name, baserulegroupname)
		}
	}
	return nil

}

var pausemap = map[bool]string{
	true:  "disabled",
	false: "enabled",
}

func ControlRuleGroup(client *http.Client, accesstoken, authtoken, rulegroupname string, paused bool) error {
	rulegroupid, err := GetRuleGroupId(client, accesstoken, authtoken, rulegroupname)
	if err != nil {
		return err
	}
	rulegroup, err := ListRules(client, accesstoken, authtoken, rulegroupid)
	if err != nil {
		return err
	}
	// stop,paused true,the state should be disabled
	// start,pause false,the state shoule be enabled
	for _, rule := range rulegroup.Rules {
		if rule.State != pausemap[paused] {
			err = ControlRule(client, accesstoken, authtoken, rule.Id, rulegroupid, paused)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
