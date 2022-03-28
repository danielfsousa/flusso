package cmd

import (
	"context"
	"fmt"

	api "github.com/danielfsousa/flusso/api/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var produceCmd = &cobra.Command{
	Use:   "produce [message]",
	Short: "Appends records to a flusso server",
	Args:  cobra.MinimumNArgs(1),
	RunE:  produceRun,
}

func init() {
	rootCmd.AddCommand(produceCmd)
	produceCmd.Flags().StringP("address", "a", ":8400", "Flusso server address")
	viper.BindPFlags(produceCmd.Flags())
}

func produceRun(cmd *cobra.Command, args []string) error {
	message := args[0]
	addr := viper.GetString("address")

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	client := api.NewLogClient(conn)
	ctx := context.Background()

	res, err := client.Produce(ctx, &api.ProduceRequest{
		Record: &api.Record{
			Value: []byte(message),
		},
	})
	if err != nil {
		return err
	}

	fmt.Println("Produced 1 message")
	fmt.Printf("Offset: %d\n", res.Offset)
	return nil
}
