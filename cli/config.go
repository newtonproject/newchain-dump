package cli

import (
	"os"

	"github.com/spf13/viper"
)

const defaultConfigFile = "./config.toml"
const defaultLogFile = "./error.log"
const defaultRPCURL = "https://rpc1.newchain.newtonproject.org"
const defaultHost = "127.0.0.1:3306"

func defaultConfig(cli *CLI) {
	viper.BindPFlag("rpcURL", cli.rootCmd.PersistentFlags().Lookup("rpcURL"))
	viper.BindPFlag("log", cli.rootCmd.PersistentFlags().Lookup("log"))
	viper.SetDefault("rpcURL", defaultRPCURL)
	viper.SetDefault("log", defaultLogFile)

	viper.BindPFlag("mysql.Host", cli.rootCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("mysql.User", cli.rootCmd.PersistentFlags().Lookup("user"))
	viper.BindPFlag("mysql.Database", cli.rootCmd.PersistentFlags().Lookup("database"))
	viper.BindPFlag("mysql.Password", cli.rootCmd.PersistentFlags().Lookup("password"))
}

func setupConfig(cli *CLI) error {

	//var ret bool
	var err error

	defaultConfig(cli)

	viper.SetConfigName(defaultConfigFile)
	viper.AddConfigPath(".")
	cfgFile := cli.config
	if cfgFile != "" {
		if _, err = os.Stat(cfgFile); err == nil {
			viper.SetConfigFile(cfgFile)
			err = viper.ReadInConfig()
		} else {
			// The default configuration is enabled.
			//fmt.Println(err)
			err = nil
		}
	} else {
		// The default configuration is enabled.
		err = nil
	}

	cli.rpcURL = viper.GetString("rpcURL")
	cli.logfile = viper.GetString("log")

	return err
}
