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
	Run: func(cmd *cobra.Command, args []string) {
		accesstoken, authtoken, err := pkg.ReadToken()
		if err != nil {
			fmt.Println(err)
			return
		}
		rulegroupresp, err := pkg.ListRuleGroups(accesstoken, authtoken)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, rulegroup := range rulegroupresp.Data {
			fmt.Println(rulegroup.Name)
		}
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
