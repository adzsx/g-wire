/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"

	"github.com/adzsx/gwire/internal/utils"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "gwire",
	Version: "v1.4.0",
	Short:   "A tool for encrypted p2p network chatting",
	Long: `Gwire is a tool for chatting on a direct connection without a middleman
It connects using a host IP and a port and uses either AES with a preset key, or generates a key via RSA.
All of the chatting is encrypted with AES.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.test.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().IntP("verbose", "v", 0, "verbose output")
	rootCmd.PersistentFlags().BoolP("setup", "s", false, "Setup a gwire chat with another gwire client")
	rootCmd.PersistentFlags().BoolP("no-fmt", "f", false, "disable ANSI formatting")

	format, err := rootCmd.Flags().GetBool("no-fmt")
	log.Println(err)
	log.Println(format)
	utils.Format = !format
	utils.Verbose, _ = rootCmd.Flags().GetInt("verbose")
	log.Println(utils.Verbose)
}
