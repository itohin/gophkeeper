package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "A brief description of your command",
	Long:  `A longer description.`,
	Run: func(cmd *cobra.Command, args []string) {
		login()
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func login() {
	loginPrompt := promptContent{}
	loginPrompt.label = "Please provide a login: "
	loginPrompt.errorMsg = "login error"
	login := promptGetInput(loginPrompt)
	passwordPrompt := promptContent{
		"password error",
		"Please provide a password: ",
		42,
	}
	password := promptGetInput(passwordPrompt)

	log.Println(login, password)
	//check credentials

	menuPrompt := promptContent{}
	menuPrompt.label = "Please provide action: "
	menuPrompt.errorMsg = "action error"

	action := promptGetSelect(menuPrompt, []string{"Add data", "Get data"})
	fmt.Println(action)
}
