package cmd

import (
	"log"

	cobra "github.com/spf13/cobra"
)

var scaffoldCmd = &cobra.Command{
	Use:   "scaffold",
	Short: "minibox.ai scaffold system",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("scaffold starting.")
		// if err := CliAuth(func(token string) error {

		// 	if prjs, err := v1.ListProjects(); err != nil {
		// 		fmt.Printf("List Project failed: %s\n", err)
		// 		return err
		// 	} else {
		// 		fmt.Printf("List Projects - \n")
		// 		for _, prj := range prjs {
		// 			fmt.Println(prj)
		// 		}
		// 	}

		// 	return nil
		// }); err != nil {
		// 	log.Printf("auth error you must do `mini login` in terminal")
		// }
	},
}

func NewScaffold() *cobra.Command {
	return scaffoldCmd
}
