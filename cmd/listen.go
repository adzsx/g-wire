package cmd

import (
	"log"
	"os"

	"github.com/adzsx/gwire/internal/host"
	"github.com/adzsx/gwire/internal/utils"
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

	listenCmd.Flags().IntP("port", "p", 0, "Port of host to connect to")
	listenCmd.Flags().StringP("source", "s", "", "Filter to single IP source IP")

	listenCmd.Flags().BoolP("time", "t", true, "Display time for each message")
	listenCmd.Flags().BoolP("encrypt", "e", false, "Encryption to set up with handshake")
}

func listen(cmd *cobra.Command, args []string) {
	port, _ := cmd.Flags().GetInt("port")

	if port == 0 {
		log.Fatalln("Port not specified")
		os.Exit(1)
	}

	src, _ := cmd.Flags().GetString("source")
	enc, _ := cmd.Flags().GetBool("encrypt")
	time, _ := cmd.Flags().GetBool("time")
	utils.Verbose, _ = rootCmd.Flags().GetInt("verbose")

	//timeout, _ := cmd.Flags().GetString("timeout")
	host.HostSetup(port, src, enc, time)
}
