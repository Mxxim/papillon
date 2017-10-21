package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
)

func init() {
	RootCmd.AddCommand(generateCmd)

}

var generateCmd = &cobra.Command{
	Use:   "gen",
	Short: "Genrate a static blog website.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2{
			fmt.Println("Errors:")
			return
		}
	},
}