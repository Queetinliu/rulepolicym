/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"rulepolicym/pkg"
)

// rmCmd represents the rm command
var stopruleCmd = &cobra.Command{
	Use:   "stoprule",
	Args:  cobra.ExactArgs(1),
	Short: "stop all rules in a rulegroup",
	Long:  `stop all rules in specific rulegroup`,
	RunE: func(cmd *cobra.Command, args []string) error {
		acctoken, authtoken, err := pkg.ReadToken(client)
		if err != nil {
			return err
		}
		err = pkg.ControlRuleGroup(client, acctoken, authtoken, args[0], true)
		if err != nil {
			return err
		}
		return nil

	},
}

func init() {
	rootCmd.AddCommand(stopruleCmd)
}
