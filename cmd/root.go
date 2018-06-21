package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/kelseyhightower/envconfig"

	"github.com/spf13/cobra"
)

var cfgFile string

var logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)

var conf struct {
	SrcPath string
	DstPath string
	Secret  string
	Port    int    `default:"3000"`
	Handler string `default:"/_update"`
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "hugo-update",
	Short: "Update server for hugo webistes",
	Long:  `A webserver that keeps your hugo updated via web hooks`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads the configuration.
func initConfig() {
	if err := envconfig.Process("", &conf); err != nil {
		logger.Fatalln("Config:", err)
	}
	if conf.SrcPath == "" || conf.DstPath == "" || conf.Secret == "" {
		logger.Fatalln("Invalid configuration: check $SRCPATH, $DSTPATH, $SECRET")
	}
	conf.SrcPath, _ = filepath.Abs(conf.SrcPath)
	conf.DstPath, _ = filepath.Abs(conf.DstPath)
}
