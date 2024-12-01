package cmd

import (
	"github.com/adzsx/gwire/internal/client"
	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "connect to a host using IP and port",
	Long:  "This command connects to a host on a specific port. \nFor example: \ngwire connect -h 192.168.0.1 -p 1337",

	Run: connect,
}

func init() {
	rootCmd.AddCommand(connectCmd)

	connectCmd.Flags().StringP("host", "d", "", "Host/destination to connect to")
	connectCmd.Flags().StringP("port", "p", "", "Port of host to connect to")

	connectCmd.Flags().BoolP("time", "t", true, "Display time for each message")
	connectCmd.Flags().StringP("encrypt", "e", "auto", "Encryption to set up with handshake")
	connectCmd.Flags().StringP("username", "u", "anonymous", "Username to perform handshake with")

}

func connect(cmd *cobra.Command, args []string) {
	host, _ := cmd.Flags().GetString("host")
	port, _ := cmd.Flags().GetString("port")

	exchangeInfo, _ := rootCmd.Flags().GetBool("setup")

	username, _ := cmd.Flags().GetString("username")
	enc, _ := cmd.Flags().GetString("encrypt")
	time, _ := cmd.Flags().GetBool("time")
	//timeout, _ := cmd.Flags().GetString("timeout")
	client.ClientSetup(host, port, exchangeInfo, enc, username, time)
}
