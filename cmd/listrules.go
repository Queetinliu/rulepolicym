/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"rulepolicym/pkg"
)

// rulesCmd represents the rules command
var rulesCmd = &cobra.Command{
	Use:   "rules",
	Short: "list rules in a rulegroup",
	Long:  `list rules in a specific rulegroup name, default is kubernetes-alert`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := pkg.HttpClient()
		accesstoken, authtoken, err := pkg.ReadToken(client)
		if err != nil {
			return err
		}
		rulegroupid, err := pkg.GetRuleGroupId(client, accesstoken, authtoken, *listrulegroup)
		if err != nil {
			return err
		}
		rulegroup, err := pkg.ListRules(client, accesstoken, authtoken, rulegroupid)
		if err != nil {
			return err
		}
		var p pkg.Printer = rulegroup
		p.Print()
		return nil
	},
}

var listrulegroup *string

func init() {
	listCmd.AddCommand(rulesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rulesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rulesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	listrulegroup = rulesCmd.Flags().StringP("rulegroup", "g", "kubernetes-alert", "specific the rulegroup name(required)")
}
