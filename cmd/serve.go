package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/danielfsousa/flusso/internal/agent"
	"github.com/danielfsousa/flusso/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cli = serveCli{}

var serveCmd = &cobra.Command{
	Use:     "serve",
	Short:   "Starts the flusso server",
	PreRunE: cli.setupConfig,
	RunE:    cli.run,
}

func init() {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dataDir := path.Join(homedir, ".flusso", "data")

	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().String("data-dir", dataDir, "Directory to store log and Raft data.")
	serveCmd.Flags().String("node-name", hostname, "Unique server ID.")
	serveCmd.Flags().String("bind-addr", "127.0.0.1:8401", "Address to bind Serf on.")
	serveCmd.Flags().Int("rpc-port", 8400, "Port for RPC clients (and Raft) connections.")
	serveCmd.Flags().StringSlice("start-join-addrs", nil, "Serf addresses to join.")
	serveCmd.Flags().Bool("bootstrap", false, "Bootstrap the cluster.")

	serveCmd.Flags().String("acl-model-file", "", "Path to ACL model.")
	serveCmd.Flags().String("acl-policy-file", "", "Path to ACL policy.")

	serveCmd.Flags().String("server-tls-cert-file", "", "Path to server tls cert.")
	serveCmd.Flags().String("server-tls-key-file", "", "Path to server tls key.")
	serveCmd.Flags().String("server-tls-ca-file", "", "Path to server certificate authority.")

	serveCmd.Flags().String("peer-tls-cert-file", "", "Path to peer tls cert.")
	serveCmd.Flags().String("peer-tls-key-file", "", "Path to peer tls key.")
	serveCmd.Flags().String("peer-tls-ca-file", "", "Path to peer certificate authority.")

	viper.BindPFlags(serveCmd.Flags())
}

type serveCli struct {
	cfg serveCfg
}

type serveCfg struct {
	agent.Config
	ServerTLSConfig config.TLSConfig
	PeerTLSConfig   config.TLSConfig
}

func (s *serveCli) setupConfig(cmd *cobra.Command, args []string) error {
	var err error

	s.cfg.DataDir = viper.GetString("data-dir")
	s.cfg.NodeName = viper.GetString("node-name")
	s.cfg.BindAddr = viper.GetString("bind-addr")
	s.cfg.RPCPort = viper.GetInt("rpc-port")
	s.cfg.StartJoinAddrs = viper.GetStringSlice("start-join-addrs")
	s.cfg.Bootstrap = viper.GetBool("bootstrap")
	s.cfg.ACLModelFile = viper.GetString("acl-model-file")
	s.cfg.ACLPolicyFile = viper.GetString("acl-policy-file")
	s.cfg.ServerTLSConfig.CertFile = viper.GetString("server-tls-cert-file")
	s.cfg.ServerTLSConfig.KeyFile = viper.GetString("server-tls-key-file")
	s.cfg.ServerTLSConfig.CAFile = viper.GetString("server-tls-ca-file")
	s.cfg.PeerTLSConfig.CertFile = viper.GetString("peer-tls-cert-file")
	s.cfg.PeerTLSConfig.KeyFile = viper.GetString("peer-tls-key-file")
	s.cfg.PeerTLSConfig.CAFile = viper.GetString("peer-tls-ca-file")

	if s.cfg.ServerTLSConfig.CertFile != "" && s.cfg.ServerTLSConfig.KeyFile != "" {
		s.cfg.ServerTLSConfig.Server = true
		s.cfg.Config.ServerTLSConfig, err = config.SetupTLSConfig(s.cfg.ServerTLSConfig)
		if err != nil {
			return err
		}
	}
	if s.cfg.PeerTLSConfig.CertFile != "" && s.cfg.PeerTLSConfig.KeyFile != "" {
		s.cfg.Config.PeerTLSConfig, err = config.SetupTLSConfig(s.cfg.PeerTLSConfig)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *serveCli) run(cmd *cobra.Command, args []string) error {
	agent, err := agent.New(s.cfg.Config)
	if err != nil {
		return err
	}
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	<-sigc
	return agent.Shutdown()
}
