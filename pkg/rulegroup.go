package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"text/tabwriter"

	"github.com/goccy/go-yaml"
	//"gopkg.in/yaml.v3"
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
	if destrulegroup.Target == "cluster" {
		destrulegroup.Target = "{\"cluster\":\"cluster\"}"
	}

	for i := range destrulegroup.Rules {
		// if rules.Severity != "critical" {
		// 	destrulegroup.Rules[i].Severity = "critical"
		// 	_, ok := destrulegroup.Rules[i].Labels["severity"]
		// 	if ok {
		// 		destrulegroup.Rules[i].Labels["severity"] = "critical"
		// 	} else {
		// 		destrulegroup.Rules[i].Labels["severity"] = "critical"
		// 	}
		// }
			_, ok := destrulegroup.Rules[i].Annotations["ruleGroupName"]
			if ok {
				destrulegroup.Rules[i].Annotations["ruleGroupName"] = dest
			} else {
				destrulegroup.Rules[i].Annotations["ruleGroupName"] = dest
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

type Annotations struct {
	Summary     string `yaml:"summary"`
	Description string `yaml:"description"`
}
type RuleLabels struct {
	Severity string `yaml:"severity"`
}
type AlertRule struct {
	Alert       string `yaml:"alert"`
	Expr        string `yaml:"expr"`
	For         string `yaml:"for"`
	RuleLabels  `yaml:"labels"`
	Annotations `yaml:"annotations"`
}

type AlertRuleGroup struct {
	Name  string      `yaml:"name"`
	Rules []AlertRule `yaml:"rules"`
}

type RuleFile struct {
	Groups []AlertRuleGroup `yaml:"groups"`
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
	Annotations  map[string]string `json:"annotations"`
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



type CreateRulePayload struct {
	Name        string            `json:"name"`
	Expr        string            `json:"expr"`
	FormExpr    string            `json:"formExpr"`
	For         string            `json:"for"`
	ForValue    string            `json:"forValue"`
	ForUnit     string            `json:"forUnit"`
	Severity    string            `json:"severity"`
	Annotations CreateRuleAnnotations `json:"annotations"`
	Index       int               `json:"__index"`
}

type CreateRuleAnnotations struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`
}

type CreateRuleGroupPayload struct {
	Name        string              `json:"name"`
	Interval    string              `json:"interval"`
	Keywords    string              `json:"keywords"`
	Labels      map[string]string   `json:"labels"`
	Notifiers   []string            `json:"notifiers"`
	Rules       []CreateRulePayload `json:"rules"`
	Target     string   `json:"target"`
	Type        string              `json:"type"`
	Notify_Type string              `json:"notify_type"`
	Ns          string              `json:"ns"`
}

func CreateRuleGroup(client *http.Client, accesstoken, authtoken, rulefilename string) error {
	rulecontent, err := os.ReadFile(rulefilename)
	if err != nil {
		return err
	}
	// unmarshal the rulefile
	var rulefile RuleFile
	err = yaml.Unmarshal(rulecontent, &rulefile)
	if err != nil {
		return err
	}

	for _, group := range rulefile.Groups {
		rulespayload := make([]CreateRulePayload, 0)
		for i, rule := range group.Rules {
			// use regex from the rule.For to extract the value and unit.
			// the value is the number part, the unit is the string part.
			regex := regexp.MustCompile(`(\d+)([a-zA-Z]+)`)
			matches := regex.FindStringSubmatch(rule.For)
			if len(matches) != 3 {
				return fmt.Errorf("invalid for format: %s", rule.For)
			}
			if rule.RuleLabels.Severity == "warning" {
				rule.RuleLabels.Severity = "warn"
			}
			createrulepayload := CreateRulePayload{
				Name:     rule.Alert,
				Expr:     rule.Expr,
				FormExpr: "",
				For:      rule.For,
				ForValue: matches[1],
				ForUnit:  matches[2],
				Severity: rule.RuleLabels.Severity,
				Annotations: CreateRuleAnnotations(rule.Annotations),
				Index: i,
			}
			rulespayload = append(rulespayload, createrulepayload)

		}

		createrulegrouppayload := CreateRuleGroupPayload{
			Name:      group.Name,
			Interval:  "10m",
			Keywords:  "",
			Labels:    map[string]string{},
			Notifiers: []string{},
			Rules:     rulespayload,
			Target: "{\"cluster\":\"cluster\"}",
			Type:        "cluster",
			Notify_Type: "",
			Ns:          "SYSTEM",
		}
		urlpath, err := GetApiMonitorUrl("rule-groups")
		if err != nil {
			return err
		}
		requestbody, err := json.Marshal(createrulegrouppayload)
		if err != nil {
			return err
		}
		//fmt.Println(string(requestbody))
		createrulegroupreq := Request{
			Client:  client,
			Method:  "POST",
			Url:     urlpath,
			Headers: apimonitorheader(accesstoken, authtoken),
			Body:    bytes.NewReader(requestbody),
		}
		createrulegroupreq.Headers["Content-Type"]="application/json"
		newrulegroup, err := NewRequest[RuleGroup](createrulegroupreq)
		if err != nil {
			return err
		}
		fmt.Printf("create rulegroup %s success,id is %s\n", newrulegroup.Name, newrulegroup.Id)
	}
	return nil
}
