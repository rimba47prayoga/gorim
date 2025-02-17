/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"gorim.org/gorim/conf"
	"gorm.io/gorm/logger"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Creates and Updates database schema.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		migrationLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // Use default log package to print to console
			logger.Config{
				SlowThreshold:             time.Second,   // Slow SQL query threshold
				LogLevel:                  logger.Silent, // Log level (e.g., Silent, Error, Warn, Info)
				IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
				Colorful:                  true,          // Enable color printing
			},
		)
		conf.DB.Logger = migrationLogger
		conf.MigrationInstance.Run()
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// migrateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
