package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/spf13/cobra"
)

func openEditor(file string) {
	cmd := exec.Command("vim", file)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func Edit(identifier int) {
	content := []byte("Your mail here")
	tmpfile, err := ioutil.TempFile("", "msg-")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}

	openEditor(tmpfile.Name())

	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
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
