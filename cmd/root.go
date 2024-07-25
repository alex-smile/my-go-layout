package cmd

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"mygo/template/pkg/server"
)

var cfgFile string

func init() {
	rootCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "config file (default is config.yaml;required)")
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	rootCmd.Flags().String("addr", "", "like 0.0.0.0:8080")

	rootCmd.MarkFlagRequired("config")

	viper.BindPFlag("server.addr", rootCmd.Flags().Lookup("addr"))
}

var rootCmd = &cobra.Command{
	Use:   "mygo",
	Short: "This is mygo",
	Run: func(cmd *cobra.Command, args []string) {
		Start()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Start() {
	initConfig()

	initDatabase()
	// initRedis()

	initSentry()
	initMetrics()

	server.Run(globalConfig)
}
