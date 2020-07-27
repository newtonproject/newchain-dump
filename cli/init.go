package cli

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/console"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func (cli *CLI) buildInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "init",
		Short:                 "Initialize config file",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {

			fmt.Println("Initialize config file")

			prompt := fmt.Sprintf("Enter file in which to save (%s): ", defaultConfigFile)
			configPath, err := console.Stdin.PromptInput(prompt)
			if err != nil {
				fmt.Println("PromptInput err:", err)
			}
			if configPath == "" {
				configPath = defaultConfigFile
			}
			cli.config = configPath

			prompt = fmt.Sprintf("Enter path of log file (%s): ", defaultLogFile)
			logfile, err := console.Stdin.PromptInput(prompt)
			if err != nil {
				fmt.Println("PromptInput err:", err)
			}
			if logfile == "" {
				logfile = defaultLogFile
			}
			cli.logfile = logfile

			rpcURLV := viper.GetString("rpcURL")
			prompt = fmt.Sprintf("Enter geth json rpc or ipc url (%s): ", rpcURLV)
			cli.rpcURL, err = console.Stdin.PromptInput(prompt)
			if err != nil {
				fmt.Println("PromptInput err:", err)
			}
			if cli.rpcURL == "" {
				cli.rpcURL = rpcURLV
			}
			viper.Set("rpcURL", cli.rpcURL)

			prompt = fmt.Sprintf("Configure MySQL database or not: [Y/n] ")
			configDB, err := console.Stdin.PromptInput(prompt)
			if err != nil {
				fmt.Println("PromptInput err:", err)
			}
			if len(configDB) <= 0 {
				configDB = "Y"
			}
			if strings.ToUpper(configDB[:1]) == "Y" {

				dbhost := defaultHost
				prompt = fmt.Sprintf("Enter database host(%s): ", dbhost)
				cli.host, err = console.Stdin.PromptInput(prompt)
				if err != nil {
					fmt.Println("PromptInput err:", err)
				}
				if cli.host == "" {
					cli.host = dbhost
				}
				viper.Set("mysql.Host", cli.host)

				prompt = fmt.Sprintf("Enter database name: ")
				cli.database, err = console.Stdin.PromptInput(prompt)
				if err != nil {
					fmt.Println("PromptInput err:", err)
				}
				viper.Set("mysql.Database", cli.database)

				prompt = fmt.Sprintf("Enter the username to connect to the database: ")
				cli.user, err = console.Stdin.PromptInput(prompt)
				if err != nil {
					fmt.Println("PromptInput err:", err)
				}
				viper.Set("mysql.User", cli.user)

				prompt = fmt.Sprintf("Enter the password for user: ")
				cli.password, err = console.Stdin.PromptPassword(prompt) // console.Stdin.PromptInput(prompt)
				if err != nil {
					fmt.Println("PromptInput err:", err)
				}
				viper.Set("mysql.Password", cli.password)

			}

			err = viper.WriteConfigAs(configPath)
			if err != nil {
				fmt.Println("WriteConfig:", err)
			} else {
				fmt.Println("Your configuration has been saved in ", configPath)
			}
		},
	}

	return cmd
}
