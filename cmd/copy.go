/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"rulepolicym/pkg"

	"github.com/spf13/cobra"
)

// copyCmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "copy a rulegroup to a new rulegroup",
	Long: `copy a rulegroup to a new rulegroup,include all rules under the rulegroup,you must
	specify the source rulegroup and the destination rulegroup`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := pkg.HttpClient()
		accesstoken, authtoken, err := pkg.ReadToken(client)
		if err != nil {
			return err
		}
		err = pkg.CopyRuleGroup(client, accesstoken, authtoken, args[0], args[1])
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// copyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// copyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
