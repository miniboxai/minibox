package cmd

import (
	"context"
	"fmt"

	cobra "github.com/spf13/cobra"
	v1 "minibox.ai/minibox/pkg/api/v1"
	types "minibox.ai/minibox/pkg/api/v1/types"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "new minibox resources, like project, kernel.",
}

var newProjectCmd = &cobra.Command{
	Use:   "project [name of project]",
	Short: "Create new project at minibox.ai",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := AuthClient(clientOpt, func(client *v1.Clients) error {
			fmt.Printf("Creating project %s...\n", args[0])
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if prj, err := client.CreateProject(ctx, &types.CreateProjectRequest{Name: args[0]}); err != nil {
				fmt.Printf("created project failed: %s\n", err)
			} else {
				fmt.Printf("Project ID: %s\n", prj.ID)
			}
			return nil
		}); err != nil {
			fmt.Printf("auth error you must do `mini login` in terminal")
		}
	},
}

func init() {
	newCmd.AddCommand(newProjectCmd)
}
