package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"minibox.ai/pkg/apiserver"
)

var serveCmd = &cobra.Command{
	Use:   "minibox",
	Short: "minibox.ai api Server",
	Long: `Minibox.AI master controller Server
                Complete documentation is available at http://docs.minibox.ai`,
	Run: func(cmd *cobra.Command, args []string) {
		apisvr := apiserver.NewApiServer()
		go func() {
			grpc := apiserver.NewGRPCServer(apisvr)
			grpc.Listen(":8080")
		}()
		log.Fatalln(apiserver.Listen(apisvr))
	},
}

func NewApiServer() *cobra.Command {
	return serveCmd
}

func init() {
	serveCmd.AddCommand(genkeyCmd)
}
