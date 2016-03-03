package cmd

import (
	"os"
	"os/exec"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func Edit(identifier int) {
	cmd := exec.Command("vim", "tmp")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error received from vim\n")
		os.Exit(-1)
	}
}

var editCmd = &cobra.Command{
	Use:   "edit <identifier>",
	Short: "Edit a draft",
	Long:  ``,
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
