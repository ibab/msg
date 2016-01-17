package cmd

import (
	"os"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func Edit(identifier int) {

}

var editCmd = &cobra.Command{
	Use:   "edit <identifier>",
	Short: "Edit a draft",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Error: Edit requires a draft identifier.")
			os.Exit(-1)
		}
		identifier, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("Error: Argument '%s' is not a valid draft identifier.\n", args[0])
			os.Exit(-1)
		}
		Edit(identifier)
	},
}

func init() {
	RootCmd.AddCommand(editCmd)
	editCmd.SetArgs([]string{"identifier"})
}
