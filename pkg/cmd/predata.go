package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	cobra "github.com/spf13/cobra"
	"github.com/spf13/viper"
	v1 "minibox.ai/pkg/api/v1"
	"minibox.ai/pkg/api/v1/types"
	"minibox.ai/pkg/predata"
)

var (
	force bool
)

var predataCmd = &cobra.Command{
	Use:   "predata",
	Short: "Prepare data for minibox Project",
	// Run: func(cmd *cobra.Command, args []string) {
	// 	if len(args) < 2 {

	// 	}

	// 	fmt.Printf("Predata src: %s to %s", args[0], args[1])
	// 	if err := CliAuth(prepareDataExec); err != nil {
	// 		fmt.Printf("auth error you must do `mini login` in terminal")
	// 	}
	// },
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync Dataset",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var importCmd = &cobra.Command{
	Use:   "import <src> <target>",
	Short: "Import Dataset to Project",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			src    = args[0]
			target = args[1]
		)

		site := viper.GetString("apiserver")
		predata.SetHubSite(site)

		source, u, err := predata.ParseURL(src)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		switch source {
		case predata.HubData, predata.UserData:
			err = predata.ImportFromHub(*u, target)
		default:
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Data from: %s url: %s\n", source, u)
	},
}

var lsCmd = &cobra.Command{
	Use:   "ls [bucket]",
	Short: "List your buckets",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var query string

		if len(args) == 1 {
			query = args[0]
		}
		clientOpt = LoadClientConfig()
		log.Printf("query: %s", query)

		if err := AuthClient(clientOpt, func(client *v1.Clients) error {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			reply, err := client.ListDatasets(ctx, &types.ListDatasetsRequest{Namespace: query})

			if err != nil {
				fmt.Printf("List Dataset failed %s\n", err)
				return err
			}

			output := os.Stdout

			w := tabwriter.NewWriter(output, 0, 0, 4, ' ', 0)

			fmt.Fprintf(output, "Datasets List ---\n\n")
			for _, ds := range reply.Datasets {
				fmt.Fprintf(w, "\t%s\n", ds.Name)
			}
			w.Flush()

			// pretty.Println(reply)
			return nil
		}); err != nil {
			fmt.Println("auth error you must do `mini login` in terminal")
		}
	},
}

var mdCmd = &cobra.Command{
	Use:   "md [dataset]",
	Short: "Make your owner dataset",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			dataset = args[0]
		)
		clientOpt = LoadClientConfig()

		site := viper.GetString("apiserver")
		predata.SetHubSite(site)
		if err := AuthClient(clientOpt, func(client *v1.Clients) error {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			reply, err := client.CreateDataset(ctx, &types.CreateDatasetRequest{Name: dataset})
			if err != nil {
				fmt.Printf("Create Dataset failed %s", err)
				os.Exit(1)
			}

			fmt.Printf("%s\n", reply)
			return nil
		}); err != nil {
			fmt.Printf("auth error you must do `mini login` in terminal")
		}
	},
}

func prepareDataExec(token string) error {
	return nil
}

func PredataCommand() *cobra.Command {
	return predataCmd
}

func init() {
	predataCmd.AddCommand(syncCmd)
	predataCmd.AddCommand(importCmd)
	predataCmd.AddCommand(lsCmd)
	predataCmd.AddCommand(mdCmd)

	flag := predataCmd.PersistentFlags()
	flag.BoolVarP(&force, "force", "f", false, "Force downloading Dataset")
}
