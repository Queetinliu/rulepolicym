/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"rulepolicym/pkg"

	"github.com/spf13/cobra"
)

var rmrulegroup *string

// ruleCmd represents the rule command
var ruleCmd = &cobra.Command{
	Use:   "rule",
	Short: "delete a rule from a rulegroup",
	Long:  `delete a rule from the specified rulegroup`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := pkg.HttpClient()
		accesstoken, authtoken, err := pkg.ReadToken(client)
		if err != nil {
			return err
		}
		err = pkg.DeleteRule(client, accesstoken, authtoken, args[0], *rmrulegroup)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rmCmd.AddCommand(ruleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ruleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ruleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rmrulegroup = ruleCmd.Flags().StringP("rulegroup", "g", "kubernetes-alert", "specific the rulegroup name(required)")
}
