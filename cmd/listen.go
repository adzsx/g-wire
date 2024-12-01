package cmd

import (
	"github.com/adzsx/gwire/internal/host"
	"github.com/spf13/cobra"
)

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "listen on a specific port.",
	Long:  "Listen on a specific port. Optionally setting up gwire chatting with the handshake option \nFor example: \ngwire listen -p 1337",

	Run: listen,
}

func init() {
	rootCmd.AddCommand(listenCmd)

	listenCmd.Flags().StringSliceP("port", "p", []string{}, "Port of host to connect to")

	listenCmd.Flags().BoolP("time", "t", true, "Display time for each message")
	listenCmd.Flags().StringP("encrypt", "e", "auto", "Encryption to set up with handshake")
	listenCmd.Flags().StringP("username", "u", "anonymous", "Username to perform handshake with")

}

func listen(cmd *cobra.Command, args []string) {
	ports, _ := cmd.Flags().GetStringSlice("port")

	exchangeInfo, _ := rootCmd.Flags().GetBool("setup")

	username, _ := cmd.Flags().GetString("username")
	enc, _ := cmd.Flags().GetString("encrypt")
	time, _ := cmd.Flags().GetBool("time")
	//timeout, _ := cmd.Flags().GetString("timeout")
	host.HostSetup(ports, exchangeInfo, enc, username, time)
}
