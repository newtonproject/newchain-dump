package cli

import (
	"github.com/spf13/cobra"
)

func (cli *CLI) buildRootCmd() {

	if cli.rootCmd != nil {
		cli.rootCmd.ResetFlags()
		cli.rootCmd.ResetCommands()
	}

	rootCmd := &cobra.Command{
		Use:              cli.Name,
		Short:            cli.Name + " is commandline client for extract NewChain blocks and transactions.",
		Run:              cli.help,
		PersistentPreRun: cli.setup,
	}
	cli.rootCmd = rootCmd

	// Global flags
	rootCmd.PersistentFlags().StringVarP(&cli.config, "config", "c", defaultConfigFile, "The `path` to config file")
	rootCmd.PersistentFlags().StringP("rpcURL", "i", defaultRPCURL, "Geth json rpc or ipc `url`")

	rootCmd.PersistentFlags().String("host", defaultHost, "The host for database")
	rootCmd.PersistentFlags().String("user", "", "The user for database")
	rootCmd.PersistentFlags().String("database", "", "The name of database")
	rootCmd.PersistentFlags().String("password", "", "The password for database")

	rootCmd.PersistentFlags().StringP("log", "l", defaultLogFile, "The path of log file")

	// Basic commands
	rootCmd.AddCommand(cli.buildInitCmd())    // init
	rootCmd.AddCommand(cli.buildVersionCmd()) // version

	// run
	rootCmd.AddCommand(cli.buildRunCmd())

}
