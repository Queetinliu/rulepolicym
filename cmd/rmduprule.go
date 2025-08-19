/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"rulepolicym/pkg"

	"github.com/spf13/cobra"
)

var sourcerulegroup, baserulegroup *string

// rmduplicateruleCmd represents the rule command
var rmduplicateruleCmd = &cobra.Command{
	Use:   "rmduprule",
	Short: "remove the duplicate rule from a rulegroup",
	Long:  `remove the duplicate rule from the specified rulegroup by compare the rule expr and name in the based rulegroup`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := pkg.HttpClient()
		accesstoken, authtoken, err := pkg.ReadToken(client)
		if err != nil {
			return err
		}
		err = pkg.RemoveDuplicateRule(client, accesstoken, authtoken, *sourcerulegroup, *baserulegroup)
		if err != nil {
			return err
		}
		fmt.Printf("remove duplicate rule from %s success\n", *sourcerulegroup)
		return nil
	},
}

func init() {
	rmCmd.AddCommand(rmduplicateruleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ruleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ruleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	sourcerulegroup = rmduplicateruleCmd.Flags().StringP("sourcerulegroup", "s", "kubernetes-alert", "specific the rulegroup name(required)")
	baserulegroup = rmduplicateruleCmd.Flags().StringP("baserulegroup", "b", "kubernetes-alert", "specific the base rulegroup name(required)")
}
