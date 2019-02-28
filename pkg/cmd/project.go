package cmd

import (
	"context"
	"fmt"
	"log"

	cobra "github.com/spf13/cobra"
	v1 "minibox.ai/pkg/api/v1"
	"minibox.ai/pkg/api/v1/types"
)

var projectCmd = &cobra.Command{
	Use:   "projects",
	Short: "Sign in minibox.ai website, use apis",
	Long: `Used minibox.ai dev account to sign in, can use apis to managament project, download/upload data, 
		running a training jobs..
        Complete documentation is available at http://docs.minibox.ai`,
	Aliases: []string{"project"},
	Run: func(cmd *cobra.Command, args []string) {
		clientOpt = LoadClientConfig()

		err := AuthClient(clientOpt, func(client *v1.Clients) error {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if prjs, err := client.ListProjects(ctx, &types.ListProjectsRequest{}); err != nil {
				fmt.Printf("List Project failed: %s\n", err)
				return err
			} else {
				fmt.Printf("List Projects - \n")
				for _, prj := range prjs.Projects {
					fmt.Println(prj)
				}
			}

			return nil
		})

		if err != nil {
			log.Printf("auth error you must do `mini login` in terminal")
		}
	},
}
