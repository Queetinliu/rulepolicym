package pkg

import "net/url"

type Rule struct {
	Active_At   int               `json:"active_at"`
	Annotation  map[string]string `json:"annotation"`
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

func ListRules(accesstoken, authtoken, rulegroupid string) ([]Rule, error) {
	urlpath, err := url.JoinPath(Api_Monitor_Url, "rule-groups", rulegroupid)
	if err != nil {
		return nil, err
	}
	listrulesreq := Request{
		Method: "GET",
		Url:    urlpath,
		Headers: apimonitorheader(accesstoken, authtoken),
	}
	listrulesresp, err := NewRequest[Rule](listrulesreq)
	if err != nil {
		return nil, err
	}
}
