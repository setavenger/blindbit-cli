/*Package cmd admaksdmalsdmalkd*/
package cmd

import (
	"context"
	"fmt"
	"github.com/setavenger/blindbit-cli/lib"
	"github.com/setavenger/blindbit-cli/lib/pb"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"google.golang.org/grpc"
	"log"
	"os"
)

// unlockCmd represents the unlock command
var unlockCmd = &cobra.Command{
	Use:   "unlock",
	Short: "unlocks the daemon",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		client, conn := lib.NewClient(socketPath)
		defer func(conn *grpc.ClientConn) {
			err := conn.Close()
			if err != nil {
				panic(err)
			}
		}(conn)

		fmt.Print("Enter password: ")
		password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Println("Error reading password")
			return
		}
		fmt.Println()

		response, err := client.Unlock(context.Background(), &pb.PasswordRequest{Password: string(password)})
		if err != nil {
			log.Fatalf("%v", err)
		}
		if response.Success {
			fmt.Println("unlock successfully")
		} else {
			fmt.Println("unlock failed")
			fmt.Println(response.Error)
		}
	},
}

func init() {
	RootCmd.AddCommand(unlockCmd)
}
