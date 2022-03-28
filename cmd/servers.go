package cmd

import (
	"context"
	"fmt"
	"log"

	api "github.com/danielfsousa/flusso/api/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var serversCmd = &cobra.Command{
	Use:   "servers",
	Short: "List all servers in a flusso cluster",
	RunE:  serversRun,
}

func init() {
	rootCmd.AddCommand(serversCmd)
	serversCmd.Flags().StringP("address", "a", ":8400", "Flusso server address")
	if err := viper.BindPFlags(serversCmd.Flags()); err != nil {
		log.Fatal(err)
	}
}

func serversRun(cmd *cobra.Command, args []string) error {
	addr := viper.GetString("address")

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	client := api.NewLogClient(conn)
	ctx := context.Background()

	res, err := client.GetServers(ctx, &api.GetServersRequest{})
	if err != nil {
		return err
	}

	fmt.Println("servers:")
	for _, server := range res.Servers {
		fmt.Printf("\t- %v\n", server)
	}

	return nil
}
