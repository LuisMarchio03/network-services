package cmd

import (
	"FTP_Server/pkg/FTP"
	"fmt"
	"github.com/spf13/cobra"
)

var ftpCmd = &cobra.Command{
	Use:   "ftp-server",
	Short: "Inicia o servidor FTP",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Iniciando o servidor FTP...")
		ftp.StartFTPServer()
	},
}

func init() {
	rootCmd.AddCommand(ftpCmd)
}
