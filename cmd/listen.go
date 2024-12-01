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

	listenCmd.Flags().StringP("port", "p", "", "Port of host to connect to")

	listenCmd.Flags().BoolP("handshake", "s", true, "Perform a handshake with another gwire client")
	listenCmd.Flags().BoolP("time", "t", true, "Display time for each message")
	listenCmd.Flags().StringP("encrypt", "e", "auto", "Encryption to set up with handshake")
	listenCmd.Flags().StringP("username", "u", "anonymous", "Username to perform handshake with")

}

func listen(cmd *cobra.Command, args []string) {
	port, _ := cmd.Flags().GetString("port")

	exchangeInfo, _ := cmd.Flags().GetBool("handshake")

	username, _ := cmd.Flags().GetString("username")
	enc, _ := cmd.Flags().GetString("encrypt")
	time, _ := cmd.Flags().GetBool("time")
	//timeout, _ := cmd.Flags().GetString("timeout")
	host.HostSetup(port, exchangeInfo, enc, username, time)
}
