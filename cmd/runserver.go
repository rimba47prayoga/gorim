/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/rimba47prayoga/gorim.git"
	"github.com/rimba47prayoga/gorim.git/conf"
	"github.com/spf13/cobra"
)

// runserverCmd represents the runserver command
var runserverCmd = &cobra.Command{
	Use:   "runserver",
	Short: "Start a lightweight Web server for development.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		server := conf.GorimServer.(*gorim.Server)
		address := fmt.Sprintf("%s:%d", conf.HOST, conf.PORT)
		versionNumber := "v1.1.0"
		printBanner(versionNumber, conf.ENV_PATH, conf.HOST)
		err := server.Start(address)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runserverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runserverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runserverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func printBanner(version string, env string, host string) {
    // Current timestamp
    now := time.Now().Format("January 02, 2006 - 15:04:05")

    // Banner content
    fmt.Println("System check identified no issues (0 silenced).")
    fmt.Println(now)
	fmt.Printf("Gorim version %s using env '%s'\n", version, env)
    fmt.Printf("Starting development server at %s\n", host)
    fmt.Println("Quit the server with CONTROL-C.")
}
