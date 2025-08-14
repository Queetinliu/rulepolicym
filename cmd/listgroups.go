/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"rulepolicym/pkg"
)

// listgroupsCmd represents the listgroups command
var listgroupsCmd = &cobra.Command{
	Use:   "rulegroups",
	Short: "list all rulegroups",
	Long:  `list all groups in the system`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := pkg.HttpClient()
		accesstoken, authtoken, err := pkg.ReadToken(client)
		if err != nil {
			return err
		}
		rulegroupresp, err := pkg.ListRuleGroups(client, accesstoken, authtoken)
		if err != nil {
			return err
		}
		for _, rulegroup := range rulegroupresp.Data {
			fmt.Println(rulegroup.Name)
		}
		return nil
	},
}

func init() {
	listCmd.AddCommand(listgroupsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listgroupsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listgroupsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
