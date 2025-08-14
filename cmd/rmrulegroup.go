/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"rulepolicym/pkg"

	"github.com/spf13/cobra"
)

// rulegroupCmd represents the rulegroup command
var rulegroupCmd = &cobra.Command{
	Use:   "rulegroup",
	Short: "delete a rulegroup",
	Long:  `delete the specified rulegroup`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := pkg.HttpClient()
		accesstoken, authtoken, err := pkg.ReadToken(client)
		if err != nil {
			return err
		}
		err = pkg.DeleteRuleGroup(client, accesstoken, authtoken, args[0])
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rmCmd.AddCommand(rulegroupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rulegroupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rulegroupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
