package main

import (
	"fmt"
	"rulepolicym/cmd"
)

func main() {
	ssoresp, err := cmd.GetToken()
	if err != nil {
		panic(err)
	}
	fmt.Println(ssoresp.Access_Token)
	authresp, err := cmd.GetAuth(ssoresp.Access_Token)
	if err != nil {
		panic(err)
	}
	fmt.Println(authresp.Token)
	rulegroupresp, err := cmd.ListRuleGroup(ssoresp.Access_Token,authresp.Token)
	if err != nil {
		panic(err)
	}
	for _,rulegroup := range rulegroupresp.Data {
		fmt.Println(rulegroup.Name)
	}
}
