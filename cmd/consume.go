package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	api "github.com/danielfsousa/flusso/api/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var consumeCmd = &cobra.Command{
	Use:   "consume",
	Short: "Consume records from a flusso server",
	RunE:  consumeRun,
}

func init() {
	consumeCmd.Flags().StringP("address", "a", ":8400", "Flusso server address")
	rootCmd.AddCommand(consumeCmd)
	viper.BindPFlags(consumeCmd.Flags())
}

func consumeRun(cmd *cobra.Command, args []string) error {
	addr := viper.GetString("address")

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	client := api.NewLogClient(conn)
	ctx := context.Background()

	stream, err := client.ConsumeStream(ctx, &api.ConsumeRequest{
		Offset: 0,
	})
	if err != nil {
		return err
	}

	done := make(chan bool, 1)
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Starting flusso consumer. Use Ctrl-C to exit.")

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Printf("Failed to receive message: %s", err)
				continue
			}
			fmt.Printf("%s\n", res.Record.Value)
		}
		done <- true
	}()

	// waits for sigchan or done
	select {
	case <-sigchan:
	case <-done:
	}

	return nil
}
