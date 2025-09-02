/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"rulepolicym/pkg"
)

var rulefile *string

func init() {
	rulefile = createruleCmd.Flags().StringP("rulefile", "f", "", "rule file")
}

// createruleCmd represents the createrule command
var createruleCmd = &cobra.Command{
	Use:   "createrule",
	Short: "create a rule in a rulegroup",
	Long:  `create a rule in specific rulegroup`,
	RunE: func(cmd *cobra.Command, args []string) error {
		acctoken, authtoken, err := pkg.ReadToken(client)
		if err != nil {
			return err
		}
		err = pkg.CreateRuleGroup(client, acctoken, authtoken, *rulefile)
		if err != nil {
			return err
		}
		return nil

	},
}

func init() {
	rootCmd.AddCommand(createruleCmd)
}
