/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/labstack/gommon/color"
	"github.com/rimba47prayoga/gorim.git/conf"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gorim.git",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Check if arguments are passed and are not "migrate"
	shouldCheckMigration := false
	if len(os.Args) > 1 {
		firstArgs := os.Args[1]
		if firstArgs != "migrate" {
			if firstArgs == "runserver" {
				if os.Args[len(os.Args) - 1] == "--nomigrationcheck" {
					shouldCheckMigration = false
				} else {
					shouldCheckMigration = true
				}
			} else {
				shouldCheckMigration = true
			}
		}
	}
	if shouldCheckMigration {
		isChanged, _ := conf.MigrationInstance.HasChanges() 
		if isChanged {
			msg := strings.TrimSpace(`
				You have unapplied migration(s). Your project may not work properly until you apply the migrations.
				Run 'go run main.go migrate' to apply them.
			`)
			msg = strings.ReplaceAll(msg, "\t\t", "")
			msg = strings.TrimSpace(msg) // Clean up leading/trailing spaces
			fmt.Println(color.Red(msg))
		}
	}
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gorim.git.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


